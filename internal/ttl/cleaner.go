package ttl

import (
	"blazeKV/internal/store"
	"time"
)

func StartCleaner(s *store.Store) {

	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for range ticker.C {
			s.CleanExpired()
		}
	}()
}
