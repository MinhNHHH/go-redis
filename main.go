package main

import (
	"github.com/MinhNHHH/redis/pkg/redis"
)

func main() {
	r := redis.New()
	redis.Start(r)
}
