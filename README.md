# BlazeKV

A high-performance in-memory key-value store written in Go, inspired by Redis.  
Designed to explore real-world backend system design concepts like concurrency, eviction, protocol parsing, and high-throughput networking.

---

## Features

- 256-shard concurrent architecture
- Approximate LRU eviction (sampling-based)
- TTL support (lazy + active expiration)
- RESP protocol implementation
- TCP server with command pipelining
- Buffered I/O for high throughput
- Probabilistic access tracking to reduce contention

---

## Supported Commands

- `PING`
- `SET key value`
- `GET key`
- `DEL key`
- `EXPIRE key seconds`

---

## Architecture

Client → TCP → RESP Parser → Command Handler → Sharded Store

---

## Performance

- ~160K ops/sec (mixed workload)
- 100K requests benchmark
- 80% reads / 20% writes
- Pipeline size: 50

---

## Getting Started

### Run Server

    go run cmd/server/main.go

### Run Benchmark

    go run cmd/bench/main.go

---

## Example (RESP)

    *2
    $3
    GET
    $3
    key

---

## Key Concepts

- Sharding for scalable concurrency
- RWMutex-based synchronization
- Approximate LRU eviction
- Probabilistic metadata updates
- Lazy + sampled TTL expiration
- Network protocol parsing (RESP)
- Command pipelining

---

## Redis CLI Compatibility

BlazeKV supports the RESP protocol, allowing it to be tested using `redis-cli`.

Example:

    redis-cli -p 8080
    SET key value
    GET key

Supported commands:
- SET
- GET
- DEL
- EXPIRE
- PING

---

## Why BlazeKV?

BlazeKV is not just a CRUD project — it demonstrates:

- System design thinking
- Performance-oriented engineering
- Real-world tradeoffs used in production systems
- Understanding of how systems like Redis work internally

---

## Future Work

- Persistence (AOF / snapshotting)
- Additional commands (MGET, INCR)
- Metrics and observability
- Distributed clustering

---

## License

MIT