package dhcp

import (
	"fmt"
	"net"
	"testing"
)

func TestDHCP(t *testing.T) {
	conn, err := net.DialUDP("udp", _SrcAddr, _DstAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

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
	msg.SetCookie([]byte("DHCP"))
	msg.SetMassageType(DHCPDISCOVER)
	msg.SetClientID(HTYPE_ETHERNET, interf.HardwareAddr)

	conn.Write(msg.Marshal())
}
