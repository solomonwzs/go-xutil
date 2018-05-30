package datalink

import (
	"net"
	"syscall"

	"github.com/solomonwzs/goxutil/net/ethernet"
	"github.com/solomonwzs/goxutil/net/xnetutil"
)

type DlSocket struct {
	fd        int
	ethP      []byte
	to        *syscall.SockaddrLinklayer
	sendFlags int
	readFlags int
}

func NewDlSocket(dev string, ethType uint16) (sock *DlSocket, err error) {
	gateway, err := xnetutil.GetGateway(dev)
	if err != nil {
		return
	}

	hardwareAddr, err := xnetutil.GetHardwareAddr(dev, gateway)
	if err == xnetutil.ERR_NOT_FOUND {
		hardwareAddr, err = GetHardwareAddr(dev, gateway, 0)
	}
	if err != nil {
		return
	}

	interf, err := net.InterfaceByName(dev)
	if err != nil {
		return
	}

	ethH := &ethernet.EthernetHeader{
		Src:  interf.HardwareAddr,
		Dst:  hardwareAddr,
		Type: ethType,
	}
	ethP, err := ethH.Marshal()
	if err != nil {
		return
	}

	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW,
		int(xnetutil.Htons(syscall.ETH_P_ALL)))
	if err != nil {
		return
	}

	sock = &DlSocket{
		fd:        fd,
		ethP:      ethP,
		sendFlags: 0,
		readFlags: 0,
		to: &syscall.SockaddrLinklayer{
			Ifindex: interf.Index,
		},
	}

	return
}

func (sock *DlSocket) Close() error {
	return syscall.Close(sock.fd)
}

func (sock *DlSocket) Write(p []byte) (n int, err error) {
	buf := make([]byte, len(sock.ethP)+len(p))
	copy(buf, sock.ethP)
	copy(buf[len(sock.ethP):], p)

	err = syscall.Sendto(sock.fd, buf, sock.sendFlags, sock.to)
	if err != nil {
		return
	}

	return len(p), err
}

func (sock *DlSocket) Read(p []byte) (n int, err error) {
	n, _, err = syscall.Recvfrom(sock.fd, p, sock.readFlags)
	return
}
