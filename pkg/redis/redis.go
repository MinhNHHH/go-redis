package redis

import (
	"fmt"
	"sync"
	"time"
)

type ExpirationItem struct {
	value      interface{}
	expiration time.Time
}

// DB represents a simple in-memory database.
type Redis struct {
	items map[string]ExpirationItem
	mu    sync.Mutex // make sure only one goroutine can access a variable at a time to avoid conflicts
}

// NewRedis creates and returns a new instance of the DB.
func NewRedis() *Redis {
	return &Redis{
		items: map[string]ExpirationItem{},
	}
}

// Get retrieves the value associated with a key in the strings database.
func (r *Redis) Get(key string) (string, error) {
	r.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer r.mu.Unlock()
	if item, exist := r.items[key]; exist {
		if _, checkType := item.value.(string); checkType {
			if exist && item.expiration.After(time.Now()) {
				return item.value.(string), nil
			} else if item.expiration.IsZero() && exist {
				return item.value.(string), nil
			}
			delete(r.items, key)
		} else {
			// This error describe key existed in this database but command get key error
			return "", fmt.Errorf("wrongtype operation against a key holding the wrong kind of value")
		}
	}

	return "", fmt.Errorf("key not found")
}

// Set adds or updates a string value in the database.
func (r *Redis) Set(key, val string, expiration time.Duration) error {
	r.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer r.mu.Unlock()
	if item, exist := r.items[key]; exist {
		if _, ok := item.value.(string); !ok {
			// This error describe key existed in this database but command get key error
			return fmt.Errorf("wrongtype operation against a key holding the wrong kind of value")
		}
	}

	if expiration > 0 {
		r.items[key] = ExpirationItem{value: val, expiration: time.Now().Add(expiration)}
	} else {
		r.items[key] = ExpirationItem{value: val}
	}
	return nil
}

func (r *Redis) SetEx(key, val string, expiration time.Duration) error {
	r.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer r.mu.Unlock()
	if item, exist := r.items[key]; exist {
		if _, ok := item.value.(string); !ok {
			// This error describe key existed in this database but command get key error
			return fmt.Errorf("wrongtype operation against a key holding the wrong kind of value")
		}
	}
	r.items[key] = ExpirationItem{value: val, expiration: time.Now().Add(expiration)}
	return nil
}

// Delete a key in the database
func (r *Redis) Del(key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, key)
	return nil
}

//--------------------------------------------------------------------------------------------
func (r *Redis) LPush(key, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Check if the underlying type is []string
	// Update the value and assign it back to the interface field
	if item, ok := r.items[key]; ok {
		if strList, checkType := item.value.([]string); checkType {
			strList = append(strList, value)
			r.items[key] = ExpirationItem{value: strList}
		} else {
			// This error describe key existed in this database but command get key error
			return fmt.Errorf("wrongtype operation against a key holding the wrong kind of value")
		}
	} else {
		strList := []string{value}
		r.items[key] = ExpirationItem{value: strList}
	}
	// Handle the case where the key doesn't exist
	return nil
}

func (r *Redis) LRange(key string, start int, stop int) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if item, oke := r.items[key]; oke {
		if _, checkType := item.value.([]string); checkType {
			if start < 0 || start > len(item.value.([]string)) {
				return "", fmt.Errorf("lists startIndex out of range")
			}

			if stop > len(item.value.([]string)) {
				return "", fmt.Errorf("lists stopIndex out of range")
			}
			// if stop = -1 is the last element, stop = -2 is the penultimate element of the list, and so forth.
			if stop < 0 {
				stop = len(item.value.([]string)) + stop + 1
			}
			// If len(val) + stop + 1 < 0 => it should be return error.
			justString := fmt.Sprint(item.value.([]string)[start:stop])
			return justString, nil
		} else {
			// This error describe key existed in this database but command get key error
			return "", fmt.Errorf("wrongtype operation against a key holding the wrong kind of value")
		}
	}
	return "", fmt.Errorf("key not found")
}

func (r *Redis) LPop(key string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if item, ok := r.items[key]; ok {
		if _, checkType := item.value.([]string); checkType {
			element := item.value.([]string)[len(item.value.([]string))-1]
			if len(item.value.([]string)) > 1 {
				r.items[key] = ExpirationItem{value: item.value.([]string)[:len(item.value.([]string))-1]}
			} else {
				r.items[key] = ExpirationItem{value: []string{}}
			}
			return element, nil
		} else {
			// This error describe key existed in this database but command get key error
			return "", fmt.Errorf("wrongtype operation against a key holding the wrong kind of value")
		}
	}
	return "", fmt.Errorf("key not found")
}

//--------------------------------------------------------------------------------------------
// UpdateData merges the data from another Redis instance into the current instance.
// It acquires locks on both the current and new instances to ensure thread safety during the merge operation.
func (r *Redis) UpdateData(new *Redis) {
	// Acquire locks on both instances to prevent concurrent modification
	r.mu.Lock()
	defer r.mu.Unlock()

	// Merge string data
	for k, v := range new.items {
		r.items[k] = v
	}
}

// DeleteData removes keys from the current Redis instance that are not present in another Redis instance.
// It acquires locks on both the current and new instances to ensure thread safety during the deletion operation.
func (r *Redis) DeleteData(new *Redis) {
	// Acquire locks on both instances to prevent concurrent modification
	r.mu.Lock()
	defer r.mu.Unlock()

	for key := range r.items {
		if _, exists := new.items[key]; !exists {
			delete(r.items, key)
		}
	}
}
