package minhash

import (
	"github.com/dgryski/go-farm"
	"github.com/dgryski/go-spooky"
)

// New MinHash instance whose signature has the specified length.
// Unlike the other constructor functions, this one uses the default
// Spooky and Farm hash functions.
func New(length int) *MinHash {
	return NewMinHash(spooky.Hash64, farm.Hash64, length)
}
