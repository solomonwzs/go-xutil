package icmp

import (
	"errors"
	"unsafe"
)

type IcmpHeader struct {
	Type     uint8
	Code     uint8
	checksum uint16
}

func BuildIcmpHeader(raw []byte, typ uint8, code uint8) (err error) {
	if len(raw) < 4 {
		return errors.New("[icmp] length of raw less than 4")
	}

	header := (*IcmpHeader)(unsafe.Pointer(&raw[0]))
	header.Type = typ
	header.Code = code

	return
}

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
