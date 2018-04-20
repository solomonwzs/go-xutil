package pubsub

import (
	"errors"
	"time"

	"github.com/solomonwzs/goxutil/closer"
)

var (
	ErrTimeout = errors.New("timeout")
	ErrAgain   = errors.New("resource temporarily unavailable")
)

type message struct {
	ready   chan struct{}
	next    *message
	payload interface{}
}

func newNilMessage() *message {
	return &message{
		ready:   make(chan struct{}),
		next:    nil,
		payload: nil,
	}
}

type Channel struct {
	msg   *message
	enter chan interface{}
	closer.Closer
}

type Publisher struct {
	ch *Channel
}

type Subscriber struct {
	msg *message
}

func NewChannel(bufSize int) *Channel {
	ch := &Channel{
		msg:    newNilMessage(),
		enter:  make(chan interface{}, bufSize),
		Closer: closer.NewCloser(nil),
	}
	go ch.loop()
	return ch
}

func (ch *Channel) loop() {
	end := false
	for !end {
		select {
		case <-ch.Done():
			end = true
		case m := <-ch.enter:
			p := ch.msg

			ch.msg.payload = m
			ch.msg.next = newNilMessage()
			ch.msg = ch.msg.next

			close(p.ready)
		}
	}
	close(ch.msg.ready)
}

func (ch *Channel) NewPublisher() *Publisher {
	if ch.IsClosed() {
		return nil
	}
	return &Publisher{
		ch: ch,
	}
}

func (pub *Publisher) Send(m interface{}, timeout time.Duration) (
	err error) {
	var deadline <-chan time.Time
	if timeout > 0 {
		timer := time.NewTimer(timeout)
		defer timer.Stop()
		deadline = timer.C
	}

	select {
	case <-pub.ch.Done():
		return closer.ErrClosed
	case pub.ch.enter <- m:
		return nil
	case <-deadline:
		return ErrTimeout
	}
}

func (ch *Channel) NewSubscriber() *Subscriber {
	if ch.IsClosed() {
		return nil
	}
	return &Subscriber{
		msg: ch.msg,
	}
}

func (sub *Subscriber) recv() (m interface{}, err error) {
	m = sub.msg.payload
	sub.msg = sub.msg.next
	if sub.msg == nil {
		err = closer.ErrClosed
	}
	return
}

func (sub *Subscriber) Recv(timeout time.Duration) (
	m interface{}, err error) {
	if timeout > 0 {
		timer := time.NewTimer(timeout)
		defer timer.Stop()

		select {
		case <-sub.msg.ready:
			return sub.recv()
		case <-timer.C:
			return nil, ErrTimeout
		}
	} else {
		<-sub.msg.ready
		return sub.recv()
	}
}

func (sub *Subscriber) NonBlockRecv() (m interface{}, err error) {
	select {
	case <-sub.msg.ready:
		return sub.recv()
	default:
		return nil, ErrAgain
	}
}

func (sub *Subscriber) IsReady() bool {
	select {
	case <-sub.msg.ready:
		return true
	default:
		return false
	}
}

func (sub *Subscriber) Ready() <-chan struct{} {
	return sub.msg.ready
}
