package memory

import (
	"context"
	"sync"
)

// Database represents an in memory Database.
type Database struct {
	data map[string]any
	mu   sync.Mutex
}

func NewDatabase() *Database {
	return &Database{
		data: make(map[string]any),
	}
}

func (db *Database) Get(_ context.Context, key string) (any, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()
	v, ok := db.data[key]

	return v, ok
}

func (db *Database) Set(_ context.Context, key string, value any) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = value
}
