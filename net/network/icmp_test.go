package network

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func testIcmpSendLoop(tb testing.TB) {
	// conn, err := net.Dial("ip4:icmp", "59.66.1.1")
	// conn, err := net.Dial("ip4:icmp", "1.1.1.1")
	conn, err := net.Dial("ip4:icmp", "120.78.185.243")
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

	for {
		_, err = conn.Write(msg)
		if err != nil {
			tb.Fatal(err)
		}
		time.Sleep(1 * time.Second)
	}
}

func testIcmpRecvLoop(tb testing.TB) {
	conn, err := net.ListenIP("ip4:icmp", nil)
	if err != nil {
		tb.Fatal(err)
	}
	var msg [512]byte

	for {
		n, err := conn.Read(msg[:])
		if err != nil {
			tb.Fatal(err)
		}

		h, _ := IPv4HeaderUnmarshal(msg[:n])
		fmt.Printf("%+v\n", h)
	}
}

func _TestIcmp(t *testing.T) {
	go testIcmpSendLoop(t)
	go testIcmpRecvLoop(t)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGINT,
		syscall.SIGTERM)
	<-c
}
