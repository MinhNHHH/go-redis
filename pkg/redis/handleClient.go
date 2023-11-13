package redis

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// HandleClient handles the incoming client connection.
func HandleClient(conn net.Conn, db *DB) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		handleCommand(scanner.Text(), db)
	}
}

// handleCommand processes the command received from the client.
func handleCommand(input string, db *DB) {
	args := strings.Split(input, " ")
	// if len(args) != 2 {
	// 	fmt.Println("Args need to be len > 3")
	// 	return
	// }

	command := args[0]
	switch command {
	case "get":
		db.handleGet(args[0])
	case "set":
		db.handleSet(args[1], args[2])
	case "delete":
		fmt.Println("delete")
	case "exit":
		fmt.Println("exit")
	case "list":
		fmt.Println("list all")
	default:
		fmt.Println("unknown command")
	}
}
