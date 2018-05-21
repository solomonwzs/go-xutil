package ethernet

import "encoding/binary"

type EthernetHeader struct {
	Src  []uint8
	Dst  []uint8
	Type uint16
}

func (h *EthernetHeader) Marshal() (b []byte, err error) {
	b = make([]byte, 14, 14)
	copy(b, h.Dst[:6])
	copy(b[6:], h.Src[:6])
	binary.BigEndian.PutUint16(b[12:], h.Type)

	return
}
