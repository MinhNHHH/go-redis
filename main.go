package main

import (
	"fmt"
	"net"
	"sync"

	"github.com/MinhNHHH/redis/pkg/redis"
)

type Redis struct {
	clients map[string]net.Conn
	store   *redis.Store
	mutex   sync.Mutex
}

func main() {
	// Listen for incoming connections
	listener, err := net.Listen("tcp", "localhost:6789")

	r := &Redis{
		clients: make(map[string]net.Conn),
		store:   redis.NewStore(),
	}

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 6789")

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Handle client connection in a goroutine
		go redis.HandleClient(conn, r.store)
	}
}

// AddClient adds a new client connection to the Redis struct
func (r *Redis) AddClient(conn net.Conn) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Generate a unique key for the connection (you may need a better way to generate keys)
	key := fmt.Sprintf("%p", conn)

	// Add the connection to the map
	r.clients[key] = conn
}

// RemoveClient removes a client connection from the Redis struct
func (r *Redis) RemoveClient(conn net.Conn) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Remove the connection from the map
	delete(r.clients, fmt.Sprintf("%p", conn))
}
