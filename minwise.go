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

// NewMinWise constructs a new instance and pushes the optional elements.
func NewMinWise(h1, h2 HashFunc, size int, elements ...interface{}) *MinWise {
	mw := &MinWise{
		minimums: defaultSignature(size), // running set of min values
		h1:       h1,
		h2:       h2,
	}

	for _, e := range elements {
		mw.Push(e)
	}

	return mw
}

func NewMinWiseFromSignature(h1, h2 HashFunc, sig []uint64) *MinWise {
	csig := make([]uint64, len(sig))
	copy(csig, sig)
	mw := MinWise{
		minimums: csig,
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

// IsEmpty reports whether the MinWise instance represents a signature
// of the empty set.  Note that it's possible for a signature of a
// non-empty set to equal the signature for the empty set in rare
// circumstances (e.g., when the hash family is not min-wise independent).
func (m *MinWise) IsEmpty() bool {
	return IsEmpty(m)
}

// Copy returns a new MinWise instance with the same type and data.
func (m *MinWise) Copy() *MinWise {
	return NewMinWiseFromSignature(m.h1, m.h2, m.Signature())
}

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
	m.PushBytes(stringIntToBytes(s))
}

// Push deals with generic data by handling byte conversion.
// It first hashes the input with each function in the instance's family,
// and then compares these values to the set of current minimums,
// updating them as necessary.
func (m *MinWise) Push(x interface{}) {
	m.PushBytes(toBytes(x))
}

// Signature returns  the  current signature.
func (m *MinWise) Signature() []uint64 {
	return m.minimums
}

// Similarity computes the similarity of two signatures represented
// as MinWise instances.  This estimates the Jaccard index of the
// two underlying sets.
func (m *MinWise) Similarity(m2 MinHash) float64 {
	return MinWiseSimilarity(m.Signature(), m2.Signature())
}

// Merge combines the signatures of the second set,
// creating the signature of their union.
func (m *MinWise) Merge(m2 MinHash) {

	for i, v := range m2.Signature() {

		if v < m.minimums[i] {
			m.minimums[i] = v
		}
	}
}

// Cardinality estimates the cardinality of the set. Both the
// signature for the empty set and the zero signature have 0 cardinality.
func (m *MinWise) Cardinality() int {

	// http://www.cohenwang.com/edith/Papers/tcest.pdf
	// Let h be a hash function.
	// With M as the the maximum uint64, the value
	// z := -ln( (M - min h(A)) / M)
	// is a sample from Z := min { X1 ~ Exp(1), ... , Xn ~ Exp(1)}.
	// where n := |A|.  But Z ~ Exp(n) has mean 1/n.
	// Thus the inverse of the mean of all the transformed signature values
	// is an estimate for |A|.
	var cardinality int
	sum := 0.0 // running sum

	// Use inverse transform to obtain a sample from Z for each sig element.
	for _, v := range m.minimums {

		d := infinity - v
		// Make sure the zero signature and sig for empty set return 0.
		if d == 0 || d == infinity {
			continue
		} else {
			sum += -math.Log(float64(d) / float64(infinity))
		}

	}

	if sum != 0.0 {
		cardinality = int(float64(len(m.minimums)) / sum)
	}

	return cardinality
}

// UnionCardinality estimates the cardinality of the union.
func (m *MinWise) UnionCardinality(m2 MinHash) int {
	union := m.Copy()
	union.Merge(m2)
	return union.Cardinality()
}

// IntersectionCardinality estimates the cardinality of the intersection.
func (m *MinWise) IntersectionCardinality(m2 MinHash) int {
	// Estimate size of the union.
	u := m.UnionCardinality(m2)

	// |A & B| + |A || B| = |A| +|B|
	est := m.Cardinality() + m2.Cardinality() - u
	// Take absolute value.
	if est < 0 {
		est = 0
	}

	return est
}

// SymmetricDifferenceCardinality estimates the difference between
// the cardinality of the union and intersection.
func (m *MinWise) SymmetricDifferenceCardinality(m2 MinHash) int {
	est := m.UnionCardinality(m2) - m.IntersectionCardinality(m2)
	if est < 0 {
		est = 0
	}

	return est
}

// LessCardinality estimates the cardinality of the left set minus
// the right set. This operator is not symmetric.
func (m *MinWise) LessCardinality(m2 MinHash) int {
	est := m.Cardinality() - m.IntersectionCardinality(m2)
	if est < 0 {
		est = 0
	}

	return est
}

// SignatureBbit returns a b-bit reduction of the signature.
// This will result in unused bits at the high-end of the words if b does not divide 64 evenly.
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
