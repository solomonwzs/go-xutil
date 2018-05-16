package network

type IPv4Header struct {
	version_ihl   uint8
	TOS           uint8
	Length        uint16
	Id            uint16
	flags_fOffset uint16
	TLL           uint8
	Protocol      uint8
	Checksum      uint16
	SrcAddr       uint32
	DstAddr       uint32
}

func (ip *IPv4Header) Version() uint8 {
	return ip.version_ihl & 0x0f
}

func (ip *IPv4Header) SetVersion(version uint8) {
	ip.version_ihl = (ip.version_ihl & 0xf0) | (version & 0x0f)
}

func (ip *IPv4Header) IHL() uint8 {
	return ip.version_ihl >> 4
}

func (ip *IPv4Header) SetIHL(ihl uint8) {
	ip.version_ihl = (ip.version_ihl & 0x0f) | (ihl << 4)
}

func (ip *IPv4Header) Flags() uint16 {
	return ip.flags_fOffset & 0x0007
}

func (ip *IPv4Header) SetFlags(flags uint16) {
	ip.flags_fOffset = (ip.flags_fOffset & 0xfff8) | (flags & 0x0007)
}

func (ip *IPv4Header) FragmentOffset() uint16 {
	return ip.flags_fOffset >> 3
}

func (ip *IPv4Header) SetFragmentOffset(offset uint16) {
	ip.flags_fOffset = (ip.flags_fOffset & 0x0007) | (offset << 3)
}
