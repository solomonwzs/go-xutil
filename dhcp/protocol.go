package dhcp

import (
	"math/rand"
	"unsafe"
)

const (
	// UDP port
	CLIENT_PORT = 68
	SERVER_PORT = 67

	// Message op
	BOOTREQUEST = 1
	BOOTREPLY   = 2

	// Hardware address type
	HTYPE_ETHERNET = 1

	// Hardware address length
	HLEN_ETHERNET = 6

	// Optionals
	OPT_PAD             = 0
	OPT_SUBNET_MASK     = 1
	OPT_ROUTER          = 3
	OPT_TIME_SERVER     = 4
	OPT_NAME_SERVER     = 5
	OPT_DOMAIN_SERVER   = 6
	OPT_HOSTNAME        = 12
	OPT_ADDR_REQUEST    = 50
	OPT_ADDR_LEASE_TIME = 51
	OPT_MSG_TYPE        = 53
	OPT_DHCP_SERVER_ID  = 54
	OPT_PARA_REQ        = 55
	OPT_RENEWAL_TIME    = 58
	OPT_REBINDING_TIME  = 59
	OPT_CLASS_ID        = 60
	OPT_CLIENT_ID       = 61
	OPT_RAPID_COMMIT    = 80
	OPT_END             = 255

	MSG_FIX_SIZE = int(unsafe.Sizeof(MessageFix{}))
)

const (
	// DHCP message type
	DHCPDISCOVER = iota + 1
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
)

var (
	_COOKIE = [4]byte{0x63, 0x82, 0x53, 0x63}

	_TransactionID uint32
	_Rand          *rand.Rand
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
