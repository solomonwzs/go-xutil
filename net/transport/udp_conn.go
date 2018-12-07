package transport

import (
	"syscall"
)

const (
	MSG_PEEK     = syscall.MSG_PEEK
	MSG_OOB      = syscall.MSG_OOB
	MSG_WAITALL  = syscall.MSG_WAITALL
	MSG_EOR      = syscall.MSG_EOR
	MSG_NOSIGNAL = syscall.MSG_NOSIGNAL
)

type UDPBroadcastConn struct {
	fd            int
	broadcastAddr *syscall.SockaddrInet4
	sendFlags     int
	recvFlags     int
}

func NewUDPBroadcastConn(srcPort int, dstPort int) (
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
		Port: srcPort,
		Addr: [4]byte{0, 0, 0, 0},
	}
	if err = syscall.Bind(fd, bindAddr); err != nil {
		return
	}

	broadcastAddr := &syscall.SockaddrInet4{
		Port: dstPort,
		Addr: [4]byte{255, 255, 255, 255},
	}
	conn = &UDPBroadcastConn{
		fd:            fd,
		broadcastAddr: broadcastAddr,
		sendFlags:     0,
		recvFlags:     0,
	}

	return
}

func (conn *UDPBroadcastConn) SetSendFlags(flags int) {
	conn.sendFlags = flags
}

func (conn *UDPBroadcastConn) SetRecvFlags(flags int) {
	conn.recvFlags = flags
}

func (conn *UDPBroadcastConn) Fd() int {
	return conn.fd
}

func (conn *UDPBroadcastConn) Close() error {
	return syscall.Close(conn.fd)
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
