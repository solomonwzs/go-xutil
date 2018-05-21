package arp

import "encoding/binary"

type Arp struct {
	HardwareType uint16
	ProtocolType uint16
	HardwareSize uint8
	ProtocolSize uint8
	Opcode       uint16
	SHA          []uint8
	SPA          []uint8
	THA          []uint8
	TPA          []uint8
}

func (a *Arp) Marshal() (b []byte, err error) {
	size := 8 + int(a.HardwareSize)*2 + int(a.ProtocolSize)*2
	b = make([]byte, size, size)

	binary.BigEndian.PutUint16(b[0:], a.HardwareType)
	binary.BigEndian.PutUint16(b[2:], a.ProtocolType)
	b[4] = a.HardwareSize
	b[5] = a.ProtocolSize
	binary.BigEndian.PutUint16(b[6:], a.Opcode)

	i := 8
	copy(b[i:], a.SHA[:a.HardwareSize])
	i += int(a.HardwareSize)
	copy(b[i:], a.SPA[:a.ProtocolSize])
	i += int(a.ProtocolSize)
	copy(b[i:], a.THA[:a.HardwareSize])
	i += int(a.HardwareSize)
	copy(b[i:], a.TPA[:a.ProtocolSize])

	return
}
