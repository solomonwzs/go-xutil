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

	msg := NewMessage()
	msg.fix.Op = BOOTREQUEST
	msg.fix.Htype = HTYPE_ETHERNET
	msg.fix.Hlen = HLEN_ETHERNET
	copy(msg.fix.Chaddr[:], interf.HardwareAddr)
	msg.SetMassageType(DHCPDISCOVER)
	msg.SetClientID(HTYPE_ETHERNET, interf.HardwareAddr)
	// msg.SetBroadcast()

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

	c := make(chan struct{})
	go func() {
		buf := make([]byte, 1024)
		if n, from, err := syscall.Recvfrom(fd, buf, 0); err != nil {
			t.Fatal(err)
		} else {
			fmt.Println(n, from)
			close(c)
		}
	}()

	interf, err := net.InterfaceByName("eno1")
	if err != nil {
		t.Fatal(err)
	}

	msg := NewMessage()
	msg.fix.Op = BOOTREQUEST
	msg.fix.Htype = HTYPE_ETHERNET
	msg.fix.Hlen = HLEN_ETHERNET
	copy(msg.fix.Chaddr[:], interf.HardwareAddr)
	msg.SetMassageType(DHCPDISCOVER)
	msg.SetBroadcast()
	msg.SetClientID(HTYPE_ETHERNET, interf.HardwareAddr)

	addr1 := syscall.SockaddrInet4{
		Port: SERVER_PORT,
		Addr: [4]byte{255, 255, 255, 255},
	}
	if err = syscall.Sendto(fd, msg.Marshal(), 0, &addr1); err != nil {
		t.Fatal(err)
	}

	<-c
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
