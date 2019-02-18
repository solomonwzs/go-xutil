package logger

import (
	"fmt"
	"os"
	"path/filepath"
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
	Name       string
	Path       string
	MaxSize    uint64
	RotateTime int
}

func (opt RotateFileOptions) GetRotateTimer() <-chan time.Time {
	now := time.Now().Unix()
	rotate := int64(0)
	if opt.RotateTime == ROTATE_HOURLY {
		rotate = int64(time.Hour)
	} else if opt.RotateTime == ROTATE_DAILY {
		rotate = int64(24 * time.Hour)
	}

	if rotate != 0 {
		n := now / rotate
		dur := (n + 1) * rotate
		return time.After(time.Duration(dur - now))
	}
	return nil
}

func (opt RotateFileOptions) GetFilename() string {
	timeStr := ""
	now := time.Now()
	if opt.RotateTime == ROTATE_HOURLY {
		timeStr = now.Format(".2006_01_02_15")
	} else if opt.RotateTime == ROTATE_DAILY {
		timeStr = now.Format(".2006_01_02")
	}
	name := filepath.Join(opt.Path, opt.Name+timeStr+".log")
	return name
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
	stat, err := os.Stat(opt.Path)
	if err != nil {
		return
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("%s is not dir", opt.Path)
	}

	rf = &RotateFile{
		opt:         opt,
		size:        0,
		rotateTimer: opt.GetRotateTimer(),
		rotateLock:  &sync.RWMutex{},
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

	name := opt.GetFilename()
	fd, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	rf.fd = fd

	go rf.serv()
	return
}

func (rf *RotateFile) serv() {
	for {
		select {
		case <-rf.Done():
			return
		case <-rf.rotateTimer:
			rf.rotateLock.Lock()
			rf.rotateLock.Unlock()
			rf.rotateTimer = rf.opt.GetRotateTimer()

			if err := rf.Rotate(); err != nil {
				rf.Close()
			}
		}
	}
}

func (rf *RotateFile) Rotate() (err error) {
	name := rf.fd.Name()
	if err = rf.fd.Close(); err != nil {
		return
	}

	for i := 0; true; i++ {
		newName := fmt.Sprintf("%s.%d", name, i)
		_, err0 := os.Stat(newName)
		if os.IsNotExist(err0) {
			if err = os.Rename(name, newName); err != nil {
				return
			}
			break
		}
	}

	fd, err := os.OpenFile(rf.opt.GetFilename(),
		os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return
	}
	rf.fd = fd
	rf.size = 0

	return
}

func (rf *RotateFile) Write(p []byte) (n int, err error) {
	rf.rotateLock.Lock()
	defer rf.rotateLock.Unlock()

	if rf.opt.MaxSize != 0 && rf.size+uint64(len(p)) > rf.opt.MaxSize {
		if err = rf.Rotate(); err != nil {
			return
		}
	}

	n, err = rf.fd.Write(p)
	rf.size += uint64(n)
	return
}
