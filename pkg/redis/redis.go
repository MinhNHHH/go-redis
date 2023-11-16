package redis

import (
	"fmt"
	"sync"
)

// DB represents a simple in-memory database.
type Redis struct {
	strings map[string]string // Data types Strings
	mu      sync.Mutex        // make sure only one goroutine can access a variable at a time to avoid conflicts
}

// NewDB creates and returns a new instance of the DB.
func NewRedis() *Redis {
	return &Redis{
		strings: map[string]string{},
	}
}

// Get retrieves the value associated with a key in the strings database.
func (r *Redis) Get(key string) (string, error) {
	r.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer r.mu.Unlock()
	if val, ok := r.strings[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("key not found")
}

// Set adds or updates a string value in the database.
func (r *Redis) Set(key string, val string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.strings[key] = val
	return nil
}
