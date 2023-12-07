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
	expirationWithExpriation := 1
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
		t.Error("Expected an error")
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

func TestHandleSet(t *testing.T) {
	redis := []*Redis{r}

	// Test case 1: Valid case without expiration.
	args1 := []string{"set", "test_set", "1"}
	handleSet(args1, redis)

	value, _ := currentRedis(redis).Get(args1[1])
	if value != args1[2] {
		t.Errorf("handlSet() expected to set value to redis")
	}

	// Test case 2: Invalid case with missing arguments
	args2 := []string{"set", "test_missing_args_set"}
	_, err2 := handleSet(args2, redis)
	if err2 == nil {
		t.Error("Expected an error")
	}

	// Test case 3: Valid case with expiration.
	args3 := []string{"set", "test_set_expriation", "1", "ex", "2"}
	handleSet(args3, redis)
	value3, _ := currentRedis(redis).Get(args3[1])
	if value3 != args1[2] {
		t.Errorf("handlSet() expected to set value to redis")
	}

	// Test case 4: Invalid case with key outdated
	time.Sleep(2 * time.Second)
	_, err4 := currentRedis(redis).Get(args3[1])
	if err4 == nil {
		t.Error("Expected an error")
	}

	// Test case 5: Invalid case with syntax error.
	args4 := []string{"set", "test_syntax_error", "1", "a"}
	_, err5 := handleSet(args4, redis)
	if err5 == nil {
		t.Error("Expected an error")
	}
}

func TestHandleSetEx(t *testing.T) {
	redis := []*Redis{r}

	// Test case 1: Valid case with minimum arguments
	args1 := []string{"setex", "test_setex", "1", "1"}
	handleSetEx(args1, redis)

	value, _ := currentRedis(redis).Get(args1[1])
	if value != args1[2] {
		t.Errorf("handleSetEx() expected to set value to redis")
	}

	// Test case 2: Invalid case with key outdated
	time.Sleep(2 * time.Second)
	_, err2 := currentRedis(redis).Get(args1[1])
	if err2 == nil {
		t.Error("Expected an error")
	}

	// Test case 3: Invalid case with missing arguments
	args3 := []string{"setex"}
	_, err3 := handleSet(args3, redis)
	if err3 == nil {
		t.Error("Expected an error")
	}

	// Test case 4: Valid case with set then setex.
	r.Set("test_set_then_setex", "1", time.Duration(0)*time.Second)

	args4 := []string{"setex", "test_set_then_setex", "1", "1"}
	handleSetEx(args4, redis)

	v, _ := currentRedis(redis).Get(args4[1])
	if v != args4[2] {
		t.Errorf("handleSetEx() expected to set value to redis")
	}

	time.Sleep(2 * time.Second)
	_, err5 := currentRedis(redis).Get(args4[1])
	if err5 == nil {
		t.Error("Expected an error")
	}
}

func TestHandDelete(t *testing.T) {
	redis := []*Redis{r}
	r.Set("test_del", "1", time.Duration(0)*time.Second)

	// Test case 1: Valid case with an existing key
	args1 := []string{"del", "test_del"}
	handleDel(args1, redis)
	_, err2 := currentRedis(redis).Get(args1[1])
	if err2 == nil {
		t.Error("Expected an error")
	}
}
