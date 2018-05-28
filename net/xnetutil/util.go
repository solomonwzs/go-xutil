package xnetutil

import "net"

type TimeoutError struct{}

func (e *TimeoutError) Error() string   { return "[net] i/o timeout" }
func (e *TimeoutError) Timeout() bool   { return true }
func (e *TimeoutError) Temporary() bool { return true }

var (
	ERR_TIMEOUT = &TimeoutError{}
)

func Checksum(raw []byte) uint16 {
	sum := uint32(0)
	i := 0
	size := len(raw)
	for ; i < size-1; i += 2 {
		sum += uint32(raw[i]) + uint32(raw[i+1])<<8
	}
	if i != size {
		sum += uint32(raw[size-1])
	}

	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	return uint16(^sum)
}

func BytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func Htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}

func Ntohs(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}

func IPlen(ip net.IP) int {
	if ipv4 := ip.To4(); ipv4 != nil {
		return net.IPv4len
	} else {
		return net.IPv6len
	}
}
