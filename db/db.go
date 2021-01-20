package db

import (
	"time"
)

// DB memory storage struct
type DB struct {
	// key namespace
	// use single goroutine interact with the data, lockfree
	data map[string]interface{}
	// expires data
	expires map[string]time.Time
}

// NewDB construct DB instance
func NewDB() *DB {
	return &DB{
		data: make(map[string]interface{}),
	}
}

// Set key=value to DB
func (db *DB) Set(key string, value interface{}) {
	db.data[key] = value
}

// Get key from DB
func (db *DB) Get(key string) interface{} {
	return db.data[key]
}
