package transport

import (
	"errors"

	"github.com/solomonwzs/goxutil/net/network"
	"github.com/solomonwzs/goxutil/net/xnetutil"
)

type Udp struct {
	IpH      *network.IPv4Header
	SrcPort  uint16
	DstPort  uint16
	Length   uint16
	Checksum uint16
	Data     []byte
}

func (u *Udp) Marshal() (b []byte, err error) {
	if u.IpH == nil {
		return nil, errors.New("miss ip header")
	}

	s := xnetutil.NewChecksumer()
	s.Write(u.IpH.SrcAddr[:4])
	s.Write(u.IpH.DstAddr[:4])

	return
}
