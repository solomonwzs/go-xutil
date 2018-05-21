package network

import (
	"bytes"
	"net"
	"syscall"
	"testing"

	"github.com/solomonwzs/goxutil/net/ethernet"
)

func _TestIPv4(t *testing.T) {
	// fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW,
	// 	syscall.IPPROTO_RAW)
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW,
		syscall.ETH_P_ALL)
	if err != nil {
		t.Fatal(err)
	}
	defer syscall.Close(fd)

	buf := new(bytes.Buffer)

	ipH := IPv4Header{
		Version:    4,
		TOS:        0,
		Id:         123,
		Flags:      IPV4_FLAG_DONT_FRAG,
		FragOffset: 0,
		TTL:        64,
		Protocol:   syscall.IPPROTO_ICMP,
		SrcAddr:    net.IPv4(192, 168, 197, 130),
		DstAddr:    net.IPv4(120, 78, 185, 243),
	}
	icmp := Icmp{
		Type: ICMP_CT_ECHO_REQUEST,
		Code: 0,
		Data: &IcmpEcho{
			Id:     123,
			SeqNum: 456,
		},
	}
	ethH := &ethernet.EthernetHeader{
		Src:  []uint8{0x4c, 0xcc, 0x6a, 0xac, 0xe5, 0x63},
		Dst:  []uint8{0x1c, 0xab, 0x34, 0x12, 0x63, 0xde},
		Type: syscall.ETH_P_IP,
	}

	ethP, _ := ethH.Marshal()
	buf.Write(ethP)

	icmpB, err := icmp.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	ipH.Length = SIZEOF_IPV4_HEADER + uint16(len(icmpB))

	if p, err := ipH.Marshal(); err != nil {
		t.Fatal(err)
	} else {
		buf.Write(p)
	}
	buf.Write(icmpB)

	// addr := syscall.SockaddrInet4{
	// 	Port: 0,
	// 	Addr: [4]byte{192, 168, 197, 130},
	// }

	interf, err := net.InterfaceByName("eno1")
	if err != nil {
		t.Fatal(err)
	}
	var addr syscall.SockaddrLinklayer
	// addr.Protocol = syscall.ETH_P_IP
	addr.Ifindex = interf.Index
	// addr.Hatype = syscall.ARPHRD_ETHER

	err = syscall.Sendto(fd, buf.Bytes(), 0, &addr)
	if err != nil {
		t.Fatal(err)
	}
}
