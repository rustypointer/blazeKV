package main

import (
	"blazeKV/internal/server"
	"blazeKV/internal/store"
	"blazeKV/internal/ttl"
)

func main() {
	totalKeysCapacity := 100000

	s := store.NewStore(totalKeysCapacity)

	ttl.StartCleaner(s)

	srv := server.NewTCPServer(s)

	srv.Start("8080")
}
