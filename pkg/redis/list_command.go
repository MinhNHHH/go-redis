package redis

import (
	"fmt"
	"strconv"
)

const (
	lpushCommand  = "lpush"
	lrangeCommand = "lrange"
	lpopCommand   = "lpop"
)

func handleLPush(args []string, redis []*Redis) (string, error) {
	if len(args) < 3 {
		return "", fmt.Errorf("lpush command requires at least two arguments")
	}
	key, value := args[1], args[2]
	r := currentRedis(redis)
	r.LPush(key, value)
	return "OK", nil
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
