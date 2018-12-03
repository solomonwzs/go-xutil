package dhcp

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync/atomic"
	"testing"
	"unsafe"
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

	buf := make([]byte, MSG_FIX_SIZE)
	msg := (*Message)(unsafe.Pointer(&buf[0]))

	msg.Op = BOOTREQUEST
	msg.Htype = HTYPE_10MB_ETH
	msg.Hlen = HLEN_10MB_ETH
	binary.BigEndian.PutUint32(
		msg.Xid[:], atomic.AddUint32(&_TransactionID, 1))
	copy(msg.Chaddr[:], interf.HardwareAddr)

	fmt.Println(buf)
}
