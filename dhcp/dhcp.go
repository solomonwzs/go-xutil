package dhcp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math/rand"
	"net"
	"sync/atomic"
	"time"
	"unsafe"
)

var (
	ERR_NOT_EXIST     = errors.New("not exist")
	ERR_INVALID_VALUE = errors.New("invalid value")
)

type Message struct {
	fixBytes [MSG_FIX_SIZE]byte
	fix      *MessageFix
	opts     map[byte][]byte
}

func NewMessage() (msg *Message) {
	msg = &Message{
		opts: map[byte][]byte{},
	}
	msg.fix = (*MessageFix)(unsafe.Pointer(&msg.fixBytes[0]))
	binary.BigEndian.PutUint32(msg.fix.Xid[:],
		atomic.AddUint32(&_TransactionID, 1))

	return
}

func NewMessageForInterface(interf *net.Interface) (msg *Message) {
	msg = NewMessage()
	msg.fix.Op = BOOTREQUEST
	msg.fix.Htype = HTYPE_ETHERNET
	msg.fix.Hlen = HLEN_ETHERNET
	copy(msg.fix.Chaddr[:], interf.HardwareAddr)
	return
}

func Unmarshal(b []byte) (msg *Message, err error) {
	if len(b) < MSG_FIX_SIZE+4+1 {
		err = errors.New("DHCP malformed packet")
		return
	}

	msg = NewMessage()
	copy(msg.fixBytes[:], b)

	optsb := b[MSG_FIX_SIZE+4:]
	for i := 0; i < len(optsb) && optsb[i] != OPT_END; {
		opt := optsb[i]
		if opt == OPT_PAD || opt == OPT_RAPID_COMMIT {
			msg.opts[opt] = nil
			i += 1
		} else {
			l := int(optsb[i+1])
			msg.opts[opt] = optsb[i+2 : i+2+l]
			i += 2 + l
		}
	}

	return
}

func (msg *Message) Marshal() []byte {
	buf := bytes.NewBuffer(msg.fixBytes[:])
	buf.Write(_COOKIE[:])

	for opt, value := range msg.opts {
		buf.WriteByte(opt)
		l := byte(len(value))
		if l > 0 {
			buf.WriteByte(l)
			buf.Write(value)
		}
	}
	buf.WriteByte(OPT_END)

	return buf.Bytes()
}

func (msg *Message) TransactionID() uint32 {
	return binary.BigEndian.Uint32(msg.fix.Xid[:])
}

func (msg *Message) ClientIP() net.IP {
	return net.IP(msg.fix.Yiaddr[:])
}

func (msg *Message) SetBroadcast() {
	msg.fix.Flags[0] |= (1 << 7)
}

func (msg *Message) SetOptions(t byte, value []byte) {
	msg.opts[t] = value
}

func (msg *Message) SetMessageType(t byte) {
	msg.opts[OPT_MSG_TYPE] = []byte{t}
}

func (msg *Message) MessageType() (t byte, err error) {
	return msg.GetOptsByte(OPT_MSG_TYPE)
}

func (msg *Message) SetClientID(t byte, id []byte) {
	value := make([]byte, byte(len(id))+1)
	value[0] = t
	copy(value[1:], id)
	msg.opts[OPT_CLIENT_ID] = value
}

func (msg *Message) ClientID() (t byte, id []byte, err error) {
	return msg.GetTypeBytes(OPT_CLIENT_ID)
}

func (msg *Message) GetTypeBytes(opt byte) (t byte, b []byte, err error) {
	value, exist := msg.opts[opt]
	if !exist {
		err = ERR_NOT_EXIST
		return
	}
	if len(value) < 1 {
		err = ERR_INVALID_VALUE
		return
	}
	t = value[0]
	b = value[1:]
	return
}

func (msg *Message) GetOptsByte(opt byte) (i byte, err error) {
	value, exist := msg.opts[opt]
	if !exist {
		err = ERR_NOT_EXIST
		return
	}
	if len(value) != 1 {
		err = ERR_INVALID_VALUE
		return
	}
	i = value[0]
	return
}

func (msg *Message) GetOptsUint16(opt byte) (i uint16, err error) {
	value, exist := msg.opts[opt]
	if !exist {
		err = ERR_NOT_EXIST
		return
	}
	if len(value) != 2 {
		err = ERR_INVALID_VALUE
		return
	}
	i = binary.BigEndian.Uint16(value)
	return
}

func (msg *Message) GetOptsUint32(opt byte) (i uint32, err error) {
	value, exist := msg.opts[opt]
	if !exist {
		err = ERR_NOT_EXIST
		return
	}
	if len(value) != 4 {
		err = ERR_INVALID_VALUE
		return
	}
	i = binary.BigEndian.Uint32(value)
	return
}

func (msg *Message) GetOptsIpv4Addr(opt byte) (addr net.IP, err error) {
	value, exist := msg.opts[opt]
	if !exist {
		err = ERR_NOT_EXIST
		return
	}
	if len(value) != 4 {
		err = ERR_INVALID_VALUE
		return
	}
	addr = net.IP(value)
	return
}

func (msg *Message) GetOptsIpv4Addrs(opt byte) (addrs []net.IP, err error) {
	value, exist := msg.opts[opt]
	if !exist {
		err = ERR_NOT_EXIST
		return
	}
	if len(value)%4 != 0 {
		err = ERR_INVALID_VALUE
		return
	}

	addrs = make([]net.IP, len(value)/4)
	for i := 0; i < len(addrs); i++ {
		j := i * 4
		addrs[i] = net.IP(value[j : j+4])
	}

	return
}

func (msg *Message) Router() (addr net.IP, err error) {
	return msg.GetOptsIpv4Addr(OPT_ROUTER)
}

func (msg *Message) DHCPServerID() (addr net.IP, err error) {
	return msg.GetOptsIpv4Addr(OPT_DHCP_SERVER_ID)
}

func (msg *Message) SubnetMask() (mask net.IP, err error) {
	return msg.GetOptsIpv4Addr(OPT_SUBNET_MASK)
}

func (msg *Message) DomainServers() (addrs []net.IP, err error) {
	return msg.GetOptsIpv4Addrs(OPT_DOMAIN_SERVER)
}

func (msg *Message) RebindingTime() (t uint32, err error) {
	return msg.GetOptsUint32(OPT_REBINDING_TIME)
}

func (msg *Message) RenewalTime() (t uint32, err error) {
	return msg.GetOptsUint32(OPT_RENEWAL_TIME)
}

func (msg *Message) AddressLeaseTime() (t uint32, err error) {
	return msg.GetOptsUint32(OPT_ADDR_LEASE_TIME)
}

func init() {
	_Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	_TransactionID = _Rand.Uint32()
}
