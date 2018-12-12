package ethernet

import (
	"encoding/binary"
	"errors"
	"net"
)

const SIZEOF_ETH_HEADER = 14

const (
	TYPE_IPV4 = 0x0800
)

var ETH_BROADCAST_ADDR = net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

type EthernetHeader struct {
	Src  net.HardwareAddr
	Dst  net.HardwareAddr
	Type uint16
}

func Unmarshal(b []byte) (ethHdr *EthernetHeader, err error) {
	if len(b) < SIZEOF_ETH_HEADER {
		return nil, errors.New("malformed packet")
	}

	ethHdr = &EthernetHeader{}
	ethHdr.Src = make([]byte, 6)
	copy(ethHdr.Src, b)
	ethHdr.Dst = make([]byte, 6)
	copy(ethHdr.Dst, b[6:])
	ethHdr.Type = binary.BigEndian.Uint16(b[12:])

	return
}

func (h *EthernetHeader) Marshal() (b []byte, err error) {
	b = make([]byte, 14, 14)
	copy(b, h.Dst[:6])
	copy(b[6:], h.Src[:6])
	binary.BigEndian.PutUint16(b[12:], h.Type)

	return
}
