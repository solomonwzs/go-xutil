package network

import (
	"fmt"
	"net"
	"sync"
	"testing"
	"unsafe"
)

func icmpRequest(tb testing.TB) {
	// conn, err := net.Dial("ip4:icmp", "59.66.1.1")
	conn, err := net.Dial("ip4:icmp", "1.1.1.1")
	if err != nil {
		tb.Fatal(err)
	}
	defer conn.Close()

	var msg [512]byte
	id := (*uint16)(unsafe.Pointer(&msg[4]))
	*id = 13
	seq := (*uint16)(unsafe.Pointer(&msg[6]))
	*seq = 37
	BuildIcmpHeader(msg[:], 8, 0)
	fmt.Println(msg[:8])

	_, err = conn.Write(msg[:8])
	if err != nil {
		tb.Fatal(err)
	}

	listener, err := net.ListenIP("ip4:icmp", nil)
	if err != nil {
		tb.Fatal(err)
	}
	defer listener.Close()

	n, err := listener.Read(msg[:])
	if err != nil {
		tb.Fatal(err)
	}
	fmt.Println("--", msg[:n])
}

func icmpReply(tb testing.TB) {
	conn, err := net.ListenIP("ip4:icmp", nil)
	if err != nil {
		tb.Fatal(err)
	}

	var msg [512]byte
	n, err := conn.Read(msg[:])
	if err != nil {
		tb.Fatal(err)
	}
	fmt.Println("<<", msg[:n])
}

func _TestDailIP(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		icmpRequest(t)
		wg.Done()
	}()
	icmpReply(t)

	wg.Wait()
}
