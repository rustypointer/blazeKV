package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

const (
	serverAddr = "localhost:8080"

	numClients        = 10
	requestsPerClient = 10000

	keyspace  = 100000
	readRatio = 80

	pipelineSize = 50
)

func main() {

	fmt.Println("Preloading keys...")

	preloadData()

	fmt.Println("Preload complete")

	start := time.Now()

	var wg sync.WaitGroup

	for i := 0; i < numClients; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			runClient(id)
		}(i)
	}

	wg.Wait()

	duration := time.Since(start)

	totalRequests := numClients * requestsPerClient

	fmt.Println("------ BENCH RESULT ------")
	fmt.Println("Clients:", numClients)
	fmt.Println("Requests per client:", requestsPerClient)
	fmt.Println("Pipeline size:", pipelineSize)
	fmt.Println("Total requests:", totalRequests)
	fmt.Println("Read ratio:", readRatio, "% GET")
	fmt.Println("Keyspace:", keyspace)
	fmt.Println("Time:", duration)

	ops := float64(totalRequests) / duration.Seconds()

	fmt.Printf("Throughput: %.2f ops/sec\n", ops)
}

func preloadData() {

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for i := 0; i < keyspace; i++ {

		cmd := respSET(
			fmt.Sprintf("key%d", i),
			fmt.Sprintf("value%d", i),
		)

		conn.Write([]byte(cmd))

		reader.ReadString('\n')
	}
}

func runClient(id int) {

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	reader := bufio.NewReader(conn)

	rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(id)))

	for i := 0; i < requestsPerClient; i += pipelineSize {

		// send pipeline
		for j := 0; j < pipelineSize && i+j < requestsPerClient; j++ {

			var cmd string

			if rng.Intn(100) < readRatio {

				key := rng.Intn(keyspace)

				cmd = respGET(
					fmt.Sprintf("key%d", key),
				)

			} else {

				key := rng.Intn(keyspace)

				cmd = respSET(
					fmt.Sprintf("key%d", key),
					fmt.Sprintf("value%d", i+j),
				)
			}

			conn.Write([]byte(cmd))
		}

		// read responses
		for j := 0; j < pipelineSize && i+j < requestsPerClient; j++ {

			_, err := reader.ReadString('\n')
			if err != nil {
				return
			}
		}
	}
}

func respGET(key string) string {

	return fmt.Sprintf(
		"*2\r\n$3\r\nGET\r\n$%d\r\n%s\r\n",
		len(key),
		key,
	)
}

func respSET(key, value string) string {

	return fmt.Sprintf(
		"*3\r\n$3\r\nSET\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n",
		len(key),
		key,
		len(value),
		value,
	)
}
