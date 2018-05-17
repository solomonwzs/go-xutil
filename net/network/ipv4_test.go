package network

import (
	"bytes"
	"net"
	"syscall"
	"testing"
)

func TestIPv4(t *testing.T) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW,
		syscall.IPPROTO_RAW)
	if err != nil {
		t.Fatal(err)
	}
	defer syscall.Close(fd)

	ipH := IPv4Header{
		Version:    4,
		TOS:        0,
		Id:         123,
		Flags:      IPV4_FLAG_DONT_FRAG,
		FragOffset: 0,
		TTL:        64,
		Protocol:   IPV4_PRO_ICMP,
		SrcAddr:    net.IPv4(192, 168, 197, 130),
		DstAddr:    net.IPv4(192, 168, 197, 128),
	}
	icmp := Icmp{
		Type: ICMP_CT_ECHO_REQUEST,
		Code: 0,
		Data: &IcmpEcho{
			Id:     123,
			SeqNum: 456,
		},
	}

	icmpB, err := icmp.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	ipH.Length = SIZEOF_IPV4_HEADER + uint16(len(icmpB))

	buf := new(bytes.Buffer)
	if p, err := ipH.Marshal(); err != nil {
		t.Fatal(err)
	} else {
		buf.Write(p)
	}
	buf.Write(icmpB)

	addr := syscall.SockaddrInet4{
		Port: 0,
		Addr: [4]byte{192, 168, 197, 130},
	}

	err = syscall.Sendto(fd, buf.Bytes(), 0, &addr)
	if err != nil {
		t.Fatal(err)
	}
}
