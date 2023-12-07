package redis

import (
	"testing"
	"time"
)

var r *Redis

func init() {
	r = NewRedis()
	// Set without expiration
	expirationWithoutExpriation := 0
	r.Set("a", "1", time.Duration(expirationWithoutExpriation)*time.Second)

	// Set with expiration
	expirationWithExpriation := 2
	r.Set("b", "1", time.Duration(expirationWithExpriation)*time.Second)
}

func TestHandleGetCommandWithoutCache(t *testing.T) {

	redis := []*Redis{r}

	// Test case 1: Valid case with an existing key
	args1 := []string{"get", "a"}
	result1, _ := handleGet(args1, redis)
	if result1 != "1" {
		t.Errorf("handleGet() expected to return value from redis")
	}

	// Test case 2: Invalid case with non-existing key
	args2 := []string{"get", "not_existing_key"}
	_, err2 := handleGet(args2, redis)
	if err2 == nil {
		t.Errorf("handleGet() expected to yield error if the key does not exist")
	}

	// Test case 3: Invalid case with missing argument
	args3 := []string{"get"}
	_, err3 := handleGet(args3, redis)
	if err3 == nil {
		t.Error("Expected an error, got nil")
	}
}

func TestHandleGetCommandWithCache(t *testing.T) {
	redis := []*Redis{r}

	// Test case 1: Valid case with an existing key
	args1 := []string{"get", "b"}
	result1, _ := handleGet(args1, redis)
	if result1 != "1" {
		t.Errorf("handleGet() expected to return value from redis")
	}

	// Test case 2: Valid case with an key outdated
	time.Sleep(2 * time.Second)
	args2 := []string{"get", "b"}
	_, err2 := handleGet(args2, redis)
	if err2 == nil {
		t.Errorf("handleGet() expected to yield error if the key does not exist")
	}
}
