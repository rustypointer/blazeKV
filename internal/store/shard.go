package store

import (
	"sync"
)

type Shard struct {
	mu       sync.RWMutex
	m        map[string]*Item
	capacity int
}

func newShard(cap int) *Shard {
	return &Shard{
		m:        make(map[string]*Item),
		capacity: cap,
	}
}
