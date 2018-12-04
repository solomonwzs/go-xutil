package dhcp

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"net"
	"sync/atomic"
	"time"
	"unsafe"
)

type Message struct {
	fixBytes [MSG_FIX_SIZE]byte
	fix      *MessageFix
	cookie   [4]byte
	opts     *bytes.Buffer
}

func NewMessage() (msg *Message) {
	msg = &Message{
		opts: &bytes.Buffer{},
	}
	msg.fix = (*MessageFix)(unsafe.Pointer(&msg.fixBytes[0]))
	binary.BigEndian.PutUint32(msg.fix.Xid[:],
		atomic.AddUint32(&_TransactionID, 1))

	return
}

func (msg *Message) Marshal() []byte {
	buf := bytes.NewBuffer(msg.fixBytes[:])
	buf.Write(msg.cookie[:])
	buf.Write(msg.opts.Bytes())
	buf.WriteByte(OPT_END)

	return buf.Bytes()
}

func (msg *Message) SetCookie(cookie []byte) {
	copy(msg.cookie[:], cookie)
}

func (msg *Message) SetMassageType(t byte) {
	msg.opts.Write([]byte{OPT_MSG_TYPE, 1, t})
}

func (msg *Message) SetClientID(t byte, id []byte) {
	msg.opts.Write([]byte{OPT_CLIENT_ID, byte(len(id)) + 1, t})
	msg.opts.Write(id)
}

func init() {
	_Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	_TransactionID = _Rand.Uint32()

	_SrcAddr = &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: SRC_PORT}
	_DstAddr = &net.UDPAddr{IP: net.IPv4(255, 255, 255, 255), Port: DST_PORT}
}

type Client struct {
}
