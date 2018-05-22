package arp

import (
	"encoding/binary"
	"net"
	"syscall"
	"time"
	"unsafe"

	"github.com/solomonwzs/goxutil/net/util"
)

type ArpRaw []byte

func (r ArpRaw) THA() net.HardwareAddr {
	return net.HardwareAddr(r)
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

func recvArpReplyPacket(fd int, targetIP net.IP, res chan []byte) {
	buf := make([]byte, 1024)
	arpRaw := ArpRaw(buf[14:])
	opcode := (*uint16)(unsafe.Pointer(&arpRaw[6]))
	for {
		_, _, err := syscall.Recvfrom(fd, buf, 0)
		if err != nil {
			return
		}
		if util.Ntohs(*opcode) == ARP_OPC_REPLY {
		}
	}
}

func GetHardwareAddr(dev string, ip net.IP, timeout time.Duration) (
	hw net.HardwareAddr, err error) {
	var (
		timer *time.Timer = nil
		end               = make(chan struct{})
		// res               = make(chan []byte, 1)
	)
	defer close(end)

	if timeout > 0 {
		timer = time.NewTimer(timeout)
		defer timer.Stop()
	}

	select {
	case <-timer.C:
		return nil, util.ERR_TIMEOUT
	}
	return
}
