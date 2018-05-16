package network

import (
	"errors"
	"unsafe"
)

type IcmpHeader struct {
	typ      uint8
	code     uint8
	checksum uint16
}

func BuildIcmpHeader(raw []byte, typ uint8, code uint8) (err error) {
	if len(raw) < 4 {
		return errors.New("[icmp] length of raw less than 4")
	}

	header := (*IcmpHeader)(unsafe.Pointer(&raw[0]))
	header.typ = typ
	header.code = code
	header.checksum = 0
	header.checksum = Checksum(raw)

	return
}
