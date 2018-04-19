package closer

import (
	"errors"
	"sync"
)

var ErrClosed = errors.New("closer was closed")

type Closer struct {
	end       chan struct{}
	endLock   *sync.Mutex
	closeFunc func() error
}

func NewCloser(closeFunc func() error) Closer {
	return Closer{
		end:       make(chan struct{}),
		endLock:   &sync.Mutex{},
		closeFunc: closeFunc,
	}
}

func (c Closer) IsClosed() bool {
	select {
	case <-c.end:
		return true
	default:
		return false
	}
}

func (c Closer) Close() (err error) {
	c.endLock.Lock()
	defer c.endLock.Unlock()

	if !c.IsClosed() {
		if c.closeFunc == nil {
			close(c.end)
		} else if err = c.closeFunc(); err == nil {
			close(c.end)
		}
		return
	}
	return ErrClosed
}

func (c Closer) Done() <-chan struct{} {
	return c.end
}
