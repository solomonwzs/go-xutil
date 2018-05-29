package xnetutil

import (
	"net"
	"unsafe"
)

type TimeoutError struct{}

func (e *TimeoutError) Error() string   { return "[net] i/o timeout" }
func (e *TimeoutError) Timeout() bool   { return true }
func (e *TimeoutError) Temporary() bool { return true }

var (
	ERR_TIMEOUT = &TimeoutError{}
)

type Checksumer struct {
	sum uint32
	i   int
}

func checksum(p []byte, sum uint32, i int) (uint32, int) {
	for _, b := range p {
		if i&0x1 == 1 {
			sum += uint32(b) << 8
		} else {
			sum += uint32(b)
		}
		i += 1
	}
	return sum, i
}

func NewChecksumer() *Checksumer {
	return &Checksumer{
		sum: 0,
		i:   0,
	}
}

func (c *Checksumer) Write(p []byte) (n int, err error) {
	c.sum, c.i = checksum(p, c.sum, c.i)
	return len(p), nil
}

func (c *Checksumer) SumU16(p []byte) uint16 {
	sum, _ := checksum(p, c.sum, c.i)
	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	return uint16(^sum)
}

func (c *Checksumer) Sum(p []byte) []byte {
	s := []byte{0, 0}
	sum := (*uint16)(unsafe.Pointer(&s[0]))
	*sum = c.SumU16(p)
	return s
}

func (c *Checksumer) Reset() {
	c.sum = 0
	c.i = 0
}

func (c *Checksumer) Size() int {
	return 2
}

func (c *Checksumer) BlockSize() int {
	return 64
}

func Checksum(p []byte) uint16 {
	sum, _ := checksum(p, 0, 0)
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
