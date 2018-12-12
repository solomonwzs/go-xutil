package transport

import (
	"encoding/binary"
	"errors"
	"net"
	"syscall"

	"github.com/solomonwzs/goxutil/net/ethernet"
	"github.com/solomonwzs/goxutil/net/network"
	"github.com/solomonwzs/goxutil/net/xnetutil"
)

const SIZEOF_UDP_HEADER = 8

type UDP struct {
	IPHeader *network.IPv4Header
	SrcPort  uint16
	DstPort  uint16
	Length   uint16
	Checksum uint16
	Data     []byte
}

func RawUDPUnmarshal(raw []byte) (
	ethHdr *ethernet.EthernetHeader, u *UDP, err error) {
	if len(raw) < ethernet.SIZEOF_ETH_HEADER+network.SIZEOF_IPV4_HEADER {
		return nil, nil, errors.New("malformed packet")
	}

	if ethHdr, err = ethernet.Unmarshal(raw); err != nil ||
		ethHdr.Type != ethernet.TYPE_IPV4 {
		return nil, nil, errors.New("malformed packet")
	}

	ipHdr, err := network.IPv4HeaderUnmarshal(
		raw[ethernet.SIZEOF_ETH_HEADER:])
	if err != nil || ipHdr.Protocol != network.IPV4_PRO_UDP {
		return nil, nil, errors.New("malformed packet")
	}

	udpRaw := raw[ethernet.SIZEOF_ETH_HEADER+ipHdr.IHL:]
	if len(udpRaw) < int(ipHdr.Length)-int(ipHdr.IHL) {
		return nil, nil, errors.New("malformed packet")
	}

	u = &UDP{IPHeader: ipHdr}
	u.SrcPort = binary.BigEndian.Uint16(udpRaw)
	u.DstPort = binary.BigEndian.Uint16(udpRaw[2:])
	u.Length = binary.BigEndian.Uint16(udpRaw[4:])
	u.Checksum = binary.BigEndian.Uint16(udpRaw[6:])
	if len(udpRaw[8:]) < int(u.Length) {
		return nil, nil, errors.New("malformed packet")
	}
	u.Data = udpRaw[8 : 8+u.Length]

	return
}

func (u *UDP) Marshal() (b []byte, err error) {
	if u.IPHeader == nil {
		return nil, errors.New("[udp] miss ip header")
	}
	if u.Length == 0 {
		u.Length = uint16(SIZEOF_UDP_HEADER + len(u.Data))
	}

	s := xnetutil.NewChecksumer()
	s.Write(u.IPHeader.SrcAddr[:4])
	s.Write(u.IPHeader.DstAddr[:4])
	s.Write([]byte{
		0,
		u.IPHeader.Protocol,
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

func NewBroadcastUDPRaw(interf *net.Interface, srcPort uint16,
	dstPort uint16, data []byte) (raw []byte, err error) {
	ethHdr := &ethernet.EthernetHeader{
		Src:  interf.HardwareAddr,
		Dst:  ethernet.ETH_BROADCAST_ADDR,
		Type: syscall.ETH_P_IP,
	}

	ipHdr := &network.IPv4Header{
		Version:    4,
		TOS:        0,
		Id:         network.NewIPHdrID(),
		Flags:      network.IPV4_FLAG_DONT_FRAG,
		FragOffset: 0,
		TTL:        64,
		Protocol:   syscall.IPPROTO_UDP,
		SrcAddr:    net.IP{0, 0, 0, 0},
		DstAddr:    net.IP{0xff, 0xff, 0xff, 0xff},
	}

	u := &UDP{
		IPHeader: ipHdr,
		SrcPort:  srcPort,
		DstPort:  dstPort,
		Data:     data,
	}

	pEth, err := ethHdr.Marshal()
	if err != nil {
		return
	}

	pUDP, err := u.Marshal()
	if err != nil {
		return
	}
	ipHdr.Length = network.SIZEOF_IPV4_HEADER + uint16(len(pUDP))

	pIP, err := ipHdr.Marshal()
	if err != nil {
		return
	}

	raw = make([]byte, len(pEth)+len(pIP)+len(pUDP))
	copy(raw, pEth)
	copy(raw[len(pEth):], pIP)
	copy(raw[len(pEth)+len(pIP):], pUDP)

	return
}
