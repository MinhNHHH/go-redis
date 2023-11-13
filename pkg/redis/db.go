package redis

import (
	"fmt"
	"sync"
)

// DB represents a simple in-memory database.
type DB struct {
	db map[string]interface{}
	mu sync.Mutex // make sure only one goroutine can access a variable at a time to avoid conflicts
}

// NewDB creates and returns a new instance of the DB.
func NewDB() *DB {
	return &DB{
		db: make(map[string]interface{}),
	}
}

// SaveToDB returns the current state of the database.
func (db *DB) SaveToDB() map[string]interface{} {
	return db.db
}

func (db *DB) handleGet(key string) (string, error) {
	db.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer db.mu.Unlock()
	if val, ok := db.db[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal, nil
		}
		return "", fmt.Errorf("value is not of type string")
	}
	return "", fmt.Errorf("key not found")
}

func (db *DB) handleSet(key string, val string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.db[key] = val
	return
}
