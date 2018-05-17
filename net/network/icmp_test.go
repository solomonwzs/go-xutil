package network

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func icmpRequest(tb testing.TB) {
	// conn, err := net.Dial("ip4:icmp", "59.66.1.1")
	conn, err := net.Dial("ip4:icmp", "1.1.1.2")
	// conn, err := net.Dial("ip4:icmp", "192.168.197.1")
	if err != nil {
		tb.Fatal(err)
	}
	defer conn.Close()

	icmp := Icmp{
		Type: ICMP_CT_ECHO_REQUEST,
		Code: 0,
		Data: &IcmpEcho{
			Id:     123,
			SeqNum: 456,
		},
	}
	msg, _ := icmp.Marshal()
	_, err = conn.Write(msg)
	if err != nil {
		tb.Fatal(err)
	}

	listener, err := net.ListenIP("ip4:icmp", nil)
	if err != nil {
		tb.Fatal(err)
	}
	defer listener.Close()

	// n, err := listener.Read(msg[:])
	// if err != nil {
	// 	tb.Fatal(err)
	// }
	// fmt.Println("--", msg[:n])
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
	// icmpReply(t)

	wg.Wait()
}
