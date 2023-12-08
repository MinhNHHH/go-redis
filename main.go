package main

import (
	"fmt"
	"net"

	"github.com/MinhNHHH/redis/pkg/redis"
)

func main() {
	// Listen for incoming connections
	listener, err := net.Listen("tcp", "localhost:6789")
	r := redis.NewRedis()
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
		go redis.HandleClient(conn, r)
	}
}
