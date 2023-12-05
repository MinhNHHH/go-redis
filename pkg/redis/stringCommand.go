package redis

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	getCommand          = "get"
	setCommand          = "set"
	getSetCommand       = "getset"
	deleteCommand       = "del"
	strLengthCommand    = "strlen"
	setAndExpireCommand = "setex"
)

func handleGet(args []string, redis []*Redis, conn net.Conn) {
	if len(args) != 2 {
		sendBackToClient(conn, "", fmt.Errorf("get command requires exactly one argument"))
		return
	}
	key := args[1]
	r := currentRedis(redis)
	result, err := r.Get(key)
	if err != nil {
		sendBackToClient(conn, "", err)
	}
	sendBackToClient(conn, result, nil)
}

func handleSet(args []string, redis []*Redis, conn net.Conn) {
	if len(args) < 3 {
		sendBackToClient(conn, "", fmt.Errorf("set command requires at least two arguments"))
		return
	}
	key, value := args[1], args[2]
	r := currentRedis(redis)
	if len(args) == 5 && strings.ToLower(args[3]) == "ex" {
		ttl, _ := strconv.ParseInt(args[4], 10, 64)
		r.Set(key, value, time.Duration(ttl)*time.Second)
	} else {
		r.Set(key, value, time.Duration(0)*time.Second)
	}
	sendBackToClient(conn, "OK", nil)
}

func handleSetEx(args []string, redis []*Redis, conn net.Conn) {
	if len(args) < 4 {
		sendBackToClient(conn, "", fmt.Errorf("set command requires at least three arguments"))
		return
	}
	ttl, _ := strconv.ParseInt(args[3], 10, 64)
	key, value := args[1], args[2]
	if ttl < 1 {
		sendBackToClient(conn, "", fmt.Errorf("invalid expire time in 'setex' command"))
	}
	r := currentRedis(redis)
	r.SetEx(key, value, time.Duration(ttl)*time.Second)
	sendBackToClient(conn, "OK", nil)
}

func handleDel(args []string, redis []*Redis, conn net.Conn) {
	if len(args) != 2 {
		sendBackToClient(conn, "", fmt.Errorf("delete command requires exactly one argument"))
		return
	}
	key := args[1]
	r := currentRedis(redis)
	r.Del(key)
	sendBackToClient(conn, "OK", nil)
}

func handleGetSet(args []string, redis []*Redis, conn net.Conn) {
	if len(args) < 3 {
		sendBackToClient(conn, "", fmt.Errorf("getset command requires at least two arguments"))
		return
	}
	key, value := args[1], args[2]
	r := currentRedis(redis)
	r.Set(key, value, time.Duration(0)*time.Second)
	result, err := r.Get(key)
	if err != nil {
		sendBackToClient(conn, "", err)
	}
	sendBackToClient(conn, result, nil)
}
