package client

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	serverURL       = "ws://localhost:4444/ws" // Change this to your server URL
	numClients      = 1000                       // Number of WebSocket clients
	messageInterval = 1 * time.Second           // Interval between messages
)

func RunParallel() {
	var wg sync.WaitGroup

	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			runClient(id)
		}(i)
	}

	wg.Wait()
}

func runClient(id int) {
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Printf("[Client %d] Connection error: %v", id, err)
		return
	}
	defer conn.Close()

	log.Printf("[Client %d] Connected", id)

	ticker := time.NewTicker(messageInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			message := fmt.Sprintf("echo \"Hello from client %d - random %d\"", id, rand.Intn(1000))
			err := conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Printf("[Client %d] Write error: %v", id, err)
				return
			}
			log.Printf("[Client %d] Sent message: %s", id, message)
		}
	}
}
