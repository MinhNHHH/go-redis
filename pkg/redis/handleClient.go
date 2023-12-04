package redis

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// HandleClient handles the incoming client connection.
// It reads commands from the client, processes them, and sends responses back to the client.
func HandleClient(conn net.Conn, r *Redis) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	// Initialize the Redis slice with the original Redis instance
	// Perform a transaction on each Redis instance in the slice
	redis := []*Redis{r}
	for scanner.Scan() {
		args := strings.Split(strings.Trim(scanner.Text(), " "), " ")

		switch strings.ToLower(args[0]) {
		case "get":
			if len(args) != 2 {
				sendBackToClient(conn, "", fmt.Errorf("get command requires exactly one argument"))
				continue
			}
			key := args[1]
			result, err := handleGet(key, redis)
			if err != nil {
				sendBackToClient(conn, "", err)
			}
			sendBackToClient(conn, result, nil)

		case "set":
			if len(args) < 3 {
				sendBackToClient(conn, "", fmt.Errorf("set command requires at least two arguments"))
				continue
			}

			key, value := args[1], args[2]
			if len(args) == 5 && strings.ToLower(args[3]) == "ex" {
				ttl, _ := strconv.ParseInt(args[4], 10, 64)
				handleSet(key, value, time.Duration(ttl), redis)
			} else {
				handleSet(key, value, 0, redis)
			}

			sendBackToClient(conn, "OK", nil)

		case "delete":
			if len(args) != 2 {
				sendBackToClient(conn, "", fmt.Errorf("delete command requires exactly one argument"))
				continue
			}
			key := args[1]
			handleDel(key, redis)
			sendBackToClient(conn, "OK", nil)

		case "lpush":
			if len(args) < 3 {
				sendBackToClient(conn, "", fmt.Errorf("lpush command requires at least two arguments"))
				continue
			}
			key, value := args[1], args[2]
			handleLPush(key, value, redis)
			sendBackToClient(conn, "OK", nil)

		case "lrange":
			if len(args) < 4 {
				sendBackToClient(conn, "", fmt.Errorf("lrange command requires exactly one argument"))
				continue
			}
			key := args[1]
			startStr, stopStr := args[2], args[3]

			start, err := strconv.Atoi(startStr)
			if err != nil {
				sendBackToClient(conn, "", err)
			}

			stop, err := strconv.Atoi(stopStr)
			if err != nil {
				sendBackToClient(conn, "", err)
			}

			result, err := handleLRange(key, start, stop, redis)
			if err != nil {
				sendBackToClient(conn, "", err)
			}
			sendBackToClient(conn, result, nil)
		case "lpop":
			if len(args) != 2 {
				sendBackToClient(conn, "", fmt.Errorf("lpop command requires exactly one argument"))
				continue
			}
			key := args[1]
			result, err := handleLPop(key, redis)
			if err != nil {
				sendBackToClient(conn, "", err)
			}
			sendBackToClient(conn, result, nil)
		case "begin":
			// BEGIN command: Start a new transaction
			transaction := NewRedis()
			transaction.UpdateData(currentRedis(redis))
			redis = append(redis, transaction) // append the new transaction to the existing slice
			sendBackToClient(conn, "Transaction started", nil)
		case "commit":
			// COMMIT command: Commit the changes made during the transaction
			if len(redis) >= 2 {
				currentTransaction := currentRedis(redis)
				originalTransaction := redis[len(redis)-2]
				// Update originalTransaction with the values from currentTransaction
				originalTransaction.UpdateData(currentTransaction)
				// Delete keys in originalTransaction that are not present in currentTransaction
				originalTransaction.DeleteData(currentTransaction)
				redis = redis[:len(redis)-1]
				sendBackToClient(conn, "Transaction stoped", nil)
			}
		default:
			sendBackToClient(conn, "", fmt.Errorf("unknown command"))
		}
	}
}

func sendBackToClient(conn net.Conn, res string, err error) {
	if err != nil {
		conn.Write([]byte(err.Error() + "\n"))
	} else {
		conn.Write([]byte(res + "\n"))
	}
}

// currentRedis returns the last Redis instance in the provided slice.
// If the slice is empty, it returns nil.
func currentRedis(redis []*Redis) *Redis {
	if len(redis) >= 1 {
		return redis[len(redis)-1]
	}
	return nil
}

func handleGet(key string, redis []*Redis) (string, error) {
	r := currentRedis(redis)
	result, err := r.Get(key)
	if err != nil {
		return "", err
	}
	return result, nil
}

func handleSet(key, value string, expiration time.Duration, redis []*Redis) error {
	r := currentRedis(redis)
	return r.Set(key, value, expiration*time.Second)
}

func handleDel(key string, redis []*Redis) error {
	r := currentRedis(redis)
	return r.Del(key)
}

func handleLPush(key, value string, redis []*Redis) error {
	r := currentRedis(redis)
	return r.LPush(key, value)
}

func handleLRange(key string, start, stop int, redis []*Redis) (string, error) {
	r := currentRedis(redis)
	result, err := r.LRange(key, start, stop)
	if err != nil {
		return "", err
	}
	return result, nil
}

func handleLPop(key string, redis []*Redis) (string, error) {
	r := currentRedis(redis)
	result, err := r.LPop(key)
	if err != nil {
		return "", err
	}
	return result, nil
}
