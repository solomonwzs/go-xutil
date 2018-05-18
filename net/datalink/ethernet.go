package datalink

type EthernetHeader struct {
	Dst  [6]uint8
	Src  [6]uint8
	Type uint32
}
