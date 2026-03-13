package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
	totalRequests := 3000
	concurrency := 10
	// We target the Docker service name 'node1'
	target := "node1:9001" 

	var wg sync.WaitGroup
	start := time.Now()

	fmt.Printf("🚀 Starting Benchmark: %d requests | Concurrency: %d\n", totalRequests, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < totalRequests/concurrency; j++ {
				conn, err := net.DialTimeout("tcp", target, 2*time.Second)
				if err != nil {
					continue
				}
				fmt.Fprintf(conn, "SET key-%d-%d value-%d\n", workerID, j, j)
				bufio.NewReader(conn).ReadString('\n')
				conn.Close()
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)
	rps := float64(totalRequests) / duration.Seconds()

	fmt.Printf("\n📊 --- BENCHMARK RESULTS ---\n")
	fmt.Printf("Total Time:  %v\n", duration)
	fmt.Printf("Throughput:  %.2f req/sec\n", rps)
	fmt.Printf("---------------------------\n")
}