package mhsig

import (
	"math"
)

// MinWise is a data structure for generating a parametric family of
// hash functions of the form h1 + i*h2 for i=1, ..., k to compute
// a MinHash signatures.
type MinWise struct {
	// size is the signature length the instance will compute
	size int
	h1 Hash64Func
	h2 Hash64Func
}

// defaultSignature will return an appropriately typed array 
func defaultSignature(size int) Signature {
	s := make(Signature, size)
	for i := range s {
		s[i] = infinity
	}
	return s
}

func (m MinWise) Sketch(Set) Signature {
	// initialize an array of infinities
	mins := make([], m.size)
}
