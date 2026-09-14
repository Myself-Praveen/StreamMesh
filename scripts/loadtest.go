package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

var (
	urlStr      = flag.String("url", "ws://localhost:8080/ws?user_id=loadtester", "WebSocket URL to connect to")
	connections = flag.Int("c", 1000, "Number of concurrent connections")
	duration    = flag.Duration("d", 10*time.Second, "Duration of the test")
	topic       = flag.String("t", "test-topic", "Topic to subscribe to")
)

func main() {
	flag.Parse()

	u, err := url.Parse(*urlStr)
	if err != nil {
		log.Fatal("Invalid URL:", err)
	}

	log.Printf("Starting load test with %d connections to %s for %v", *connections, u.String(), *duration)

	var (
		connectedCount atomic.Int32
		errorCount     atomic.Int32
		msgReceived    atomic.Int64
		wg             sync.WaitGroup
	)

	stopChan := make(chan struct{})

	// Start connections
	for i := 0; i < *connections; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				errorCount.Add(1)
				return
			}
			defer c.Close()
			connectedCount.Add(1)

			// Subscribe to topic
			subMsg := fmt.Sprintf(`{"action":"subscribe","topic":"%s"}`, *topic)
			if err := c.WriteMessage(websocket.TextMessage, []byte(subMsg)); err != nil {
				return
			}

			// Read messages loop
			go func() {
				for {
					_, _, err := c.ReadMessage()
					if err != nil {
						return
					}
					msgReceived.Add(1)
				}
			}()

			<-stopChan
		}(i)

		// Small delay to prevent overwhelming the server all at once
		time.Sleep(1 * time.Millisecond)
	}

	log.Printf("All %d connection attempts started. Active connections: %d. Errored: %d", *connections, connectedCount.Load(), errorCount.Load())
	log.Printf("Waiting %v for test to complete...", *duration)
	
	time.Sleep(*duration)
	close(stopChan)

	// Wait for all to exit cleanly (with timeout)
	doneChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneChan)
	}()

	select {
	case <-doneChan:
	case <-time.After(5 * time.Second):
		log.Println("Timeout waiting for connections to close")
	}

	// Print results
	log.Println("--- Load Test Results ---")
	log.Printf("Total Connections Attempted: %d", *connections)
	log.Printf("Successful Connections: %d", connectedCount.Load())
	log.Printf("Failed Connections: %d", errorCount.Load())
	log.Printf("Messages Received: %d", msgReceived.Load())
	
	if connectedCount.Load() > 0 {
		rate := float64(msgReceived.Load()) / duration.Seconds()
		log.Printf("Throughput: %.2f msg/sec", rate)
	}
}
