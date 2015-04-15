package minhash

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// dummy hash function
func h1(bs []byte) uint64 {
	return uint64(len(bs) + 1)
}

// Produce a new signature from input, if it's specified,
// else use a default size.
func newDummyMinWise(sig []uint64) *MinWise {
	var m *MinWise
	if len(sig) > 0 {
		m = NewMinWiseFromSignature(h1, h1, sig)
	} else {
		m = NewMinWise(h1, h1, 5)
	}

	return m
}

func TestCardinality(t *testing.T) {
	var testCases = []struct {
		sig  *MinWise
		card int
	}{
		{
			sig:  newDummyMinWise(nil),
			card: 0,
		},
		{
			sig:  newDummyMinWise([]uint64{0, 0}),
			card: 0,
		},
	}

	for _, tt := range testCases {
		assert.Equal(t, tt.card, tt.sig.Cardinality())
	}
}

func TestIsEmpty(t *testing.T) {
	var testCases = []*MinWise{
		newDummyMinWise(nil),
	}

	for _, tt := range testCases {
		assert.True(t, tt.IsEmpty())
	}
}

func TestSimilarity(t *testing.T) {
	var testCases = []struct {
		s1  *MinWise
		s2  *MinWise
		sim float64
	}{
		{
			s1:  newDummyMinWise(nil),
			s2:  newDummyMinWise(nil),
			sim: 1.0,
		},
		{
			s1:  newDummyMinWise([]uint64{1, 2}),
			s2:  newDummyMinWise([]uint64{1, 3}),
			sim: 0.5,
		},
		{
			s1:  newDummyMinWise(nil),
			s2:  newDummyMinWise([]uint64{1, 2, 3, 4, 5}),
			sim: 0.0,
		},
	}

	for _, tt := range testCases {
		assert.Equal(t, tt.sim, tt.s1.Similarity(tt.s2))
	}
}
