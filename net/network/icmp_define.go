package network

// Control messages type
const (
	ICMP_CT_ECHO_REPLY      = 0
	ICMP_CT_DST_UNREACHABLE = 3
	ICMP_CT_ECHO_REQUEST    = 8
	ICMP_CT_TIME_EXCEEDED   = 11
)

// Control messages code
const (
	// Destination Unreachable
	ICMP_CC_DST_NETWORK_UR     = 0
	ICMP_CC_DST_HOST_UR        = 1
	ICMP_CC_DST_PROTOCOL_UR    = 2
	ICMP_CC_DST_PORT_UR        = 3
	ICMP_CC_FRAG_REQ_DF_SET    = 4
	ICMP_CC_SRC_ROUTE_FAILED   = 5
	ICMP_CC_DST_NETWORK_UK     = 6
	ICMP_CC_DST_HOST_UK        = 7
	ICMP_CC_SRC_HOST_ISO       = 8
	ICMP_CC_NETWORK_PRO        = 9
	ICMP_CC_HOST_PRO           = 10
	ICMP_CC_DST_NETWORK_UR_TOS = 11
	ICMP_CC_DST_HOST_UR_TOS    = 12
	ICMP_CC_COMM_PRO           = 13
	ICMP_CC_HOST_PRE_VIOLATION = 14
	ICMP_CC_PRE_CUTOFF         = 15
)

const (
	SIZEOF_ICMP_HEADER = 4
	SIZEOF_ICMP_ECHO   = 4
)
