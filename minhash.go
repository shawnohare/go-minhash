package minhash

import (
	"errors"

	"github.com/dgryski/go-farm"
	"github.com/dgryski/go-spooky"
)

// MinHash is a data structure for generating a min-wise independent
// parametric family of hash functions of the form h1 + i*h2 for i=1, ..., k
// in order to to compute a MinHash signature.  Each instance is tied to a single
// streamed set and hence single signature.  As the instance
// ingests elements it keeps track of the minimum for the ith hash function.
type MinHash struct {
	mins []uint64 // mins[i] is the current min-value of ith hash func.
	h1   HashFunc
	h2   HashFunc
}

// New MinHash instance.  It is an alias for NewMinHash.
func New(h1, h2 HashFunc, size int) *MinHash {
	return NewMinHash(spooky.Hash64, farm.Hash64, size)
}

// NewMinHash constructs a new instance and pushes the optional elements.
func NewMinHash(h1, h2 HashFunc, size int) *MinHash {
	mw := &MinHash{
		mins: emptySetSignature(size), // running set of min values
		h1:   h1,
		h2:   h2,
	}

	return mw
}

func NewMinHashFromSignature(h1, h2 HashFunc, sig []uint64) *MinHash {
	csig := make([]uint64, len(sig))
	copy(csig, sig)
	mw := MinHash{
		mins: csig,
		h1:   h1,
		h2:   h2,
	}
	return &mw
}

// Signature returns the underlying signature slice.
func (m *MinHash) Signature() []uint64 {
	return m.mins
}

// Merge combines the signatures of the second set,
// creating the signature of their union.
func (m *MinHash) Merge(m2 Interface) error {

	if len(m.Signature()) != len(m2.Signature()) {
		return errors.New("Cannot merge signatures due to size mismatch.")
	}

	for i, v := range m2.Signature() {
		if v < m.mins[i] {
			m.mins[i] = v
		}
	}
	return nil
}

// IsEmpty reports whether the MinHash instance represents a signature
// of the empty set.  Note that it's possible for a signature of a
// non-empty set to equal the signature for the empty set in rare
// circumstances (e.g., when the hash family is not min-wise independent).
func (m *MinHash) IsEmpty() bool {
	return IsEmpty(m)
}

// Copy returns a new MinHash instance with the same type and data.
func (m *MinHash) Copy() *MinHash {
	return NewMinHashFromSignature(m.h1, m.h2, m.Signature())
}

// Push deals with generic data by handling byte conversion.
// It first hashes the input with each function in the instance's family,
// and then compares these values to the set of current mins,
// updating them as necessary.
func (m *MinHash) Push(x interface{}) {
	m.PushBytes(toBytes(x))
}

// PushBytes updates the set's signature.  It hashes the input
// with each function in the family and compares these values
// with the current set of mins, replacing them as necessary.
func (m *MinHash) PushBytes(b []byte) {

	v1 := m.h1(b)
	v2 := m.h2(b)

	// Compare minimal values
	for i, min := range m.mins {
		// Compute hi(b) for ith hash function hi
		hb := v1 + uint64(i)*v2
		// Ensure 0 is never pushed.
		if 0 < hb && hb < min {
			m.mins[i] = hb
		}
	}
}

// PushStringInt first converts a string into a uint64 before pushing.
func (m *MinHash) PushStringInt(s string) {
	m.PushBytes(stringIntToBytes(s))
}

// PushString casts the input as a []byte and pushes the element.
func (m *MinHash) PushString(s string) {
	m.PushBytes([]byte(s))
}

func (m *MinHash) PushStrings(ss ...string) {
	for _, s := range ss {
		m.PushString(s)
	}
}
