package store

import (
	"blazeKV/pkg/hash"
	"math/rand"
	"time"
)

const shardCount = 256
const sampleSize = 5

type Store struct {
	shards        [shardCount]*Shard
	totalCapacity int
}

func NewStore(totalCap int) *Store {

	s := &Store{
		totalCapacity: totalCap,
	}

	perShardCapacity := totalCap / shardCount

	for i := 0; i < shardCount; i++ {
		s.shards[i] = newShard(perShardCapacity)
	}

	return s
}

func (s *Store) getShard(key string) *Shard {

	index := hash.Hash(key) % shardCount

	return s.shards[index]
}

func (s *Store) Set(key, value string) {

	shard := s.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	if item, ok := shard.m[key]; ok {
		item.Value = value
		item.LastAccess = time.Now().UnixNano()
		return
	}

	item := &Item{
		Key:        key,
		Value:      value,
		LastAccess: time.Now().UnixNano(),
	}

	shard.m[key] = item

	if len(shard.m) > shard.capacity {
		s.evict(shard)
	}
}

func (s *Store) Get(key string) (string, bool) {

	shard := s.getShard(key)

	shard.mu.RLock()

	item, ok := shard.m[key]
	if !ok {
		shard.mu.RUnlock()
		return "", false
	}

	if !item.ExpiryTime.IsZero() && time.Now().After(item.ExpiryTime) {
		shard.mu.RUnlock()

		shard.mu.Lock()
		delete(shard.m, key)
		shard.mu.Unlock()

		return "", false
	}

	val := item.Value
	shard.mu.RUnlock()

	// probabilistic access update for approximate lru eviction
	if rand.Intn(100) < 20 {
		shard.mu.Lock()
		if item, ok := shard.m[key]; ok {
			item.LastAccess = time.Now().UnixNano()
		}
		shard.mu.Unlock()
	}

	return val, true
}

func (s *Store) Del(key string) {

	shard := s.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	if _, ok := shard.m[key]; !ok {
		return
	}

	delete(shard.m, key)
}

func (s *Store) Expire(key string, seconds int) {

	shard := s.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	item, ok := shard.m[key]
	if !ok {
		return
	}

	item.ExpiryTime = time.Now().Add(time.Duration(seconds) * time.Second)
}

func (s *Store) CleanExpired() {

	now := time.Now()

	for _, shard := range s.shards {

		shard.mu.Lock()

		i := 0

		for key, item := range shard.m {

			if i >= sampleSize {
				break
			}

			if !item.ExpiryTime.IsZero() && now.After(item.ExpiryTime) {
				delete(shard.m, key)
			}

			i++
		}

		shard.mu.Unlock()
	}
}

func (s *Store) evict(shard *Shard) {

	if len(shard.m) == 0 {
		return
	}

	var oldestKey string
	oldestTime := time.Now().UnixNano()

	i := 0

	for k, item := range shard.m {

		if i >= sampleSize {
			break
		}

		if item.LastAccess < oldestTime {
			oldestTime = item.LastAccess
			oldestKey = k
		}

		i++
	}

	if oldestKey == "" {
		return
	}

	delete(shard.m, oldestKey)
}
