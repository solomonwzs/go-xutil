package network

import "fmt"

type NetworkData interface {
	Marshal() ([]byte, error)
}

func Checksum(raw []byte) uint16 {
	sum := uint32(0)
	i := 0
	size := len(raw)
	for ; i < size-1; i += 2 {
		sum += uint32(raw[i]) + uint32(raw[i+1])<<8
	}
	if i != size {
		sum += uint32(raw[size-1])
	}

	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	return uint16(^sum)
}

func PrintHex(p []byte) {
	for _, b := range p {
		fmt.Printf("%02x ", b)
	}
	fmt.Printf("\n")
}
