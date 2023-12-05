package redis

import (
	"fmt"
	"net"
	"strconv"
)

const (
	lpushCommand  = "lpush"
	lrangeCommand = "lrange"
	lpopCommand   = "lpop"
)

func handleLPush(args []string, redis []*Redis, conn net.Conn) {
	if len(args) < 3 {
		sendBackToClient(conn, "", fmt.Errorf("lpush command requires at least two arguments"))
		return
	}
	key, value := args[1], args[2]
	r := currentRedis(redis)
	r.LPush(key, value)
	sendBackToClient(conn, "OK", nil)
}

func handleLRange(args []string, redis []*Redis, conn net.Conn) {
	if len(args) < 4 {
		sendBackToClient(conn, "", fmt.Errorf("lrange command requires exactly one argument"))
		return
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

	r := currentRedis(redis)
	result, err := r.LRange(key, start, stop)
	if err != nil {
		sendBackToClient(conn, "", err)
	}
	sendBackToClient(conn, result, nil)
}

func handleLPop(args []string, redis []*Redis, conn net.Conn) {
	if len(args) != 2 {
		sendBackToClient(conn, "", fmt.Errorf("lpop command requires exactly one argument"))
		return
	}
	key := args[1]
	r := currentRedis(redis)
	result, err := r.LPop(key)
	if err != nil {
		sendBackToClient(conn, "", err)
	}
	sendBackToClient(conn, result, nil)
}
