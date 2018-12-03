package dhcp

import (
	"math/rand"
	"net"
	"time"
)

func init() {
	_Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	_TransactionID = _Rand.Uint32()

	_SrcAddr = &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: SRC_PORT}
	_DstAddr = &net.UDPAddr{IP: net.IPv4(255, 255, 255, 255), Port: DST_PORT}
}

type Client struct {
}
