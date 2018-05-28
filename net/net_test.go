package net

import (
	"bytes"
	"net"
	"syscall"
	"testing"
	"time"

	"github.com/solomonwzs/goxutil/net/datalink"
	"github.com/solomonwzs/goxutil/net/ethernet"
	"github.com/solomonwzs/goxutil/net/network"
	"github.com/solomonwzs/goxutil/net/xnetutil"
)

func TestICMP(t *testing.T) {
	dev := "eno1"

	gateway, err := xnetutil.GetGateway(dev)
	if err != nil {
		t.Fatal(err)
	}

	hardwareAddr, err := xnetutil.GetHardwareAddr(dev, gateway)
	if err == xnetutil.ERR_NOT_FOUND {
		hardwareAddr, err = datalink.GetHardwareAddr(dev, gateway, 0)
	}
	if err != nil {
		t.Fatal(err)
	}

	interf, err := net.InterfaceByName(dev)
	if err != nil {
		t.Fatal(err)
	}
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

	ethH := &ethernet.EthernetHeader{
		Src:  interf.HardwareAddr,
		Dst:  hardwareAddr,
		Type: syscall.ETH_P_IP,
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
		// DstAddr:    net.IPv4(120, 78, 185, 243),
		DstAddr: net.IPv4(59, 66, 1, 1),
		// DstAddr: net.IPv4(1, 1, 1, 1),
	}
	icmp := &network.Icmp{
		Type: network.ICMP_CT_ECHO_REQUEST,
		Code: 0,
		Data: &network.IcmpEcho{
			Id:     456,
			SeqNum: 789,
		},
	}

	p0, _ := ethH.Marshal()
	p2, _ := icmp.Marshal()
	ipH.Length = network.SIZEOF_IPV4_HEADER + uint16(len(p2))
	p1, _ := ipH.Marshal()

	buf := new(bytes.Buffer)
	buf.Write(p0)
	buf.Write(p1)
	buf.Write(p2)

	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW,
		int(xnetutil.Htons(syscall.ETH_P_ALL)))
	if err != nil {
		t.Fatal(err)
	}
	defer syscall.Close(fd)
	to := syscall.SockaddrLinklayer{
		Ifindex: interf.Index,
	}

	for {
		err = syscall.Sendto(fd, buf.Bytes(), 0, &to)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(1 * time.Second)
	}
}
