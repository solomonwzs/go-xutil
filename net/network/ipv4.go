package network

import (
	"encoding/binary"
	"errors"
	"net"
	"unsafe"
)

type IPv4Header struct {
	Version    uint8
	IHL        uint8
	TOS        uint8
	Length     uint16
	Id         uint16
	Flags      uint16
	FragOffset uint16
	TTL        uint8
	Protocol   uint8
	checksum   uint16
	SrcAddr    net.IP
	DstAddr    net.IP
	Options    []byte
}

func (h *IPv4Header) Marshal() (b []byte, err error) {
	if len(h.Options) == 0 {
		b = make([]byte, SIZEOF_IPV4_HEADER)
	} else {
		b = make([]byte, SIZEOF_IPV4_HEADER+4)
		copy(b[SIZEOF_IPV4_HEADER:], h.Options)
	}
	h.IHL = uint8(len(b))

	b[0] = h.Version<<4 | (h.IHL >> 2 & 0x0f)
	b[1] = h.TOS
	binary.BigEndian.PutUint16(b[2:], h.Length)
	binary.BigEndian.PutUint16(b[4:], h.Id)
	binary.BigEndian.PutUint16(b[6:],
		h.Flags<<13|(h.FragOffset&0x1fff))
	b[8] = h.TTL
	b[9] = h.Protocol
	checksum := (*uint16)(unsafe.Pointer(&b[10]))

	if ip := h.SrcAddr.To4(); ip != nil {
		copy(b[12:16], ip[:net.IPv4len])
	}
	if ip := h.DstAddr.To4(); ip != nil {
		copy(b[16:20], ip[:net.IPv4len])
	} else {
		return nil, errors.New("[ipv4] missing address")
	}

	*checksum = 0
	*checksum = Checksum(b)
	return
}
