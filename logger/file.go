package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/solomonwzs/goxutil/closer"
)

const (
	ROTATE_NONE   = 0x00
	ROTATE_HOURLY = 0x01
	ROTATE_DAILY  = 0x02
)

type RotateFileOptions struct {
	path       string
	maxSize    uint64
	rotateTime int
}

func (opt RotateFileOptions) GetRotateTimer() <-chan time.Time {
	now := time.Now().Unix()
	rotate := int64(0)
	if opt.rotateTime == ROTATE_HOURLY {
		rotate = int64(time.Hour)
	} else if opt.rotateTime == ROTATE_DAILY {
		rotate = int64(24 * time.Hour)
	}

	if rotate != 0 {
		n := now / rotate
		dur := (n + 1) * rotate
		return time.After(time.Duration(dur - now))
	}
	return nil
}

type RotateFile struct {
	opt         RotateFileOptions
	fd          *os.File
	size        uint64
	rotateLock  *sync.RWMutex
	rotateTimer <-chan time.Time
	closer.Closer
}

func NewRotateFile(opt RotateFileOptions) (rf *RotateFile, err error) {
	stat, err := os.Stat(opt.path)
	if err != nil {
		return
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("%s is not dir", opt.path)
	}

	var rotateTimer <-chan time.Time = nil
	var rotateLock *sync.RWMutex = nil
	if rotateTimer = opt.GetRotateTimer(); rotateTimer != nil {
		rotateLock = &sync.RWMutex{}
	}

	rf = &RotateFile{
		opt:         opt,
		size:        0,
		rotateTimer: rotateTimer,
		rotateLock:  rotateLock,
	}
	rf.Closer = closer.NewCloser(func() error {
		if rf.rotateLock != nil {
			rf.rotateLock.Lock()
			defer rf.rotateLock.Unlock()
		}
		if rf.fd != nil {
			return rf.fd.Close()
		}
		return nil
	})

	return
}

func (rf *RotateFile) serv() {
	for {
		select {
		case <-rf.Done():
			return
		case <-rf.rotateTimer:
			return
		}
	}
}
