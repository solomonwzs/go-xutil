package dhcp

import (
	"fmt"
	"net"
	"syscall"
	"testing"
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
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM,
		syscall.IPPROTO_UDP)
	if err != nil {
		return
	}
	defer syscall.Close(fd)

	if err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET,
		syscall.SO_BROADCAST, 1); err != nil {
		t.Fatal(err)
	}

	addr0 := syscall.SockaddrInet4{
		Port: CLIENT_PORT,
		Addr: [4]byte{0, 0, 0, 0},
	}
	if err = syscall.Bind(fd, &addr0); err != nil {
		t.Fatal(err)
	}

	interf, err := net.InterfaceByName("eno1")
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1024)
	addr1 := syscall.SockaddrInet4{
		Port: SERVER_PORT,
		Addr: [4]byte{255, 255, 255, 255},
	}

	// discover
	msg := NewMessaageForInterface(interf)
	msg.SetMessageType(DHCPDISCOVER)
	msg.SetBroadcast()
	if err = syscall.Sendto(fd, msg.Marshal(), 0, &addr1); err != nil {
		t.Fatal(err)
	}
	fmt.Println(">>> discover")

	// offer
	n, from0, err := syscall.Recvfrom(fd, buf, 0)
	if err != nil {
		t.Fatal(err)
	}
	from := from0.(*syscall.SockaddrInet4)
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
	if err = syscall.Sendto(fd, msg.Marshal(), 0, &addr1); err != nil {
		t.Fatal(err)
	}
	fmt.Println(">>> request")

	// ack
	n, from0, err = syscall.Recvfrom(fd, buf, 0)
	if err != nil {
		t.Fatal(err)
	}
	from = from0.(*syscall.SockaddrInet4)
	respMsg, err = Unmarshal(buf[:n])
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("<<< ack")
	fmt.Println("from:", net.IP(from.Addr[:]))
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
