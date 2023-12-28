package redis

import (
	"fmt"
	"net"
	"sync"
)

type RedisServer struct {
	clients map[string]*RedisClient
	store   *Store
	mutex   sync.Mutex
}
type RedisClient struct {
	ID   string
	conn net.Conn
}

func New() *RedisServer {
	return &RedisServer{
		clients: make(map[string]*RedisClient),
		store:   NewStore(),
	}
}

// AddClient adds a new client connection to the Redis struct
func (r *RedisServer) AddClient(conn *RedisClient) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	fmt.Printf("client %s connected \n", conn.ID)
	// Add the connection to the map
	r.clients[conn.ID] = conn
}

// RemoveClient removes a client connection from the Redis struct
func (r *RedisServer) RemoveClient(conn RedisClient) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	// Remove the connection from the map
	delete(r.clients, conn.ID)
	fmt.Printf("client %s disconnected \n", conn.ID)
	defer conn.conn.Close()
}

func Start(r *RedisServer) {
	// Listen for incoming connections
	listener, err := net.Listen("tcp", "localhost:6789")

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
		go HandleClient(conn, r)
	}
}
