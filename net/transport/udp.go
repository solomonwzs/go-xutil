package transport

import (
	"encoding/binary"
	"errors"

	"github.com/solomonwzs/goxutil/net/network"
	"github.com/solomonwzs/goxutil/net/xnetutil"
)

const SIZEOF_UDP_HEADER = 8

type UDP struct {
	IpH      *network.IPv4Header
	SrcPort  uint16
	DstPort  uint16
	Length   uint16
	Checksum uint16
	Data     []byte
}

func (u *UDP) Marshal() (b []byte, err error) {
	if u.IpH == nil {
		return nil, errors.New("[udp] miss ip header")
	}
	if u.Length == 0 {
		u.Length = uint16(SIZEOF_UDP_HEADER + len(u.Data))
	}

	s := xnetutil.NewChecksumer()
	s.Write(u.IpH.SrcAddr[:4])
	s.Write(u.IpH.DstAddr[:4])
	s.Write([]byte{
		0,
		u.IpH.Protocol,
		byte(u.Length >> 8), byte(u.Length & 0xff),
		byte(u.SrcPort >> 8), byte(u.SrcPort & 0xff),
		byte(u.DstPort >> 8), byte(u.DstPort & 0xff),
		byte(u.Length >> 8), byte(u.Length & 0xff),
		0, 0,
	})
	s.Write(u.Data)
	u.Checksum = s.SumU16(nil)

	b = make([]byte, u.Length)
	binary.BigEndian.PutUint16(b, u.SrcPort)
	binary.BigEndian.PutUint16(b[2:], u.DstPort)
	binary.BigEndian.PutUint16(b[4:], u.Length)
	binary.BigEndian.PutUint16(b[6:], u.Checksum)
	copy(b[8:], u.Data)

	return
}
