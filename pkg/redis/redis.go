package redis

import (
	"fmt"
	"sync"
	"time"
)

type cacheItem struct {
	value      string
	expiration time.Time
}

// DB represents a simple in-memory database.
type Redis struct {
	strings  map[string]string   // Data types Strings
	lists    map[string][]string // Data types Lists
	cache    map[string]cacheItem
	mu       sync.Mutex // make sure only one goroutine can access a variable at a time to avoid conflicts
	cacheMux sync.Mutex
}

// NewRedis creates and returns a new instance of the DB.
func NewRedis() *Redis {
	return &Redis{
		strings: map[string]string{},
		lists:   map[string][]string{},
		cache:   map[string]cacheItem{},
	}
}

// Get retrieves the value associated with a key in the strings database.
func (r *Redis) Get(key string) (string, error) {
	r.mu.Lock()
	r.cacheMux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer r.mu.Unlock()
	defer r.cacheMux.Unlock()

	if cache, ok := r.cache[key]; ok {
		return cache.value, nil
	}

	if val, ok := r.strings[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("key not found")
}

// Set adds or updates a string value in the database.
func (r *Redis) Set(key, val string, expiration time.Duration) error {
	r.mu.Lock()
	r.cacheMux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer r.mu.Unlock()
	defer r.cacheMux.Unlock()

	r.strings[key] = val
	r.cache[key] = cacheItem{value: val, expiration: time.Now().Add((expiration))}
	return nil
}

// Delete a key in the database
func (r *Redis) Del(key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.strings, key)
	return nil
}

func (r *Redis) LPush(key, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lists[key] = append(r.lists[key], value)
	return nil
}

func (r *Redis) LRange(key string, start int, stop int) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if val, oke := r.lists[key]; oke {

		if start < 0 || start > len(val) {
			return "", fmt.Errorf("lists startIndex out of range")
		}

		if stop > len(val) {
			return "", fmt.Errorf("lists stopIndex out of range")
		}
		// if stop = -1 is the last element, stop = -2 is the penultimate element of the list, and so forth.
		if stop < 0 {
			stop = len(val) + stop + 1
		}
		// If len(val) + stop + 1 < 0 => it should be return error.
		justString := fmt.Sprint(val[start:stop])
		return justString, nil
	}
	return "", fmt.Errorf("key not found")
}

func (r *Redis) LPop(key string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if val, ok := r.lists[key]; ok {
		element := val[len(val)-1]
		if len(val) > 1 {
			r.lists[key] = val[:len(val)-1]
		} else {
			r.lists[key] = []string{}
		}
		return element, nil
	}
	return "", fmt.Errorf("key not found")
}

// UpdateData merges the data from another Redis instance into the current instance.
// It acquires locks on both the current and new instances to ensure thread safety during the merge operation.
func (r *Redis) UpdateData(new *Redis) {
	// Acquire locks on both instances to prevent concurrent modification
	r.mu.Lock()
	new.mu.Lock()
	defer r.mu.Unlock()
	defer new.mu.Unlock()
	// Merge string data
	for k, v := range new.strings {
		r.strings[k] = v
	}

	for k, v := range new.lists {
		r.lists[k] = v
	}
}

// DeleteData removes keys from the current Redis instance that are not present in another Redis instance.
// It acquires locks on both the current and new instances to ensure thread safety during the deletion operation.
func (r *Redis) DeleteData(new *Redis) {
	// Acquire locks on both instances to prevent concurrent modification
	r.mu.Lock()
	new.mu.Lock()
	defer r.mu.Unlock()
	defer new.mu.Unlock()

	for key := range r.strings {
		if _, exists := new.strings[key]; !exists {
			delete(r.strings, key)
		}
	}

	for key := range r.lists {
		if _, exists := new.lists[key]; !exists {
			delete(r.lists, key)
		}
	}
}
