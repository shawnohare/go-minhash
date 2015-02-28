package minhash

import (
	"encoding/binary"
	"math"
	"strconv"
	"sync"
)

// MinWise is a data structure for generating a parametric family of
// hash functions of the form h1 + i*h2 for i=1, ..., k to compute
// a MinHash signature.  Each instance is tied to a single
// streamed set and hence single signature.  As it ingests
// elements it will update the current signature.

// The ith element of the signature is the current minimal value
// of the ith hash function.
type MinWise struct {
	minimums Signature
	h1       HashFunc
	h2       HashFunc
}

// NewMinWise constructs a new MinWise instance that will
// compute a signature of the specified size.
func NewMinWise(h1, h2 HashFunc, size int) *MinWise {
	mw := MinWise{
		minimums: defaultSignature(size), // running set of min values
		h1:       h1,
		h2:       h2,
	}
	return &mw
}

// Push updates the set's signature.  It hashes the input
// with each function in the family and compares these values
// with the current set of minimums, replacing them as necessary.
func (m *MinWise) Push(b []byte) {

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

// PushStringInt converts a string representation of an integer.
func (m *MinWise) PushStringInt(s string) {
	n, err := strconv.ParseUint(s, 0, 64)
	// Gracefully do not push if cannot convert.
	var b []byte
	if err != nil {
		log.Println("Could not convert string to uint64.")
		b = []bytes(s)
	} else {
		b = toBytes(n)
	}
	m.Push(b)
}

// PushGeneric deals with generic data by handling byte conversion.
func (m *MinWise) PushGeneric(x interface{}) {
	m.Push(toBytes(x))
}

func (m *MinWise) Signature() Signature {
	return m.minimums
}

func (m *MinWise) Similarity(m2 *MinWise) float64 {
	return MinWiseSimilarity(m.Signature(), m2.Signature())
}

// MinWiseSimilarity computes an estimate for the
// Jaccard similarity of two sets given their MinWise signatures.
func MinWiseSimilarity(sig1, sig2 Signature) float64 {
	if len(sig1) != len(sig2) {
		panic("signature size mismatch")
	}

	intersect := 0 // counter for number of elements in intersection

	for i := range sig1 {
		if sig1[i] == sig2[i] {
			intersect++
		}
	}

	return float64(intersect) / float64(len(sig1))
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

// Similarity computes an estimate for the similarity between the two sets.
func (m *MinWise) Similarity(m2 *MinWise) float64 {

	if len(m.minimums) != len(m2.minimums) {
		panic("minhash minimums size mismatch")
	}

	intersect := 0

	for i := range m.minimums {
		if m.minimums[i] == m2.minimums[i] {
			intersect++
		}
	}

	return float64(intersect) / float64(len(m.minimums))
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
