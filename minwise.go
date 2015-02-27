package mhsig

import (
	"math"
)

// MinWise is a data structure for generating a parametric family of
// hash functions of the form h1 + i*h2 for i=1, ..., k to compute
// a MinHash signatures.
type MinWise struct {
	h1 HashFunc
	h2 HashFunc
}
