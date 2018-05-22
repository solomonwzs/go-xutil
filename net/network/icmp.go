package network

import (
	"encoding/binary"
	"unsafe"

	"github.com/solomonwzs/goxutil/net/ethernet"
	"github.com/solomonwzs/goxutil/net/util"
)

type Icmp struct {
	Type     uint8
	Code     uint8
	checksum uint16
	Data     ethernet.NetworkData
}

func (icmp *Icmp) Marshal() (b []byte, err error) {
	if icmp.Data == nil {
		b = make([]byte, SIZEOF_ICMP_HEADER)
	} else if b0, err0 := icmp.Data.Marshal(); err0 != nil {
		return b, err0
	} else {
		b = make([]byte, int(SIZEOF_ICMP_HEADER)+len(b0))
		copy(b[SIZEOF_ICMP_HEADER:], b0)
	}
	b[0] = icmp.Type
	b[1] = icmp.Code
	checksum := (*uint16)(unsafe.Pointer(&b[2]))

	*checksum = 0
	*checksum = util.Checksum(b)
	return
}

type IcmpEcho struct {
	Id     uint16
	SeqNum uint16
}

func (e *IcmpEcho) Marshal() (b []byte, err error) {
	b = make([]byte, SIZEOF_ICMP_ECHO)

	binary.BigEndian.PutUint16(b[0:], e.Id)
	binary.BigEndian.PutUint16(b[2:], e.SeqNum)

	return
}
