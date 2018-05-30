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
	"time"
	"unsafe"
)

type Handle struct {
	handle *C.pcap_t
	interf *net.Interface
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

func UnsafeString(p unsafe.Pointer) string {
	var l uintptr = 0
	for ; *(*byte)(unsafe.Pointer(uintptr(p) + l)) != 0; l++ {
	}
	h := [2]unsafe.Pointer{p, unsafe.Pointer(l)}
	return *(*string)(unsafe.Pointer(&h))
}

func UnsafeBytes(p unsafe.Pointer, size int) []byte {
	l := unsafe.Pointer(uintptr(size))
	h := [3]unsafe.Pointer{p, l, l}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func PcapLookupDev() (string, error) {
	errBuf := make([]byte, C.PCAP_ERRBUF_SIZE)
	dev := C.pcap_lookupdev((*C.char)(unsafe.Pointer(&errBuf[0])))

	if dev == nil {
		return "", pcapError(errBuf)
	}

	size := C.strlen(dev)
	b := UnsafeBytes(unsafe.Pointer(dev), int(size))

	return string(b), nil
}

func OpenLive(dev string, snaplen int, promisc bool, toMs time.Duration) (
	h *Handle, err error) {
	device := C.CString(dev)
	defer C.free(unsafe.Pointer(device))

	var pro C.int = 0
	if promisc {
		pro = 1
	}

	errBuf := make([]byte, C.PCAP_ERRBUF_SIZE)
	handle := C.pcap_open_live(
		device,
		C.int(snaplen),
		pro,
		C.int(toMs/time.Millisecond),
		(*C.char)(unsafe.Pointer(&errBuf[0])),
	)

	if handle == nil {
		return nil, pcapError(errBuf)
	}

	return &Handle{
		handle: handle,
	}, nil
}

func (h *Handle) GetDatalink() (typ int) {
	typ = int(C.pcap_datalink(h.handle))
	return
}

func (h *Handle) Compile(filter string) {
}
