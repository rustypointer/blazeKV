package store

import (
	"strconv"
	"testing"
)

// run -> go test ./internal/store -bench=. -benchmem

var result string

const totalKeysCapacity = 100000

func BenchmarkGet(b *testing.B) {

	store := NewStore(totalKeysCapacity)

	for i := 0; i < 100000; i++ {
		store.Set("key"+strconv.Itoa(i), "value")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		key := "key" + strconv.Itoa(i%100000)

		result, _ = store.Get(key)
	}

	_ = result
}

func BenchmarkSet(b *testing.B) {

	store := NewStore(totalKeysCapacity)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		key := "key" + strconv.Itoa(i)

		store.Set(key, "value")
	}
}

func BenchmarkDel(b *testing.B) {

	store := NewStore(totalKeysCapacity)

	for i := 0; i < b.N; i++ {
		store.Set("key"+strconv.Itoa(i), "value")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		key := "key" + strconv.Itoa(i)

		store.Del(key)
	}
}

func BenchmarkSetParallel(b *testing.B) {

	store := NewStore(totalKeysCapacity)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {

		i := 0

		for pb.Next() {

			key := "key" + strconv.Itoa(i)

			store.Set(key, "value")

			i++
		}
	})
}

func BenchmarkGetParallel(b *testing.B) {

	store := NewStore(totalKeysCapacity)

	for i := 0; i < 100000; i++ {
		store.Set("key"+strconv.Itoa(i), "value")
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {

		i := 0

		for pb.Next() {

			key := "key" + strconv.Itoa(i%100000)

			store.Get(key)

			i++
		}
	})
}
