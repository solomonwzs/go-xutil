package pcap

/*
#cgo linux LDFLAGS: -lpcap

#include <pcap.h>
#include <stdint.h>
#include <stdlib.h>
#include <sys/select.h>

#define SIZEOF_PCAP_IF_T sizeof(pcap_if_t)

int
_pcap_next_ex(pcap_t *p, uintptr_t hdr, uintptr_t pkt) {
	return pcap_next_ex(p, (struct pcap_pkthdr**)hdr, (const u_char**) pkt);
}

int
_pcap_wait(pcap_t *p, int to_us) {
	int fd = pcap_get_selectable_fd(p);
	if (fd < 0) {
		return -1;
	}

	fd_set fds;
	FD_ZERO(&fds);
	FD_SET(fd, &fds);

	int n;
	if (to_us != 0) {
		struct timeval tv;
		tv.tv_sec = to_us / 1000000;
		tv.tv_usec = to_us % 1000000;
		n = select(fd + 1, &fds, NULL, NULL, &tv);
	} else {
		n = select(fd + 1, &fds, NULL, NULL, NULL);
	}
	if (n == -1) {
		perror("select: ");
	}
	return n;
}
*/
import "C"
import (
	"errors"
	"net"
	"reflect"
	"sync"
	"time"
	"unsafe"

	"github.com/solomonwzs/goxutil/closer"
)

var (
	ErrHandlerWasClosed = errors.New("[pcap] handler was closed")
)

const (
	PcapNextOK         = 1
	PcapNextErrTimeout = 0
	PcapNextErrRead    = -1
	PcapNextErrNoMore  = -2
)

type CapturePacket struct {
	Ts     time.Time
	Caplen uint32
	Len    uint32
	Data   []byte
}

type PcapNextError int

func (e PcapNextError) Error() string {
	switch e {
	case PcapNextOK:
		return "OK"
	case PcapNextErrTimeout:
		return "Timeout"
	case PcapNextErrRead:
		return "Read Error"
	case PcapNextErrNoMore:
		return "No more packet"
	}
	return "Unknown"
}

var pcapCompileLock = &sync.Mutex{}

type Handle struct {
	closer.Closer
	handle *C.pcap_t
	interf *net.Interface
	netp   C.bpf_u_int32
	maskp  C.bpf_u_int32

	pkthdr  *C.struct_pcap_pkthdr
	packet  *C.u_char
	pktLock *sync.Mutex
}

func UnsafeSlice(ptr unsafe.Pointer, len int) (p []byte) {
	slice := (*reflect.SliceHeader)(unsafe.Pointer(&p))
	slice.Data = uintptr(ptr)
	slice.Len = len
	slice.Cap = len

	return
}

func DatalinkName(typ int) string {
	name := C.pcap_datalink_val_to_name(C.int(typ))
	return C.GoString(name)
}

func DatalinkDesc(typ int) string {
	desc := C.pcap_datalink_val_to_description(C.int(typ))
	return C.GoString(desc)
}

func pcapError(errBuf []byte) error {
	i := 0
	for ; errBuf[i] != 0 && i < len(errBuf); i++ {
	}
	return errors.New(string(errBuf[:i]))
}

func charptr(errBuf []byte) *C.char {
	return (*C.char)(unsafe.Pointer(&errBuf[0]))
}

func PcapFindAllDevs() ([]string, error) {
	var devsp *C.pcap_if_t
	errBuf := make([]byte, C.PCAP_ERRBUF_SIZE)

	if C.pcap_findalldevs(&devsp, charptr(errBuf)) != 0 {
		return nil, pcapError(errBuf)
	} else {
		defer C.pcap_freealldevs(devsp)

		devspn := []string{C.GoString(devsp.name)}
		next := devsp.next
		for next != nil {
			devspn = append(devspn, C.GoString(next.name))
			next = next.next
		}
		return devspn, nil
	}
}

// func PcapLookupDev() (string, error) {
// 	errBuf := make([]byte, C.PCAP_ERRBUF_SIZE)
// 	dev := C.pcap_lookupdev(charptr(errBuf))

// 	if dev == nil {
// 		return "", pcapError(errBuf)
// 	}

// 	return C.GoString(dev), nil
// }

func OpenLive(dev string, snaplen int, promisc bool, timeout time.Duration) (
	h *Handle, err error) {
	interf, err := net.InterfaceByName(dev)
	if err != nil {
		return
	}
	h = &Handle{interf: interf, pktLock: &sync.Mutex{}}

	cDev := C.CString(dev)
	defer C.free(unsafe.Pointer(cDev))

	var pro C.int = 0
	if promisc {
		pro = 1
	}

	errBuf := make([]byte, C.PCAP_ERRBUF_SIZE)
	if C.pcap_lookupnet(cDev, &h.netp, &h.maskp, charptr(errBuf)) == -1 {
		return nil, pcapError(errBuf)
	}

	h.handle = C.pcap_open_live(
		cDev,
		C.int(snaplen),
		pro,
		C.int(timeout/time.Millisecond),
		charptr(errBuf),
	)
	if h.handle == nil {
		return nil, pcapError(errBuf)
	}

	// if h.timeout > 0 {
	// 	if C.pcap_setnonblock(h.handle, 1, charptr(errBuf)) == -1 {
	// 		C.pcap_close(h.handle)
	// 		return nil, pcapError(errBuf)
	// 	}
	// }

	h.Closer = closer.NewCloser(func() error {
		C.pcap_close(h.handle)
		return nil
	})

	return
}

func (h *Handle) Error() error {
	return errors.New(C.GoString(C.pcap_geterr(h.handle)))
}

func (h *Handle) DatalinkType() int {
	return int(C.pcap_datalink(h.handle))
}

func (h *Handle) compile(expr string) (
	bpf C.struct_bpf_program, err error) {
	cExpr := C.CString(expr)
	defer C.free(unsafe.Pointer(cExpr))

	pcapCompileLock.Lock()
	defer pcapCompileLock.Unlock()

	if C.pcap_compile(h.handle, &bpf, cExpr, 1, h.netp) == -1 {
		return bpf, h.Error()
	}
	return
}

func (h *Handle) SetFilter(expr string) (err error) {
	bpf, err := h.compile(expr)
	if err != nil {
		return
	}
	defer C.pcap_freecode(&bpf)

	if C.pcap_setfilter(h.handle, &bpf) == -1 {
		return h.Error()
	}
	return nil
}

func (h *Handle) ReadPacket(copy bool) (packet CapturePacket, err error) {
	// if h.timeout > 0 {
	// 	var ready bool
	// 	if ready, err = h.waitForRead(h.timeout); err != nil {
	// 		return
	// 	} else if !ready {
	// 		return packet, PcapNextError(PcapNextErrTimeout)
	// 	}
	// }

	h.pktLock.Lock()
	defer h.pktLock.Unlock()

	hdrP := C.uintptr_t(uintptr(unsafe.Pointer(&h.pkthdr)))
	pktP := C.uintptr_t(uintptr(unsafe.Pointer(&h.packet)))
	if res := PcapNextError(
		C._pcap_next_ex(h.handle, hdrP, pktP)); res != PcapNextOK {
		return packet, res
	}

	packet.Ts = time.Unix(
		int64(h.pkthdr.ts.tv_sec),
		int64(h.pkthdr.ts.tv_usec)*1000,
	)
	packet.Caplen = uint32(h.pkthdr.caplen)
	packet.Len = uint32(h.pkthdr.len)
	if copy {
		packet.Data = C.GoBytes(unsafe.Pointer(h.packet),
			C.int(h.pkthdr.len))
	} else {
		packet.Data = UnsafeSlice(unsafe.Pointer(h.packet),
			int(h.pkthdr.len))
	}
	return
}

func (h *Handle) WaitForRead(timeout time.Duration) (ready bool, err error) {
	if n := C._pcap_wait(
		h.handle, C.int(timeout/time.Microsecond)); n < 0 {
		err = errors.New("[pcap] wait fail")
		return
	} else {
		return n != 0, nil
	}
}
