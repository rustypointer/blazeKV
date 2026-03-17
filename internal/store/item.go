package store

import "time"

type Item struct {
	Key        string
	Value      string
	ExpiryTime time.Time
	LastAccess int64
}
