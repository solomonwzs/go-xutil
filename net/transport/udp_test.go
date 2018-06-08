package transport

import (
	"net"
	"syscall"
	"testing"
	"time"

	"github.com/solomonwzs/goxutil/logger"
	"github.com/solomonwzs/goxutil/net/datalink"
	"github.com/solomonwzs/goxutil/net/network"
)

func TestChecksum(t *testing.T) {
	ipRaw := []byte{
		0x45, 0x00, 0x00, 0x2b, 0x7a, 0x8f, 0x00, 0x00,
		0x40, 0x11, 0x79, 0xf2, 0xc0, 0xa8, 0xc5, 0x98,
		0xff, 0xff, 0xff, 0xff,
	}
	ipH, _ := network.IPv4HeaderUnmarshal(ipRaw)

	u := &Udp{
		IpH:     ipH,
		SrcPort: 9999,
		DstPort: 10000,
		Length:  23,
		Data: []byte{
			0x56, 0x53, 0x54, 0x41, 0x52, 0x43, 0x41, 0x4d,
			0x51, 0x55, 0x45, 0x52, 0x59, 0x2c, 0x30,
		},
	}
	u.Marshal()
	logger.DPrintf("%x\n", u.Checksum)
}

func TestUdp(t *testing.T) {
	dev := "eno1"

	sock, err := datalink.NewDlSocket(dev, syscall.ETH_P_IP)
	if err != nil {
		t.Fatal(err)
	}
	defer sock.Close()

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

	ipH := &network.IPv4Header{
		Version:    4,
		TOS:        0,
		Id:         123,
		Flags:      network.IPV4_FLAG_DONT_FRAG,
		FragOffset: 0,
		TTL:        64,
		Protocol:   syscall.IPPROTO_UDP,
		SrcAddr:    localIP,
		DstAddr:    net.IPv4(120, 78, 185, 243),
	}

	u := &Udp{
		IpH:     ipH,
		SrcPort: 7777,
		DstPort: 8888,
		Data:    []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	}

	p1, _ := u.Marshal()
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
