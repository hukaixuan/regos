package internal

import (
	"sync"
	"time"
)

// TODO rebalance when load disk data and this value changed
const BucketNumber = 256

// TODO: change to sync.Map?
type db struct {
	// key namespace
	syncData   sync.Map
	dataBucket [BucketNumber]map[string]interface{}
	// expires data
	expires map[string]time.Time
}

func NewDB() *db {
	db := db{
		syncData:   sync.Map{},
		dataBucket: [BucketNumber]map[string]interface{}{},
	}
	for i := 0; i < len(db.dataBucket); i++ {
		db.dataBucket[i] = make(map[string]interface{})
	}
	return &db
}

func (db *db) Set(key string, value interface{}) error {
	return nil
}

func (db *db) Get(key string) (interface{}, error) {
	return nil, nil
}
