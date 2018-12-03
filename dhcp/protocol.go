package dhcp

import (
	"math/rand"
	"net"
	"unsafe"
)

const (
	SRC_PORT = 68
	DST_PORT = 67

	BOOTREQUEST = 1
	BOOTREPLY   = 2

	HTYPE_10MB_ETH = 1

	HLEN_10MB_ETH = 6

	MSG_FIX_SIZE = unsafe.Sizeof(Message{})
)

var (
	_TransactionID uint32
	_Rand          *rand.Rand

	_SrcAddr *net.UDPAddr
	_DstAddr *net.UDPAddr
)

type Message struct {
	Op     byte
	Htype  byte
	Hlen   byte
	Hops   byte
	Xid    [4]byte
	Secs   [2]byte
	Flags  [2]byte
	Ciaddr [4]byte
	Yiaddr [4]byte
	Siaddr [4]byte
	Giaddr [4]byte
	Chaddr [16]byte
	Sname  [64]byte
	File   [128]byte
}
