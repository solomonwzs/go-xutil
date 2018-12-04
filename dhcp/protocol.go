package dhcp

import (
	"math/rand"
	"net"
	"unsafe"
)

const (
	// UDP port
	SRC_PORT = 68
	DST_PORT = 67

	// Message op
	BOOTREQUEST = 1
	BOOTREPLY   = 2

	// Hardware address type
	HTYPE_ETHERNET = 1

	// Hardware address length
	HLEN_ETHERNET = 6

	// Optionals
	OPT_HOSTNAME  = 12
	OPT_MSG_TYPE  = 53
	OPT_PARA_REQ  = 55
	OPT_CLASS_ID  = 60
	OPT_CLIENT_ID = 61
	OPT_END       = 255

	// DHCP message type
	DHCPDISCOVER = 1
	DHCPOFFER
	DHCPREQUEST
	DHCPDECLINE
	DHCPACK
	DHCPNAK
	DHCPRELEASE
	DHCPINFORM
	DHCPFORCERENEW
	DHCPLEASEQUERY
	DHCPLEASEUNASSIGNED
	DHCPLEASEUNKNOWN
	DHCPLEASEACTIVE
	DHCPBULKLEASEQUERY
	DHCPLEASEQUERYDONE
	DHCPACTIVELEASEQUERY
	DHCPLEASEQUERYSTATUS
	DHCPTLS

	MSG_FIX_SIZE = unsafe.Sizeof(MessageFix{})
)

var (
	_TransactionID uint32
	_Rand          *rand.Rand

	_SrcAddr *net.UDPAddr
	_DstAddr *net.UDPAddr
)

type MessageFix struct {
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
