package minhash

import (
	"encoding/binary"
	"strconv"
)

// toBytes converts various types to a byte slice
// so they can be pushed into a MinHash instance.
func toBytes(x interface{}) []byte {
	b := make([]byte, 8)
	switch t := x.(type) {
	case []byte:
		b = t
	case string:
		b = []byte(t)
	case uint:
		binary.LittleEndian.PutUint64(b, uint64(t))
	case uint16:
		binary.LittleEndian.PutUint64(b, uint64(t))
	case uint32:
		binary.LittleEndian.PutUint64(b, uint64(t))
	case uint64:
		binary.LittleEndian.PutUint64(b, t)
	case int:
		binary.LittleEndian.PutUint64(b, uint64(t))
	case int16:
		binary.LittleEndian.PutUint64(b, uint64(t))
	case int32:
		binary.LittleEndian.PutUint64(b, uint64(t))
	case int64:
		binary.LittleEndian.PutUint64(b, uint64(t))
	}
	return b
}

// stringIntToByte converts a string representation of an integer to a byte slice.
func stringIntToBytes(s string) []byte {
	n, err := strconv.ParseUint(s, 0, 64)
	var b []byte
	if err != nil {
		// Use literal conversion if string conversion failed.
		b = []byte(s)
	} else {
		b = toBytes(n)
	}
	return b
}
