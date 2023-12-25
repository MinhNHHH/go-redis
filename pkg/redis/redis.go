package redis

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type ExpirationItem struct {
	value      interface{}
	expiration time.Time
}

// DB represents a simple in-memory database.
type Store struct {
	items map[string]ExpirationItem
	mu    sync.Mutex // make sure only one goroutine can access a variable at a time to avoid conflicts
}

// NewStore creates and returns a new instance of the DB.
func NewStore() *Store {
	return &Store{
		items: map[string]ExpirationItem{},
	}
}

// Get retrieves the value associated with a key in the strings database.
func (r *Store) Get(key string) (string, error) {
	r.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer r.mu.Unlock()
	if item, exist := r.items[key]; exist {
		if exist && item.expiration.After(time.Now()) {
			return item.value.(string), nil
		} else if item.expiration.IsZero() && exist {
			return item.value.(string), nil
		}
		delete(r.items, key)
	}
	return "", fmt.Errorf("key not found")
}

// Set adds or updates a string value in the database.
func (r *Store) Set(key, val string, expiration time.Duration) error {
	r.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer r.mu.Unlock()
	if expiration > 0 {
		r.items[key] = ExpirationItem{value: val, expiration: time.Now().Add(expiration)}
	} else {
		r.items[key] = ExpirationItem{value: val}
	}
	return nil
}

func (r *Store) SetEx(key, val string, expiration time.Duration) error {
	r.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer r.mu.Unlock()
	r.items[key] = ExpirationItem{value: val, expiration: time.Now().Add(expiration)}
	return nil
}

// Delete a key in the database
func (r *Store) Del(key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, key)
	return nil
}

func (r *Store) Incre(key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if item, exist := r.items[key]; exist {
		incrNumber, err := strconv.Atoi(item.value.(string))
		if err != nil {
			return err
		}
		incrNumber += 1
		value := strconv.Itoa(incrNumber)
		r.items[key] = ExpirationItem{value: value}
	} else {
		r.items[key] = ExpirationItem{value: "1"}
	}
	return nil
}

func (r *Store) IncreBy(key, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if item, exist := r.items[key]; exist {
		incrNumber, err := strconv.Atoi(item.value.(string))
		if err != nil {
			return err
		}
		valueIncred, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		incrNumber += valueIncred
		value := strconv.Itoa(incrNumber)
		r.items[key] = ExpirationItem{value: value}
	} else {
		r.items[key] = ExpirationItem{value: value}
	}
	return nil
}

func (r *Store) Decre(key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if item, exist := r.items[key]; exist {
		decrNumber, err := strconv.Atoi(item.value.(string))
		if err != nil {
			return err
		}
		decrNumber -= 1
		value := strconv.Itoa(decrNumber)
		r.items[key] = ExpirationItem{value: value}
	} else {
		r.items[key] = ExpirationItem{value: "-1"}
	}
	return nil
}

func (r *Store) DecreBy(key, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if item, exist := r.items[key]; exist {
		decrNumber, err := strconv.Atoi(item.value.(string))
		if err != nil {
			return err
		}
		valueIncred, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		decrNumber -= valueIncred
		value := strconv.Itoa(decrNumber)
		r.items[key] = ExpirationItem{value: value}
	} else {
		value = "-" + value
		r.items[key] = ExpirationItem{value: value}
	}
	return nil
}

// --------------------------------------------------------------------------------------------
func (r *Store) LPush(key, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Check if the underlying type is []string
	// Update the value and assign it back to the interface field
	if item, ok := r.items[key]; ok {
		if strList, checkType := item.value.([]string); checkType {
			strList = append(strList, value)
			r.items[key] = ExpirationItem{value: strList}
		}
	} else {
		strList := []string{value}
		r.items[key] = ExpirationItem{value: strList}
	}
	// Handle the case where the key doesn't exist
	return nil
}

func (r *Store) LRange(key string, start int, stop int) (string, error) {
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
		}
	}
	return "", fmt.Errorf("key not found")
}

func (r *Store) LPop(key string) (string, error) {
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
		}
	}
	return "", fmt.Errorf("key not found")
}

// --------------------------------------------------------------------------------------------
// UpdateData merges the data from another Store instance into the current instance.
// It acquires locks on both the current and new instances to ensure thread safety during the merge operation.
func (r *Store) UpdateData(new *Store) {
	// Acquire locks on both instances to prevent concurrent modification
	r.mu.Lock()
	defer r.mu.Unlock()

	// Merge string data
	for k, v := range new.items {
		r.items[k] = v
	}
}

// DeleteData removes keys from the current Store instance that are not present in another Store instance.
// It acquires locks on both the current and new instances to ensure thread safety during the deletion operation.
func (r *Store) DeleteData(new *Store) {
	// Acquire locks on both instances to prevent concurrent modification
	r.mu.Lock()
	defer r.mu.Unlock()

	for key := range r.items {
		if _, exists := new.items[key]; !exists {
			delete(r.items, key)
		}
	}
}
