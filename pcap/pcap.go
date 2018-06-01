package pcap

/*
#cgo linux LDFLAGS: -lpcap

#include <pcap.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"errors"
	"net"
	"sync"
	"time"
	"unsafe"
)

var pcapCompileLock = &sync.Mutex{}

type Handle struct {
	handle *C.pcap_t
	interf *net.Interface
	netp   C.bpf_u_int32
	maskp  C.bpf_u_int32
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

// func UnsafeString(p unsafe.Pointer) string {
// 	var l uintptr = 0
// 	for ; *(*byte)(unsafe.Pointer(uintptr(p) + l)) != 0; l++ {
// 	}
// 	h := [2]unsafe.Pointer{p, unsafe.Pointer(l)}
// 	return *(*string)(unsafe.Pointer(&h))
// }

// func UnsafeBytes(p unsafe.Pointer, size int) []byte {
// 	l := unsafe.Pointer(uintptr(size))
// 	h := [3]unsafe.Pointer{p, l, l}
// 	return *(*[]byte)(unsafe.Pointer(&h))
// }

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
	h = &Handle{interf: interf}

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

	if C.pcap_compile(h.handle, &bpf, cExpr, 1, h.maskp) == -1 {
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
