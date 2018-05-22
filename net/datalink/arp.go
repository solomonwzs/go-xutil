package arp

import (
	"encoding/binary"
	"net"
	"syscall"
	"time"

	"github.com/solomonwzs/goxutil/net/util"
)

type ArpRaw []byte

func (r ArpRaw) HardwareSize() uint8 {
	return r[4]
}

func (r ArpRaw) ProtocolSize() uint8 {
	return r[5]
}

func (r ArpRaw) Opcode() uint16 {
	return binary.BigEndian.Uint16(r[6:])
}

func (r ArpRaw) THA() net.HardwareAddr {
	hs := r.HardwareSize()
	ps := r.ProtocolSize()
	i := 8 + hs + ps
	return net.HardwareAddr(r[i : i+hs])
}

func (r ArpRaw) TPA() net.IP {
	hs := r.HardwareSize()
	ps := r.ProtocolSize()
	i := 8 + hs + ps + hs
	return net.IP(r[i : i+ps])
}

type Arp struct {
	HardwareType uint16
	ProtocolType uint16
	HardwareSize uint8
	ProtocolSize uint8
	Opcode       uint16
	SHA          net.HardwareAddr
	SPA          net.IP
	THA          net.HardwareAddr
	TPA          net.IP
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

func recvArpReplyPacket(fd int, targetIP net.IP, res chan net.HardwareAddr) {
	buf := make([]byte, 1024)
	arpRaw := ArpRaw(buf[14:])
	reply := util.Htons(ARP_OPC_REPLY)
	for {
		_, _, err := syscall.Recvfrom(fd, buf, 0)
		if err != nil {
			return
		}
		if arpRaw.Opcode() == reply &&
			util.BytesEqual(arpRaw.TPA(), targetIP) {
			select {
			case res <- arpRaw.THA():
			default:
			}
			return
		}
	}
}

func broadcastArpRequest(fd int) {
}

func GetHardwareAddr(dev string, ip net.IP, timeout time.Duration) (
	hw net.HardwareAddr, err error) {
	var (
		timer *time.Timer = nil
		end               = make(chan struct{})
		res               = make(chan net.HardwareAddr, 1)
	)
	defer close(end)

	if timeout > 0 {
		timer = time.NewTimer(timeout)
		defer timer.Stop()
	}

	recvFd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW,
		int(util.Htons(syscall.ETH_P_ARP)))
	if err != nil {
		return
	}
	defer syscall.Close(recvFd)
	go recvArpReplyPacket(recvFd, ip, res)

	sendFd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW,
		int(util.Htons(syscall.ETH_P_ALL)))
	if err != nil {
		return
	}
	defer syscall.Close(sendFd)

	select {
	case <-timer.C:
		return nil, util.ERR_TIMEOUT
	}
	return
}
