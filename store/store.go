// Package store handles the data persistence layer for the URL shortener.
// Because it is in a folder named "shorturl/store", this package name matches the folder.
package store

import (
	"fmt"  // Imported to use fmt.Errorf for creating dynamic error messages
	"sync" // Imported to use sync.RWMutex for safe concurrent access to the map
)

var _ Store = (*MemoryStore)(nil)

// Store defines the contract for our storage engine.
// Any struct that implements these two methods automatically satisfies this interface.
// Because "Store", "Save", and "Get" start with capital letters, they are public (exported).
type Store interface {
	Save(code string, original string) error
	Get(code string) (string, error)
}

// MemoryStore is an in-memory implementation of the Store interface.
// It wraps a standard Go map.
type MemoryStore struct {
	// mu protects the urls map from concurrent read/write access.
	// RWMutex lets multiple readers proceed together while still blocking writes.
	mu sync.RWMutex

	// urls is lowercase, making it private (unexported) to this package.
	// External packages cannot read or modify this map directly.
	urls map[string]string
}

// NewMemoryStore is a constructor function (factory).
// It safely initializes and returns a ready-to-use MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		// Crucial step: make() allocates memory for the map.
		// Writing to a nil (uninitialized) map causes a runtime panic.
		urls: make(map[string]string),
	}
}

// Save inserts or updates a short code mapping to an original URL.
// We use a pointer receiver (*MemoryStore) so we can modify the struct's internal map.
func (m *MemoryStore) Save(code string, original string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Store the original URL string inside the map, using the code as the unique key.
	m.urls[code] = original

	// Return nil to indicate that the operation completed successfully without errors.
	return nil
}

// Get looks up an original URL by its short code.
// We use a pointer receiver (*MemoryStore) for consistency across our method set.
func (m *MemoryStore) Get(code string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Use the "comma, ok" idiom to check if the key exists in the map.
	// 'value' will hold the original URL string if found.
	// 'ok' will be true if the key exists, or false if it does not.
	value, ok := m.urls[code]

	// If 'ok' is false, the requested short code is not in our map.
	if !ok {
		// Return an empty string and a formatted error explaining what went wrong.
		return "", fmt.Errorf("code not found: %s", code)
	}

	// If it exists, return the retrieved URL value and nil for the error.
	return value, nil
}
