package redis

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const (
	beginCommand  = "begin"
	commitCommand = "command"
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
		case getCommand:
			handleGet(args, redis, conn)
		case setCommand:
			handleSet(args, redis, conn)
		case setAndExpireCommand:
			handleSetEx(args, redis, conn)
		case deleteCommand:
			handleDel(args, redis, conn)
		case getSetCommand:
			handleGetSet(args, redis, conn)
		case lpushCommand:
			handleLPush(args, redis, conn)
		case lrangeCommand:
			handleLRange(args, redis, conn)
		case lpopCommand:
			handleLPop(args, redis, conn)
		case beginCommand:
			// BEGIN command: Start a new transaction
			transaction := NewRedis()
			transaction.UpdateData(currentRedis(redis))
			redis = append(redis, transaction) // append the new transaction to the existing slice
			sendBackToClient(conn, "Transaction started", nil)
		case commitCommand:
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
