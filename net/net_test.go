package net

import (
	"net"
	"syscall"
	"testing"
	"time"

	"github.com/solomonwzs/goxutil/net/datalink"
	"github.com/solomonwzs/goxutil/net/network"
)

func TestICMP0(t *testing.T) {
	interf, err := net.InterfaceByName("eno1")
	if err != nil {
		t.Fatal(err)
	}

	sock, err := datalink.NewDlSocket(interf, syscall.ETH_P_IP)
	if err != nil {
		t.Fatal(err)
	}
	defer sock.Close()

	addrs, err := interf.Addrs()
	if err != nil {
		return
	}
	var localIP net.IP = nil
	for _, addr := range addrs {
		if ip, _, err := net.ParseCIDR(addr.String()); err != nil {
			t.Fatal(err)
		} else if ipv4 := ip.To4(); ipv4 != nil {
			localIP = ipv4
			break
		}
	}
	if localIP == nil {
		t.Fatal("can not get local ip")
	}

	ipH := &network.IPv4Header{
		Version:    4,
		TOS:        0,
		Id:         123,
		Flags:      network.IPV4_FLAG_DONT_FRAG,
		FragOffset: 0,
		TTL:        64,
		Protocol:   syscall.IPPROTO_ICMP,
		SrcAddr:    localIP,
		DstAddr:    net.IPv4(59, 66, 1, 1),
	}
	icmp := &network.Icmp{
		Type: network.ICMP_CT_ECHO_REQUEST,
		Code: 0,
		Data: &network.IcmpEcho{
			Id:     456,
			SeqNum: 789,
		},
	}

	p1, _ := icmp.Marshal()
	ipH.Length = network.SIZEOF_IPV4_HEADER + uint16(len(p1))
	p0, _ := ipH.Marshal()
	p0 = append(p0, p1...)

	for {
		_, err = sock.Write(p0)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(1 * time.Second)
	}
}
