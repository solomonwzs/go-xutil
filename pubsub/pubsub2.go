package pubsub

import (
	"sync"
	"time"

	"github.com/solomonwzs/goxutil/closer"
)

type _Channel struct {
	enter chan interface{}
	m     map[uint32]*_Subscriber
	mlock *sync.Mutex
	pid   uint32
	closer.Closer
}

type _Publisher struct {
	ch *_Channel
}

type _Subscriber struct {
	ch    *_Channel
	msgCh chan interface{}
	closer.Closer
}

func _newChannel() *_Channel {
	ch := &_Channel{
		enter:  make(chan interface{}, 100),
		m:      make(map[uint32]*_Subscriber),
		mlock:  &sync.Mutex{},
		pid:    0,
		Closer: closer.NewCloser(nil),
	}
	go ch.loop()
	return ch
}

func (ch *_Channel) loop() {
	end := false
	for !end {
		select {
		case <-ch.Done():
			end = true
		case m := <-ch.enter:
			ch.mlock.Lock()
			closedIDs := []uint32{}
			for id, sub := range ch.m {
				select {
				case <-sub.Done():
					closedIDs = append(closedIDs, id)
				case sub.msgCh <- m:
				}
			}

			for _, id := range closedIDs {
				delete(ch.m, id)
			}
			ch.mlock.Unlock()
		}
	}
	ch.mlock.Lock()
	for _, sub := range ch.m {
		close(sub.msgCh)
	}
	ch.mlock.Unlock()
}

func (ch *_Channel) newPublisher() *_Publisher {
	if ch.IsClosed() {
		return nil
	}
	return &_Publisher{
		ch: ch,
	}
}

func (pub *_Publisher) send(m interface{}, timeout time.Duration) (
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

func (ch *_Channel) newSubscriber() *_Subscriber {
	if ch.IsClosed() {
		return nil
	}
	sub := &_Subscriber{
		msgCh:  make(chan interface{}, 10),
		Closer: closer.NewCloser(nil),
		ch:     ch,
	}

	ch.mlock.Lock()
	defer ch.mlock.Unlock()

	ch.pid += 1
	ch.m[ch.pid] = sub

	return sub
}

func (sub *_Subscriber) recv(timeout time.Duration) (interface{}, error) {
	if timeout > 0 {
		timer := time.NewTimer(timeout)
		defer timer.Stop()

		select {
		case m, ok := <-sub.msgCh:
			if !ok {
				return nil, closer.ErrClosed
			}
			return m, nil
		case <-timer.C:
			return nil, ErrTimeout
		}
	} else {
		m, ok := <-sub.msgCh
		if !ok {
			return nil, closer.ErrClosed
		}
		return m, nil
	}
}
