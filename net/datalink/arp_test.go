package arp

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"testing"

	"github.com/solomonwzs/goxutil/net/ethernet"
)

func TestRecv(t *testing.T) {
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW,
		syscall.ETH_P_ALL)
	if err != nil {
		t.Fatal(err)
	}

	interf, err := net.InterfaceByName("eno1")
	if err != nil {
		t.Fatal(err)
	}

	addr := syscall.SockaddrLinklayer{
		Protocol: syscall.ETH_P_ARP,
		Ifindex:  interf.Index,
	}

	syscall.Bind(fd, &addr)
	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))

	for {
		buf := make([]byte, 1024)
		numRead, err := f.Read(buf)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("% X\n", buf[:numRead])
	}
}

func _TestARP(t *testing.T) {
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW,
		syscall.ETH_P_ALL)
	if err != nil {
		t.Fatal(err)
	}
	defer syscall.Close(fd)

	interf, err := net.InterfaceByName("eno1")
	if err != nil {
		t.Fatal(err)
	}

	ethH := &ethernet.EthernetHeader{
		Src:  []uint8(interf.HardwareAddr),
		Dst:  []uint8{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		Type: syscall.ETH_P_ARP,
	}

	arp := Arp{
		HardwareType: 1,
		ProtocolType: syscall.ETH_P_IP,
		HardwareSize: 6,
		ProtocolSize: 4,
		Opcode:       ARP_OPC_REQUEST,
		SHA:          []uint8(interf.HardwareAddr),
		SPA:          []uint8{192, 168, 197, 130},
		THA:          make([]uint8, 6, 6),
		TPA:          []uint8{192, 168, 197, 252},
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
