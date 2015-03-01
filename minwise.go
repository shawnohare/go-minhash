package minhash

import (
	"math"
)

// MinWise is a data structure for generating a parametric family of
// hash functions of the form h1 + i*h2 for i=1, ..., k to compute
// a MinHash signature.  Each instance is tied to a single
// streamed set and hence single signature.  As it ingests
// elements it will update the current signature by calculating
// the ith hash function on the input and replacing the current
// ith element of the signature with the minimum of the hash
// and current value.

// The ith element of the signature is the current minimal value
// of the ith hash function.
type MinWise struct {
	minimums Signature
	h1       HashFunc
	h2       HashFunc
}

// NOTE MinWise constructors.

// NewMinWise constructs a new MinWise instance initialized with
// the empty set.
func NewMinWise(h1, h2 HashFunc, size int) *MinWise {
	mw := MinWise{
		minimums: defaultSignature(size), // running set of min values
		h1:       h1,
		h2:       h2,
	}
	return &mw
}

func NewMinWiseFromSignature(h1, h2 HashFunc, sig []uint64) *MinWise {
	mw := MinWise{
		minimums: sig,
		h1:       h1,
		h2:       h2,
	}
	return &mw
}

// InitStringIntMinWise creates a new MinWise instance and pushes a
// set of integers (represented as strings).  The returned instance
// contains the signature for the input set.
// func NewMinWiseFromStringInts(h1, h2 HashFunc, size int, xs []string) *MinWise {
// 	mw := NewMinWise(h1, h2, size)
// 	for _, x := range xs {
// 		mw.PushStringInt(x)
// 	}
// 	return mw
// }

// NOTE MinWise methods

// PushBytes updates the set's signature.  It hashes the input
// with each function in the family and compares these values
// with the current set of minimums, replacing them as necessary.
func (m *MinWise) PushBytes(b []byte) {

	v1 := m.h1(b)
	v2 := m.h2(b)

	// Compare minimal values
	for i, min := range m.minimums {
		// Compute hi(b) for ith hash function hi
		hb := v1 + uint64(i)*v2
		if hb < min {
			m.minimums[i] = hb
		}
	}
}

// PushStringInt first converts a string into a uint64 before pushing.
func (m *MinWise) PushStringInt(s string) {
	m.PushBytes(stringIntToByte(s))
}

// Push deals with generic data by handling byte conversion.
// It first hashes the input with each function in the instance's family,
// and then compares these values to the set of current minimums,
// updating them as necessary.
func (m *MinWise) Push(x interface{}) {
	m.PushBytes(toBytes(x))
}

// Signature returns the current signature.
func (m *MinWise) Signature() []uint64 {
	return m.minimums
}

// Similarity computes the similarity of two signatures represented
// as MinWise instances.  This estimates the Jaccard index of the
// two underlying sets.
func (m *MinWise) Similarity(m2 *MinWise) float64 {
	return MinWiseSimilarity(m.Signature(), m2.Signature())
}

// Merge combines the signatures of the second set,
// creating the signature of their union.
func (m *MinWise) Merge(m2 *MinWise) {

	for i, v := range m2.minimums {

		if v < m.minimums[i] {
			m.minimums[i] = v
		}
	}
}

// Cardinality estimates the cardinality of the set.
func (m *MinWise) Cardinality() int {

	// http://www.cohenwang.com/edith/Papers/tcest.pdf

	sum := 0.0

	for _, v := range m.minimums {
		sum += -math.Log(float64(infinity-v) / float64(infinity))
	}

	return int(float64(len(m.minimums)-1) / sum)
}

// SignatureBbit returns a b-bit reduction of the signature.  This will result in unused bits at the high-end of the words if b does not divide 64 evenly.
func (m *MinWise) SignatureBbit(b uint) []uint64 {

	var sig []uint64 // full signature
	var w uint64     // current word
	bits := uint(64) // bits free in current word

	mask := uint64(1<<b) - 1

	for _, v := range m.minimums {
		if bits >= b {
			w <<= b
			w |= v & mask
			bits -= b
		} else {
			sig = append(sig, w)
			w = 0
			bits = 64
		}
	}

	if bits != 64 {
		sig = append(sig, w)
	}

	return sig
}

// MinWiseSimilarity computes an estimate for the
// Jaccard similarity of two sets given their MinWise signatures.
func MinWiseSimilarity(sig1, sig2 []uint64) float64 {
	if len(sig1) != len(sig2) {
		panic("Signature size mismatch.")
	}

	intersect := 0 // counter for number of elements in intersection

	for i := range sig1 {
		if sig1[i] == sig2[i] {
			intersect++
		}
	}

	return float64(intersect) / float64(len(sig1))
}

// SimilarityBbit computes an estimate for the similarity between two b-bit signatures
func SimilarityBbit(sig1, sig2 []uint64, b uint) float64 {

	if len(sig1) != len(sig2) {
		panic("signature size mismatch")
	}

	intersect := 0
	count := 0

	mask := uint64(1<<b) - 1

	for i := range sig1 {
		w1 := sig1[i]
		w2 := sig2[i]

		bits := uint(64)

		for bits >= b {
			v1 := (w1 & mask)
			v2 := (w2 & mask)

			count++
			if v1 == v2 {
				intersect++
			}

			bits -= b
			w1 >>= b
			w2 >>= b
		}
	}

	return float64(intersect) / float64(count)
}
