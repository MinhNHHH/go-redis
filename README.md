<img class="badge" alt="Go report card" tag="github.com/tmpaul/golang-redis-mock" src="https://goreportcard.com/badge/github.com/tmpaul/golang-redis-mock">

<img width="50%" height="50%" alt="Redis logo" src="https://upload.wikimedia.org/wikipedia/en/thumb/6/6b/Redis_Logo.svg/1200px-Redis_Logo.svg.png"/>

# golang-redis
A minimal functional Redis server written in Golang. I built this to learn Golang while simultaneouly
building out a functional product that begs good code practices, moderate use of concurrent goroutines
and dynamic type management.


# Architecture
- The redis is a hashtable with both key, value are string, lists
- The global store is initialized when the server starts and stored in RAM
- Each connection will be handled by a go-coroutine

## Running locally
Make sure that you have Go installed, and that it supports go modules.
```bash
go run main.go
```


Allowed commands are `GET`, `SET`, `DEL`, `GETSET`, `SETEX`, `LPUSH`, `LRANGE` and `TRANSACTION`.

## Running tests

```bash
go test -v ./...
```