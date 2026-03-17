package server

import (
	"blazeKV/internal/protocol"
	"blazeKV/internal/store"
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

type TCPServer struct {
	store *store.Store
}

func NewTCPServer(s *store.Store) *TCPServer {
	return &TCPServer{store: s}
}

func (t *TCPServer) Start(port string) error {

	ln, err := net.Listen("tcp", ":"+port)

	if err != nil {
		return err
	}

	fmt.Println("BlazeKV store running on port", port)

	for {

		conn, err := ln.Accept()

		if err != nil {
			continue
		}

		go t.handleConnection(conn)
	}
}

func (t *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Panic recovery so a bad command or internal error doesn't kill the goroutine
	defer func() {
		if r := recover(); r != nil {
			protocol.WriteError(writer, "ERR internal server error")
			writer.Flush()
		}
	}()

	for {
		// Optional: timeout to prevent stuck connections
		conn.SetReadDeadline(time.Now().Add(5 * time.Minute))

		// Read as many commands as are available in the TCP stream
		args, err := protocol.ReadCommand(reader)
		if err != nil {
			if err == io.EOF {
				return // client closed connection
			}
			protocol.WriteError(writer, "ERR invalid command")
			writer.Flush()
			return
		}

		if len(args) == 0 {
			protocol.WriteError(writer, "ERR empty command")
			continue
		}

		cmd := strings.ToUpper(args[0])

		switch cmd {
		case "PING":
			protocol.WriteSimpleString(writer, "PONG")

		case "SET":
			if len(args) < 3 {
				protocol.WriteError(writer, "ERR wrong number of arguments")
				break
			}
			key := args[1]
			val := args[2]
			t.store.Set(key, val)
			protocol.WriteSimpleString(writer, "OK")

		case "GET":
			if len(args) < 2 {
				protocol.WriteError(writer, "ERR wrong number of arguments")
				break
			}
			key := args[1]
			val, ok := t.store.Get(key)
			if !ok {
				protocol.WriteNull(writer)
			} else {
				protocol.WriteBulkString(writer, val)
			}

		case "DEL":
			if len(args) < 2 {
				protocol.WriteError(writer, "ERR wrong number of arguments")
				break
			}
			t.store.Del(args[1])
			protocol.WriteSimpleString(writer, "OK")

		case "EXPIRE":
			if len(args) < 3 {
				protocol.WriteError(writer, "ERR wrong number of arguments")
				break
			}
			sec, err := strconv.Atoi(args[2])
			if err != nil {
				protocol.WriteError(writer, "ERR invalid expire time")
				break
			}
			t.store.Expire(args[1], sec)
			protocol.WriteSimpleString(writer, "OK")

		default:
			protocol.WriteError(writer, "ERR unknown command")
		}

		// Flush after each command to maintain correct order for pipelining
		writer.Flush()
	}
}
