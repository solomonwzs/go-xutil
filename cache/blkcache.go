package cache

import (
	"sync"

	"github.com/solomonwzs/goxutil/closer"
)

type block struct {
	data  interface{}
	atime int
	acnt  int
}

type BlkCacheOption struct {
}

type BlkCache struct {
	blkList []*block
	cap     int32
	lock    *sync.RWMutex
	closer.Closer
}
