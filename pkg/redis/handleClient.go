package redis

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// HandleClient handles the incoming client connection.
func HandleClient(conn net.Conn, r *Redis) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		res, err := handleCommand(scanner.Text(), r)
		if err != nil {
			conn.Write([]byte(err.Error() + "\n"))
		} else {
			conn.Write([]byte(res + "\n"))
		}
	}
}

// handleCommand processes the command received from the client.
func handleCommand(command string, r *Redis) (string, error) {
	args := strings.Split(command, " ")
	if len(args) < 2 {
		return "", fmt.Errorf("error command")
	}
	switch strings.ToLower(args[0]) {
	case "get":
		if len(args) != 2 {
			return "", fmt.Errorf("get command requires exactly one argument")
		}
		key := args[1]
		result, err := r.Get(key)
		if err != nil {
			return "", err
		}
		return result, nil

	case "set":
		if len(args) < 3 {
			return "", fmt.Errorf("set command requires at least two arguments")
		}
		key, value := args[1], args[2]
		r.Set(key, value)
		return "OK", nil
	default:
		return "", fmt.Errorf("unknown command: %s", command)
	}
}
