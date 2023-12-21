package redis

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type CommandType string

const (
	multiCommand        CommandType = "multi"
	execCommand         CommandType = "exec"
	discardCommand      CommandType = "discard"
	getCommand          CommandType = "get"
	setCommand          CommandType = "set"
	getSetCommand       CommandType = "getset"
	deleteCommand       CommandType = "del"
	strLengthCommand    CommandType = "strlen"
	setAndExpireCommand CommandType = "setex"
	lpushCommand        CommandType = "lpush"
	lrangeCommand       CommandType = "lrange"
	lpopCommand         CommandType = "lpop"
)

type Client struct {
	conn  net.Conn
	redis []*Redis
}

// HandleClient handles the incoming client connection.
// It reads commands from the client, processes them, and sends responses back to the client.
func HandleClient(conn net.Conn, r *Redis) {
	client := &Client{conn: conn, redis: []*Redis{r}}
	defer client.conn.Close()

	scanner := bufio.NewScanner(conn)
	// Initialize the Redis slice with the original Redis instance
	// Perform a transaction on each Redis instance in the slice
	for scanner.Scan() {
		client.handleCommand(scanner.Text())
	}
}

func (client *Client) handleCommand(command string) {

	args := strings.Split(strings.Trim(command, " "), " ")

	switch CommandType(strings.ToLower(args[0])) {
	case getCommand:
		result, err := handleGet(args, client.redis)

		if err != nil {
			sendBackToClient(client.conn, err.Error())
		} else {
			sendBackToClient(client.conn, result)
		}
	case setCommand:
		_, err := handleSet(args, client.redis)
		if err != nil {
			sendBackToClient(client.conn, err.Error())
		}
	case setAndExpireCommand:
		_, err := handleSetEx(args, client.redis)
		if err != nil {
			sendBackToClient(client.conn, err.Error())
		}
	case deleteCommand:
		_, err := handleDel(args, client.redis)
		if err != nil {
			sendBackToClient(client.conn, err.Error())
		}
	case getSetCommand:
		result, err := handleGetSet(args, client.redis)
		if err != nil {
			sendBackToClient(client.conn, err.Error())
		} else {
			sendBackToClient(client.conn, result)
		}
	case lpushCommand:
		_, err := handleLPush(args, client.redis)
		if err != nil {
			sendBackToClient(client.conn, err.Error())
		}
	case lrangeCommand:
		result, err := handleLRange(args, client.redis)
		if err != nil {
			sendBackToClient(client.conn, err.Error())
		} else {
			sendBackToClient(client.conn, result)
		}
	case lpopCommand:
		result, err := handleLPop(args, client.redis)
		if err != nil {
			sendBackToClient(client.conn, err.Error())
		} else {
			sendBackToClient(client.conn, result)
		}

	case multiCommand:
		// Multi command: Start a new transaction
		client.redis = handleMulti(client.redis)
		sendBackToClient(client.conn, "started transaction")
	case execCommand:
		// Exec command: Commit the changes made during the transaction
		redis, err := handleExec(client.redis)
		if err != nil {
			sendBackToClient(client.conn, err.Error())
		} else {
			client.redis = redis
			sendBackToClient(client.conn, "stoped transaction")
		}
	case discardCommand:
		// Discard command: Revert state the changes made during a transaction to bring the system back to a consistent state.
		redis, err := handleDiscard(client.redis)
		if err != nil {
			sendBackToClient(client.conn, err.Error())
		} else {
			client.redis = redis
			sendBackToClient(client.conn, "discarded transaction")
		}
	default:
		sendBackToClient(client.conn, "unknown command")
	}
}

func sendBackToClient(conn net.Conn, message string) {
	conn.Write([]byte(message + "\n"))
}

// currentRedis returns the last Redis instance in the provided slice.
// If the slice is empty, it returns nil.
func currentRedis(redis []*Redis) *Redis {
	if len(redis) >= 1 {
		return redis[len(redis)-1]
	}
	return nil
}

func handleMulti(redis []*Redis) []*Redis {
	transaction := NewRedis()
	transaction.UpdateData(currentRedis(redis))
	redis = append(redis, transaction)
	return redis
}

func handleExec(redis []*Redis) ([]*Redis, error) {
	if len(redis) >= 2 {
		currentTransaction := currentRedis(redis)
		originalTransaction := redis[len(redis)-2]
		// Update originalTransaction with the values from currentTransaction
		originalTransaction.UpdateData(currentTransaction)
		// Delete keys in originalTransaction that are not present in currentTransaction
		originalTransaction.DeleteData(currentTransaction)
		redis = redis[:len(redis)-1]
		return redis, nil
	}
	return redis, fmt.Errorf("exec without multi")
}

func handleDiscard(redis []*Redis) ([]*Redis, error) {
	if len(redis) >= 2 {
		return redis[:len(redis)-1], nil
	}
	return redis, fmt.Errorf("discard without multi")
}

func handleGet(args []string, redis []*Redis) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("get command requires exactly one argument")
	}
	key := args[1]
	r := currentRedis(redis)
	result, err := r.Get(key)
	if err != nil {
		return "", err
	}
	return result, nil
}

func handleSet(args []string, redis []*Redis) (string, error) {
	if len(args) < 3 {
		return "", fmt.Errorf("set command requires at least two arguments")
	}
	key, value := args[1], args[2]
	r := currentRedis(redis)
	var ttl int
	if len(args) == 5 && strings.EqualFold(args[3], "ex") {
		timeToLive, err := strconv.ParseInt(args[4], 10, 64)
		if err != nil {
			return "", fmt.Errorf("failed to parse TTL: %v", err)
		}
		ttl = int(timeToLive)
	} else if len(args) == 3 {
		ttl = 0
	} else {
		return "", fmt.Errorf("syntax error")
	}
	err := r.Set(key, value, time.Duration(ttl)*time.Second)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	return "", nil
}

func handleSetEx(args []string, redis []*Redis) (string, error) {
	if len(args) < 4 {
		return "", fmt.Errorf("setex command requires at least three arguments")
	}
	ttl, _ := strconv.ParseInt(args[3], 10, 64)
	key, value := args[1], args[2]
	if ttl < 1 {
		return "", fmt.Errorf("invalid expire time in 'setex' command")
	}
	r := currentRedis(redis)
	err := r.SetEx(key, value, time.Duration(ttl)*time.Second)
	if err != nil {
		return "", err
	}
	return "", nil
}

func handleDel(args []string, redis []*Redis) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("delete command requires exactly one argument")
	}
	key := args[1] //
	r := currentRedis(redis)
	r.Del(key)
	return "", nil
}

func handleGetSet(args []string, redis []*Redis) (string, error) {
	if len(args) < 3 {
		return "", fmt.Errorf("getset command requires at least two arguments")
	}
	key, value := args[1], args[2]
	r := currentRedis(redis)
	r.Set(key, value, time.Duration(0)*time.Second)
	result, err := r.Get(key)
	if err != nil {
		return "", err
	}
	return result, nil
}

func handleLPush(args []string, redis []*Redis) (string, error) {
	if len(args) < 3 {
		return "", fmt.Errorf("lpush command requires at least two arguments")
	}
	key, value := args[1], args[2]
	r := currentRedis(redis)
	err := r.LPush(key, value)
	if err != nil {
		return "", err
	}
	return "", nil
}

func handleLRange(args []string, redis []*Redis) (string, error) {
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

	r := currentRedis(redis)
	result, err := r.LRange(key, start, stop)
	if err != nil {
		return "", err
	}
	return result, nil
}

func handleLPop(args []string, redis []*Redis) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("lpop command requires exactly one argument")
	}
	key := args[1]
	r := currentRedis(redis)
	result, err := r.LPop(key)
	if err != nil {
		return "", err
	}
	return result, nil
}
