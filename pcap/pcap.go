package pcap

/*
#cgo linux LDFLAGS: -lpcap

#include <pcap.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

int
_pcap_next_ex(pcap_t *p, uintptr_t hdr, uintptr_t pkt) {
	return pcap_next_ex(p, (struct pcap_pkthdr**)hdr, (const u_char**) pkt);
}
*/
import "C"
import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
	"unsafe"

	"github.com/solomonwzs/goxutil/closer"
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

func PcapLookupDev() (string, error) {
	errBuf := make([]byte, C.PCAP_ERRBUF_SIZE)
	dev := C.pcap_lookupdev(charptr(errBuf))

	if dev == nil {
		return "", pcapError(errBuf)
	}

	return C.GoString(dev), nil
}

func OpenLive(dev string, snaplen int, promisc bool, toMs time.Duration) (
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
		C.int(toMs/time.Millisecond),
		charptr(errBuf),
	)

	if h.handle == nil {
		return nil, pcapError(errBuf)
	}
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

func (h *Handle) ReadPacket() (packet CapturePacket, err error) {
	h.pktLock.Lock()
	defer h.pktLock.Unlock()

	hdrP := C.uintptr_t(uintptr(unsafe.Pointer(&h.pkthdr)))
	pktP := C.uintptr_t(uintptr(unsafe.Pointer(&h.packet)))
	if res := PcapNextError(
		C._pcap_next_ex(h.handle, hdrP, pktP)); res != PcapNextOK {
		return packet, res
	}
	defer C.free(unsafe.Pointer(h.packet))

	packet.Ts = time.Unix(
		int64(h.pkthdr.ts.tv_sec),
		int64(h.pkthdr.ts.tv_usec)*1000,
	)
	packet.Caplen = uint32(h.pkthdr.caplen)
	packet.Len = uint32(h.pkthdr.len)
	packet.Data = C.GoBytes(unsafe.Pointer(h.packet),
		C.int(h.pkthdr.len))
	fmt.Println(uintptr(unsafe.Pointer(&packet.Data[0])))
	fmt.Println(uintptr(unsafe.Pointer(unsafe.Pointer(h.packet))))
	return
}
