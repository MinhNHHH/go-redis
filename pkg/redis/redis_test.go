package redis

import (
	"reflect"
	"testing"
	"time"
)

var store *Store

func init() {
	store = NewStore()
}

func TestGet(t *testing.T) {

	// Test case 1: Valid case with an existing key
	key := "existingKey"
	value := "someValue"
	expiration := time.Duration(0) * time.Second
	store.Set(key, value, expiration)

	result, err := store.Get(key)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != value {
		t.Errorf("Expected value %s, got %s", value, result)
	}

	// Test case 2: Invalid case with non-existing key
	key1 := "not_existing_key"
	_, err1 := store.Get(key1)
	if err1 == nil {
		t.Errorf("store.Get expected to yield error if the key does not exist")
	}
}

func TestGetWithExpireTime(t *testing.T) {
	key := "existingKeyWithExprieTime"
	value := "someValue"
	expiration := time.Duration(2) * time.Second
	store.Set(key, value, expiration)

	// Test case 1: Valid case with an existing key and has expried
	result, err := store.Get(key)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != value {
		t.Errorf("Expected value %s, got %s", value, result)
	}

	// Test case 2: Valid case with an key outdated
	time.Sleep(2 * time.Second)
	_, err2 := store.Get(key)
	if err2 == nil {
		t.Errorf("expected to yield error if the key does not exist: %v", err)
	}
}

func TestSet(t *testing.T) {
	// Test case 1: Valid case without expiration.
	key := "testSet"
	value := "1"
	expiration := time.Duration(0) * time.Second

	err := store.Set(key, value, expiration)
	if err != nil {
		t.Errorf("Expected an error: %v", err)
	}

	// Test case 2: Valid case with err expiration.
	key1 := "testSetWithErrorExp"
	value1 := "1"
	expiration1 := time.Duration(-1) * time.Second

	err1 := store.Set(key1, value1, expiration1)
	if err1 != nil {
		t.Errorf("Expected an error: %v", err)
	}
}

func TestSetEx(t *testing.T) {
	// Test case 1: Valid case without expiration.
	key := "testSetex"
	value := "1"
	expiration := time.Duration(1) * time.Second

	err := store.SetEx(key, value, expiration)
	if err != nil {
		t.Errorf("Expected an error: %v", err)
	}

	// Test case 2: Valid case with err expiration.
	key1 := "testSetWithErrorExp"
	value1 := "1"
	expiration1 := time.Duration(-1) * time.Second

	err1 := store.SetEx(key1, value1, expiration1)
	if err1 != nil {
		t.Errorf("Expected an error: %v", err)
	}
}

func TestDel(t *testing.T) {
	key := "testDel"
	value := "1"
	expiration := time.Duration(0) * time.Second
	store.Set(key, value, expiration)

	err := store.Del(key)
	if err != nil {
		t.Errorf("Expected an error: %v", err)
	}
}

func TestIncre(t *testing.T) {
	expiration := time.Duration(0) * time.Second

	// Test case 1: Valid with key is not exist
	key := "testInc"
	err := store.Incre(key)
	if err != nil {
		t.Errorf("Expected an error: %v", err)
	}

	// Test case 2: Valid with key is existing
	key1 := "testInc1"
	value1 := "1"
	store.Set(key1, value1, expiration)
	err2 := store.Incre(key1)
	if err2 != nil {
		t.Errorf("Expected an error: %v", err)
	}

	// Test case 3: Valid with key is not number
	key2 := "testInc2"
	value2 := "a"
	store.Set(key2, value2, expiration)
	err3 := store.Incre(key2)
	if err3 == nil {
		t.Error("Expected an error")
	}
}

func TestIncreBy(t *testing.T) {
	expiration := time.Duration(0) * time.Second
	increByValue := "10"

	// Test case 1: Valid with key is not exist
	key := "testIncBy"
	err := store.IncreBy(key, increByValue)
	if err != nil {
		t.Errorf("Expected an error: %v", err)
	}

	// Test case 2: Valid with key is existing
	key1 := "testIncBy1"
	value1 := "1"
	store.Set(key1, value1, expiration)
	err2 := store.IncreBy(key1, increByValue)
	if err2 != nil {
		t.Errorf("Expected an error: %v", err)
	}

	// Test case 3: Valid with key is not number
	key2 := "testIncBy2"
	value2 := "a"
	store.Set(key2, value2, expiration)
	err3 := store.IncreBy(key2, increByValue)
	if err3 == nil {
		t.Error("Expected an error")
	}
}

func TestDecre(t *testing.T) {
	expiration := time.Duration(0) * time.Second

	// Test case 1: Valid with key is not exist
	key := "testDecre"
	err := store.Decre(key)
	if err != nil {
		t.Errorf("Expected an error: %v", err)
	}

	// Test case 2: Valid with key is existing
	key1 := "testDecre1"
	value1 := "1"
	store.Set(key1, value1, expiration)
	err2 := store.Decre(key1)
	if err2 != nil {
		t.Errorf("Expected an error: %v", err)
	}

	// Test case 3: Valid with key is not number
	key2 := "testDecre2"
	value2 := "a"
	store.Set(key2, value2, expiration)
	err3 := store.Decre(key2)
	if err3 == nil {
		t.Error("Expected an error")
	}
}

func TestDecreBy(t *testing.T) {
	expiration := time.Duration(0) * time.Second
	DecreByValue := "-10"

	// Test case 1: Valid with key is not exist
	key := "testDecreBy"
	err := store.DecreBy(key, DecreByValue)
	if err != nil {
		t.Errorf("Expected an error: %v", err)
	}

	// Test case 2: Valid with key is existing
	key1 := "testDecreBy1"
	value1 := "1"
	store.Set(key1, value1, expiration)
	err2 := store.DecreBy(key1, DecreByValue)
	if err2 != nil {
		t.Errorf("Expected an error: %v", err)
	}

	// Test case 3: Valid with key is not number
	key2 := "testDecreBy2"
	value2 := "a"
	store.Set(key2, value2, expiration)
	err3 := store.DecreBy(key2, DecreByValue)
	if err3 == nil {
		t.Error("Expected an error")
	}
}

func TestUpdateData(t *testing.T) {
	// Create two stores with some initial data
	store1 := NewStore()
	store1.Set("a", "1", time.Duration(0)*time.Second)
	store1.Set("b", "2", time.Duration(0)*time.Second)

	store2 := NewStore()
	store2.Set("b", "3", time.Duration(0)*time.Second)
	store2.Set("c", "4", time.Duration(0)*time.Second)

	// Expected result after calling UpdateData on store1 with store2
	expectedStore := NewStore()
	expectedStore.Set("a", "1", time.Duration(0)*time.Second)
	expectedStore.Set("b", "3", time.Duration(0)*time.Second)
	expectedStore.Set("c", "4", time.Duration(0)*time.Second)

	// Call UpdateData on store1 with store2
	store1.UpdateData(store2)

	// Check if the modified store1 matches the expected result
	if !reflect.DeepEqual(store1.items, expectedStore.items) {
		t.Errorf("Store.UpdateData() did not merge the data as expected")
	}
}

func TestDeleteData(t *testing.T) {
	// Create two stores with some initial data
	store1 := NewStore()
	store1.Set("a", "1", time.Duration(0)*time.Second)
	store1.Set("b", "2", time.Duration(0)*time.Second)
	store1.Set("c", "3", time.Duration(0)*time.Second)

	store2 := NewStore()
	store2.Set("b", "2", time.Duration(0)*time.Second)

	// Expected result after calling DeleteData on store1 with store2
	expectedStore := NewStore()
	expectedStore.Set("b", "2", time.Duration(0)*time.Second)

	// Call DeleteData on store1 with store2
	store1.DeleteData(store2)

	// Check if the modified store1 matches the expected result
	if !reflect.DeepEqual(store1.items, expectedStore.items) {
		t.Errorf("Store.DeleteData() did not delete the expected keys")
	}
}
