package dhcp

import (
	"fmt"
	"net"
	"syscall"
	"testing"

	"github.com/solomonwzs/goxutil/net/ethernet"
	"github.com/solomonwzs/goxutil/net/transport"
	"github.com/solomonwzs/goxutil/net/xnetutil"
)

func _TestDHCP0(t *testing.T) {
	src := &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: CLIENT_PORT}
	dst := &net.UDPAddr{IP: net.IPv4(255, 255, 255, 255), Port: SERVER_PORT}

	conn, err := net.DialUDP("udp", src, dst)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	f, err := conn.File()
	if err != nil {
		t.Fatal(err)
	}
	fd := int(f.Fd())
	syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)

	interf, err := net.InterfaceByName("eno1")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(interf.HardwareAddr)

	msg := NewMessaageForInterface(interf)
	msg.SetMessageType(DHCPDISCOVER)
	msg.SetBroadcast()

	conn.Write(msg.Marshal())

	buf := make([]byte, 1024)
	f.Read(buf)
}

func TestDHCP1(t *testing.T) {
	conn, err := transport.NewUDPBroadcastConn(CLIENT_PORT, SERVER_PORT)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	interf, err := net.InterfaceByName("eno1")
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 1024)

	// discover
	msg := NewMessaageForInterface(interf)
	msg.SetMessageType(DHCPDISCOVER)
	msg.SetBroadcast()
	if _, err = conn.Write(msg.Marshal()); err != nil {
		t.Fatal(err)
	}
	fmt.Println(">>> discover")

	// offer
	n, from, err := conn.Readfrom(buf)
	if err != nil {
		t.Fatal(err)
	}
	respMsg, err := Unmarshal(buf[:n])
	if err != nil {
		t.Fatal(err)
	}
	clientIP := respMsg.ClientIP()
	dhcpServerIP, _ := respMsg.DHCPServerID()
	mask, _ := respMsg.SubnetMask()
	router, _ := respMsg.Router()
	leaseTime, _ := respMsg.AddressLeaseTime()
	fmt.Println("<<< offer")
	fmt.Printf("from:        %s\n", net.IP(from.Addr[:]))
	fmt.Printf("client IP:   %s\n", clientIP)
	fmt.Printf("DHCP server: %s\n", dhcpServerIP)
	fmt.Printf("subnet mask: %s\n", mask)
	fmt.Printf("router:      %s\n", router)
	fmt.Printf("lease time:  %d\n", leaseTime)

	// request
	msg = NewMessaageForInterface(interf)
	msg.SetBroadcast()
	msg.SetMessageType(DHCPREQUEST)
	msg.SetOptions(OPT_ADDR_REQUEST, []byte(clientIP))
	msg.SetOptions(OPT_DHCP_SERVER_ID, []byte(dhcpServerIP))
	if _, err = conn.Write(msg.Marshal()); err != nil {
		t.Fatal(err)
	}
	fmt.Println(">>> request")

	// ack
	n, from, err = conn.Readfrom(buf)
	if err != nil {
		t.Fatal(err)
	}
	respMsg, err = Unmarshal(buf[:n])
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("<<< ack")
	clientIP = respMsg.ClientIP()
	dhcpServerIP, _ = respMsg.DHCPServerID()
	mask, _ = respMsg.SubnetMask()
	router, _ = respMsg.Router()
	leaseTime, _ = respMsg.AddressLeaseTime()
	fmt.Printf("from:        %s\n", net.IP(from.Addr[:]))
	fmt.Printf("client IP:   %s\n", clientIP)
	fmt.Printf("DHCP server: %s\n", dhcpServerIP)
	fmt.Printf("subnet mask: %s\n", mask)
	fmt.Printf("router:      %s\n", router)
	fmt.Printf("lease time:  %d\n", leaseTime)
}

func _TestUDP0(t *testing.T) {
	addr := &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: SERVER_PORT}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(n)
	}
}

func _TestUDP1(t *testing.T) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM,
		syscall.IPPROTO_UDP)
	if err != nil {
		return
	}
	defer syscall.Close(fd)

	if err = syscall.SetNonblock(fd, true); err != nil {
		t.Fatal(err)
	}

	// if err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET,
	// 	syscall.SO_BROADCAST, 1); err != nil {
	// 	t.Fatal(err)
	// }

	addr := syscall.SockaddrInet4{
		Port: SERVER_PORT,
		Addr: [4]byte{0, 0, 0, 0},
	}
	if err = syscall.Bind(fd, &addr); err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1024)
	for {
		n, from, err := syscall.Recvfrom(fd, buf, 0)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(n, from)
	}
}

func TestDHCP2(t *testing.T) {
	interf, err := net.InterfaceByName("eno1")
	if err != nil {
		panic(err)
	}

	msg := NewMessaageForInterface(interf)
	msg.SetMessageType(DHCPDISCOVER)
	msg.SetBroadcast()
	raw, err := transport.NewBroadcastUDPRaw(interf, CLIENT_PORT,
		SERVER_PORT, msg.Marshal())
	if err != nil {
		panic(err)
	}
	fmt.Printf("% x\n", raw)

	fd0, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW,
		int(xnetutil.Htons(syscall.ETH_P_ALL)))
	if err != nil {
		panic(err)
	}
	defer syscall.Close(fd0)

	fd1, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW,
		int(xnetutil.Htons(syscall.ETH_P_IP)))
	if err != nil {
		panic(err)
	}
	defer syscall.Close(fd1)

	addr := &syscall.SockaddrLinklayer{
		Ifindex: interf.Index,
	}
	err = syscall.Sendto(fd0, raw, 0, addr)
	if err != nil {
		panic(err)
	}

	p := make([]byte, 1024)
	for {
		_, _, err := syscall.Recvfrom(fd1, p, 0)
		if err != nil {
			panic(err)
		}

		if len(p) < ethernet.SIZEOF_ETH_HEADER {
			continue
		}

		_, u, err := transport.RawUDPUnmarshal(p)
		if err != nil || u.DstPort != CLIENT_PORT {
			continue
		}

		rMsg, err := Unmarshal(u.Data)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%+v\n", rMsg)
		break
	}
}
