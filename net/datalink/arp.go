package datalink

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"syscall"
	"time"

	"github.com/solomonwzs/goxutil/net/ethernet"
	"github.com/solomonwzs/goxutil/net/xnetutil"
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

func (r ArpRaw) SHA() net.HardwareAddr {
	hs := r.HardwareSize()
	return net.HardwareAddr(r[8 : 8+hs])
}

func (r ArpRaw) SPA() net.IP {
	hs := r.HardwareSize()
	ps := r.ProtocolSize()
	return net.IP(r[8+hs : 8+hs+ps])
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

	if len(a.SHA) != int(a.HardwareSize) ||
		len(a.THA) != int(a.HardwareSize) {
		return nil, errors.New("hardware size error")
	}

	var (
		spa net.IP
		tpa net.IP
	)
	if a.ProtocolSize == net.IPv4len {
		spa = a.SPA.To4()
		tpa = a.TPA.To4()
	} else {
		spa = a.SPA.To16()
		tpa = a.TPA.To16()
	}
	if spa == nil || tpa == nil {
		return nil, errors.New("protocol size error")
	}

	i := 8
	copy(b[i:], a.SHA[:a.HardwareSize])
	i += int(a.HardwareSize)
	copy(b[i:], spa)
	i += int(a.ProtocolSize)
	copy(b[i:], a.THA[:a.HardwareSize])
	i += int(a.HardwareSize)
	copy(b[i:], tpa)

	return
}

func recvArpReplyPacket(fd int, targetIP net.IP, res chan net.HardwareAddr,
	errChan chan error) {
	buf := make([]byte, 1024)
	arpRaw := ArpRaw(buf[14:])
	for {
		_, _, err := syscall.Recvfrom(fd, buf, 0)
		if err != nil {
			errChan <- err
			return
		}
		if arpRaw.Opcode() == ARP_OPC_REPLY &&
			xnetutil.BytesEqual(arpRaw.SPA(), targetIP) {
			select {
			case res <- arpRaw.SHA():
			default:
			}
			return
		}
	}
}

func broadcastArpRequest(interf *net.Interface, fd int, targetIP net.IP,
	errChan chan error) {
	var err error = nil
	defer func() {
		if err != nil {
			errChan <- err
		}
	}()

	addrs, err := interf.Addrs()
	if err != nil {
		return
	}

	ipLen := xnetutil.IPlen(targetIP)
	haLen := len(interf.HardwareAddr)
	p := [][]byte{}
	ethH := &ethernet.EthernetHeader{
		Src:  interf.HardwareAddr,
		Dst:  ethernet.ETH_BROADCAST_ADDR,
		Type: syscall.ETH_P_ARP,
	}
	tha := make([]uint8, haLen)
	for _, addr := range addrs {
		if spa, _, err := net.ParseCIDR(addr.String()); err != nil {
			return
		} else if xnetutil.IPlen(spa) == ipLen {
			arp := Arp{
				HardwareType: 1,
				ProtocolType: syscall.ETH_P_IP,
				HardwareSize: uint8(haLen),
				ProtocolSize: uint8(ipLen),
				Opcode:       ARP_OPC_REQUEST,
				SHA:          interf.HardwareAddr,
				SPA:          spa,
				THA:          tha,
				TPA:          targetIP,
			}
			p0, _ := ethH.Marshal()
			p1, _ := arp.Marshal()
			p0 = append(p0, p1...)
			p = append(p, p0)
		}

	}

	to := syscall.SockaddrLinklayer{
		Ifindex: interf.Index,
	}
	for {
		for _, p0 := range p {
			fmt.Println(len(p0))
			if err = syscall.Sendto(fd, p0, 0, &to); err != nil {
				return
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func GetHardwareAddr(interf *net.Interface, targetIP net.IP,
	timeout time.Duration) (
	hw net.HardwareAddr, err error) {
	var (
		timer   *time.Timer = nil
		res                 = make(chan net.HardwareAddr, 1)
		errChan             = make(chan error)
	)

	if timeout > 0 {
		timer = time.NewTimer(timeout)
		defer timer.Stop()
	}

	recvFd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW,
		int(xnetutil.Htons(syscall.ETH_P_ARP)))
	if err != nil {
		return
	}
	defer syscall.Close(recvFd)
	go recvArpReplyPacket(recvFd, targetIP, res, errChan)

	sendFd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW,
		int(xnetutil.Htons(syscall.ETH_P_ALL)))
	if err != nil {
		return
	}
	defer syscall.Close(sendFd)
	go broadcastArpRequest(interf, sendFd, targetIP, errChan)

	if timer != nil {
		select {
		case err = <-errChan:
		case hw = <-res:
		case <-timer.C:
			return nil, xnetutil.ERR_TIMEOUT
		}
	} else {
		select {
		case err = <-errChan:
		case hw = <-res:
		}
	}
	return
}
