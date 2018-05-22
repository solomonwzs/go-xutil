package ethernet

import (
	"encoding/binary"
	"net"
)

var ETH_BROADCAST_ADDR = net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

type EthernetHeader struct {
	Src  net.HardwareAddr
	Dst  net.HardwareAddr
	Type uint16
}

func (h *EthernetHeader) Marshal() (b []byte, err error) {
	b = make([]byte, 14, 14)
	copy(b, h.Dst[:6])
	copy(b[6:], h.Src[:6])
	binary.BigEndian.PutUint16(b[12:], h.Type)

	return
}
