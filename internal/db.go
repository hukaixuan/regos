package internal

import (
	"time"
)

// TODO: change to sync.Map?
type db struct {
	// key namespace
	data map[string]interface{}
	// expires data
	expires map[string]time.Time
}

func NewDB() *db {
	return &db{
		data: make(map[string]interface{}),
	}
}

func (db *db) Set(key string, value interface{}) error {
	return nil
}

func (db *db) Get(key string) (interface{}, error) {
	return nil, nil
}
