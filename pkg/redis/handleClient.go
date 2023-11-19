package redis

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// HandleClient handles the incoming client connection.
func HandleClient(conn net.Conn, r *Redis) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	redis := []*Redis{r} // Perform a transaction on each Redis instance in the slice
	for scanner.Scan() {
		res, err := handleCommand(scanner.Text(), redis)
		if err != nil {
			conn.Write([]byte(err.Error() + "\n"))
		} else {
			conn.Write([]byte(res + "\n"))
		}
	}
}

// lastInstanceRedis returns the last Redis instance in the provided slice.
// If the slice is empty, it returns nil.
func lastInstanceRedis(redis []*Redis) *Redis {
	if len(redis) >= 1 {
		return redis[len(redis)-1]
	}
	return nil
}

func handleGet(key string, redis []*Redis) (string, error) {
	r := lastInstanceRedis(redis)
	result, err := r.Get(key)
	if err != nil {
		return "", err
	}
	return result, nil
}

func handleSet(key, value string, redis []*Redis) error {
	r := lastInstanceRedis(redis)
	return r.Set(key, value)
}

func handleDel(key string, redis []*Redis) error {
	r := lastInstanceRedis(redis)
	return r.Del(key)
}

func handleLPush(key, value string, redis []*Redis) error {
	r := lastInstanceRedis(redis)
	return r.LPush(key, value)
}

func handleLRange(key string, start, stop int, redis []*Redis) (string, error) {
	r := lastInstanceRedis(redis)
	result, err := r.LRange(key, start, stop)
	if err != nil {
		return "", err
	}
	return result, nil
}

func handleLPop(key string, redis []*Redis) (string, error) {
	r := lastInstanceRedis(redis)
	result, err := r.LPop(key)
	if err != nil {
		return "", err
	}
	return result, nil
}

// handleCommand processes the command received from the client.
func handleCommand(command string, redis []*Redis) (string, error) {
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
		result, err := handleGet(key, redis)
		if err != nil {
			return "", err
		}
		return result, nil

	case "set":
		if len(args) < 3 {
			return "", fmt.Errorf("set command requires at least two arguments")
		}
		key, value := args[1], args[2]
		handleSet(key, value, redis)
		return "OK", nil

	case "delete":
		if len(args) != 2 {
			return "", fmt.Errorf("delete command requires exactly one argument")
		}
		key := args[1]
		handleDel(key, redis)
		return "OK", nil

	case "lpush":
		if len(args) < 3 {
			return "", fmt.Errorf("lpush command requires at least two arguments")
		}
		key, value := args[1], args[2]
		handleLPush(key, value, redis)
		return "OK", nil

	case "lrange":
		if len(args) < 4 {
			return "", fmt.Errorf("lrange command requires exactly one argument")
		}
		key := args[1]
		startStr, stopStr := args[2], args[3]

		start, err := strconv.Atoi(startStr)
		if err != nil {
			return "", err
		}

		stop, err := strconv.Atoi(stopStr)
		if err != nil {
			return "", err
		}

		result, err := handleLRange(key, start, stop, redis)
		if err != nil {
			return "", err
		}
		return result, nil
	case "lpop":
		if len(args) != 2 {
			return "", fmt.Errorf("lpop command requires exactly one argument")
		}
		key := args[1]
		result, err := handleLPop(key, redis)
		if err != nil {
			return "", err
		}
		return result, nil
	default:
		return "", fmt.Errorf("unknown command: %s", command)
	}
}
