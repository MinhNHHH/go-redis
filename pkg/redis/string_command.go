package redis

import (
	"fmt"
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
	if len(args) == 5 && strings.ToLower(args[3]) == "ex" {
		ttl, _ := strconv.ParseInt(args[4], 10, 64)
		r.Set(key, value, time.Duration(ttl)*time.Second)
	} else if len(args) == 3 {
		r.Set(key, value, time.Duration(0)*time.Second)
	} else {
		return "", fmt.Errorf("syntax error")
	}
	return "OK", nil
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
	r.SetEx(key, value, time.Duration(ttl)*time.Second)
	return "OK", nil
}

func handleDel(args []string, redis []*Redis) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("delete command requires exactly one argument")
	}
	key := args[1]
	r := currentRedis(redis)
	r.Del(key)
	return "OK", nil
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
