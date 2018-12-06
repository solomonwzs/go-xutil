package datalink

import (
	"fmt"
	"net"
	"syscall"
	"testing"

	"github.com/solomonwzs/goxutil/net/ethernet"
	"github.com/solomonwzs/goxutil/net/xnetutil"
)

func _TestRecv(t *testing.T) {
	// fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW,
	// 	int(ethernet.Htons(syscall.ETH_P_ALL)))
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW,
		int(xnetutil.Htons(syscall.ETH_P_ARP)))
	if err != nil {
		t.Fatal(err)
	}
	defer syscall.Close(fd)

	buf := make([]byte, 1024)
	// typ := (*uint16)(unsafe.Pointer(&buf[12]))
	for {
		numRead, _, err := syscall.Recvfrom(fd, buf, 0)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("% x\n", buf[:numRead])
		// if ethernet.Ntohs(*typ) == syscall.ETH_P_ARP {
		// 	fmt.Printf("% x\n", buf[:numRead])
		// }
	}
}

func _TestARP(t *testing.T) {
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW,
		int(xnetutil.Htons(syscall.ETH_P_ALL)))
	if err != nil {
		t.Fatal(err)
	}
	defer syscall.Close(fd)

	interf, err := net.InterfaceByName("eno1")
	if err != nil {
		t.Fatal(err)
	}
	addrs, _ := interf.Addrs()
	for _, a := range addrs {
		fmt.Println(a)
		fmt.Println(net.ParseCIDR(a.String()))
	}

	ethH := &ethernet.EthernetHeader{
		Src:  interf.HardwareAddr,
		Dst:  ethernet.ETH_BROADCAST_ADDR,
		Type: syscall.ETH_P_ARP,
	}

	arp := Arp{
		HardwareType: 1,
		ProtocolType: syscall.ETH_P_IP,
		HardwareSize: 6,
		ProtocolSize: 4,
		Opcode:       ARP_OPC_REQUEST,
		SHA:          interf.HardwareAddr,
		SPA:          net.IP{192, 168, 197, 130},
		THA:          make([]uint8, 6, 6),
		TPA:          net.IP{192, 168, 197, 252},
	}

	addr := syscall.SockaddrLinklayer{
		Ifindex: interf.Index,
	}

	p0, _ := ethH.Marshal()
	p1, _ := arp.Marshal()
	p := append(p0, p1...)
	err = syscall.Sendto(fd, p, 0, &addr)
	if err != nil {
		t.Fatal(err)
	}
}

func _TestT(t *testing.T) {
	interf, err := net.InterfaceByName("eno1")
	if err != nil {
		t.Fatal(err)
	}

	addrs, _ := interf.Addrs()
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok {
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				addr := &net.IPNet{
					IP:   ip4,
					Mask: ipnet.Mask[len(ipnet.Mask)-4:],
				}
				fmt.Println(ipnet.Mask)
				fmt.Println(addr.Mask)
				break
			}
		}
	}
}

func TestGetHardwareAddr(t *testing.T) {
	fmt.Println(GetHardwareAddr("eno1", net.IP{192, 168, 197, 252}, 0))
}
