package transport

import (
	"net"
	"syscall"

	"github.com/solomonwzs/goxutil/net/ethernet"
	"github.com/solomonwzs/goxutil/net/xnetutil"
)

const (
	MSG_PEEK     = syscall.MSG_PEEK
	MSG_OOB      = syscall.MSG_OOB
	MSG_WAITALL  = syscall.MSG_WAITALL
	MSG_EOR      = syscall.MSG_EOR
	MSG_NOSIGNAL = syscall.MSG_NOSIGNAL
)

type broadcastConn struct {
	fd        int
	sendFlags int
	recvFlags int
}

func (conn *broadcastConn) SetSendFlags(flags int) {
	conn.sendFlags = flags
}

func (conn *broadcastConn) SetRecvFlags(flags int) {
	conn.recvFlags = flags
}

func (conn *broadcastConn) Fd() int {
	return conn.fd
}

func (conn *broadcastConn) Close() error {
	return syscall.Close(conn.fd)
}

type UDPBroadcastConn struct {
	broadcastConn
	broadcastAddr *syscall.SockaddrInet4
}

func NewUDPBroadcastConn(srcPort uint16, dstPort uint16) (
	conn *UDPBroadcastConn, err error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM,
		syscall.IPPROTO_UDP)
	if err != nil {
		return
	}

	if err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET,
		syscall.SO_BROADCAST, 1); err != nil {
		return
	}

	bindAddr := &syscall.SockaddrInet4{
		Port: int(srcPort),
		Addr: [4]byte{0, 0, 0, 0},
	}
	if err = syscall.Bind(fd, bindAddr); err != nil {
		return
	}

	broadcastAddr := &syscall.SockaddrInet4{
		Port: int(dstPort),
		Addr: [4]byte{255, 255, 255, 255},
	}
	conn = &UDPBroadcastConn{
		broadcastConn: broadcastConn{
			fd:        fd,
			sendFlags: 0,
			recvFlags: 0,
		},
		broadcastAddr: broadcastAddr,
	}

	return
}

func (conn *UDPBroadcastConn) Write(p []byte) (n int, err error) {
	err = syscall.Sendto(conn.fd, p, conn.sendFlags, conn.broadcastAddr)
	return len(p), err
}

func (conn *UDPBroadcastConn) Read(p []byte) (n int, err error) {
	n, _, err = syscall.Recvfrom(conn.fd, p, conn.recvFlags)
	return
}

func (conn *UDPBroadcastConn) Readfrom(p []byte) (
	n int, from *syscall.SockaddrInet4, err error) {
	n, from0, err := syscall.Recvfrom(conn.fd, p, conn.recvFlags)
	if err == nil {
		from = from0.(*syscall.SockaddrInet4)
	}
	return
}

type UDPBroadcastRawConn struct {
	broadcastConn
	broadcastAddr *syscall.SockaddrLinklayer
	interf        *net.Interface
	srcPort       uint16
	dstPort       uint16

	hFrom *syscall.SockaddrLinklayer
	iFrom *syscall.SockaddrInet4
	cache []byte
}

func NewUDPBroadcastRawConn(interf *net.Interface, srcPort, dstPort uint16) (
	conn *UDPBroadcastRawConn, err error) {
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW,
		int(xnetutil.Htons(syscall.ETH_P_IP)))
	if err != nil {
		return
	}

	addr := &syscall.SockaddrLinklayer{
		Protocol: xnetutil.Htons(syscall.ETH_P_IP),
		Ifindex:  interf.Index,
		Halen:    6,
		Addr:     [8]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
	}

	conn = &UDPBroadcastRawConn{
		broadcastConn: broadcastConn{
			fd:        fd,
			sendFlags: 0,
			recvFlags: 0,
		},
		broadcastAddr: addr,
		interf:        interf,
		srcPort:       srcPort,
		dstPort:       dstPort,
	}
	return
}

func (conn *UDPBroadcastRawConn) Write(p []byte) (n int, err error) {
	raw, err := NewBroadcastUDPRaw(conn.interf, conn.srcPort,
		conn.dstPort, p)
	if err != nil {
		return
	}
	err = syscall.Sendto(conn.fd, raw, conn.sendFlags, conn.broadcastAddr)
	return len(p), err
}

func (conn *UDPBroadcastRawConn) ReadFrom(p []byte) (
	n int, hFrom *syscall.SockaddrLinklayer, iFrom *syscall.SockaddrInet4,
	err error) {
	if len(conn.cache) == 0 {
		buf := make([]byte, 0xffff)
		for {
			n0, from, err0 := syscall.Recvfrom(conn.fd, buf, conn.recvFlags)
			if err0 != nil {
				err = err0
				return
			}

			if n0 < ethernet.SIZEOF_ETH_HEADER {
				continue
			}

			_, u, err0 := RawUDPUnmarshal(buf)
			if err0 != nil || u.DstPort != conn.srcPort ||
				u.SrcPort != conn.dstPort {
				continue
			}

			conn.hFrom = from.(*syscall.SockaddrLinklayer)
			conn.iFrom = &syscall.SockaddrInet4{Port: int(conn.dstPort)}
			copy(conn.iFrom.Addr[:], u.IPHeader.SrcAddr)
			conn.cache = u.Data
			break
		}
	}
	n = copy(p, conn.cache)
	hFrom = conn.hFrom
	iFrom = conn.iFrom
	conn.cache = conn.cache[n:]
	return
}

func (conn *UDPBroadcastRawConn) Read(p []byte) (n int, err error) {
	n, _, _, err = conn.ReadFrom(p)
	return
}
