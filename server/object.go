package server

// basic data types
const (
	TYPE_STRING = iota
	TYPE_LIST
	TYPE_HASH
	TYPE_SET
	TYPE_ZSET
)

// encoding types
const (
	ENCODING_INT = iota
	ENCODING_STRING
	ENCODING_RAW
	ENCODING_HASHTABLE
	ENCODING_LINKEDLIST
	ENCODING_ZIPLIST
	ENCODING_INTSET
	ENCODING_SKIPLIST
)

type RegosObject struct {
	Type     uint8
	Encoding uint8
}

func (o *RegosObject) Exists() bool { return true }
