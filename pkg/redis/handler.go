package redis

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const (
	beginCommand  = "begin"
	commitCommand = "commit"
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
		r, result, err := handleCommand(scanner.Text(), redis)
		redis = r
		sendBackToClient(conn, result, err)
	}
}

func handleCommand(command string, redis []*Redis) ([]*Redis, string, error) {
	args := strings.Split(strings.Trim(command, " "), " ")

	switch strings.ToLower(args[0]) {
	case getCommand:
		result, err := handleGet(args, redis)
		return redis, result, err
	case setCommand:
		result, err := handleSet(args, redis)
		return redis, result, err
	case setAndExpireCommand:
		result, err := handleSetEx(args, redis)
		return redis, result, err
	case deleteCommand:
		result, err := handleDel(args, redis)
		return redis, result, err
	case getSetCommand:
		result, err := handleGetSet(args, redis)
		return redis, result, err
	case lpushCommand:
		result, err := handleLPush(args, redis)
		return redis, result, err
	case lrangeCommand:
		result, err := handleLRange(args, redis)
		return redis, result, err
	case lpopCommand:
		result, err := handleLPop(args, redis)
		return redis, result, err
	case beginCommand:
		// BEGIN command: Start a new transaction
		redis = handleBegin(redis)
		return redis, "start transaction", nil
	case commitCommand:
		// COMMIT command: Commit the changes made during the transaction
		redis = handleCommit(redis)
		return redis, "stoped transaction", nil
	default:
		return redis, "", fmt.Errorf("unknown command")
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

func handleBegin(redis []*Redis) []*Redis {
	transaction := NewRedis()
	transaction.UpdateData(currentRedis(redis))
	redis = append(redis, transaction)
	return redis
}

func handleCommit(redis []*Redis) []*Redis {
	if len(redis) >= 2 {
		currentTransaction := currentRedis(redis)
		originalTransaction := redis[len(redis)-2]
		// Update originalTransaction with the values from currentTransaction
		originalTransaction.UpdateData(currentTransaction)
		// Delete keys in originalTransaction that are not present in currentTransaction
		originalTransaction.DeleteData(currentTransaction)
		redis = redis[:len(redis)-1]
	}
	return redis
}
