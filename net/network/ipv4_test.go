package network

import (
	"fmt"
	"syscall"
	"testing"
	"unsafe"
)

func TestIPv4(t *testing.T) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW,
		syscall.IPPROTO_RAW)
	if err != nil {
		t.Fatal(err)
	}
	defer syscall.Close(fd)

	addr := syscall.SockaddrInet4{
		Port: 0,
		Addr: [4]byte{127, 0, 0, 1},
	}

	syscall.Sendto(fd, []byte{}, 0, &addr)

	h := IPv4Header{}
	fmt.Println(unsafe.Sizeof(h))
}
