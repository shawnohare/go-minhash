package minhash

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dgryski/go-farm"
	"github.com/dgryski/go-spooky"
)

var h1 = farm.Hash64
var h2 = spooky.Hash64

// Two signatures.

func makeSigOfInts() *MinWise {
	var sig = NewMinWise(h1, h2, 400)
	for i := 0; i <= 10000; i++ {
		sig.Push(i)
	}
	return sig
}

func makeSigOfEvens() *MinWise {
	var sig = NewMinWise(h1, h2, 400)
	for i := 0; i <= 10000; i++ {
		if i%2 == 0 {
			sig.Push(i)
		}
	}
	return sig
}

func makeSigOfOdds() *MinWise {
	var sig = NewMinWise(h1, h2, 400)
	for i := 0; i <= 10000; i++ {
		if i%2 == 1 {
			sig.Push(i)
		}
	}
	return sig
}

// Produce a new signature from input, if it's specified,
// else use a default size.
func newDummyMinWise(sig []uint64) *MinWise {
	var m *MinWise
	if len(sig) > 0 {
		m = NewMinWiseFromSignature(h1, h2, sig)
	} else {
		m = NewMinWise(h1, h2, 5)
	}

	return m
}

func TestLen(t *testing.T) {
	s := makeSigOfInts()
	assert.Equal(t, 400, s.Len())
}

func TestPush(t *testing.T) {
	// Test that 0 values are never pushed.
	h := func(bs []byte) uint64 { return 0 }
	s := NewMinWise(h, h, 2)
	s.Push(1)
	assert.Equal(t, []uint64{infinity, infinity}, s.Signature())
	assert.True(t, s.IsEmpty())
}

func TestCardinality(t *testing.T) {

	sigInts := makeSigOfInts()   // I
	sigEvens := makeSigOfEvens() // E
	sigOdds := makeSigOfOdds()   // O

	// Empty signature should have cardinality 0.
	assert.Equal(t, 0, NewMinWise(h1, h2, 400).Cardinality())

	// Zero signature should also have cardinality 0.
	assert.Equal(t, 0, NewMinWiseFromSignature(h1, h2, []uint64{0, 0, 0}).Cardinality())

	assert.Equal(t, 11001, sigInts.Cardinality())
	assert.Equal(t, 0, sigEvens.IntersectionCardinality(sigOdds))
	assert.Equal(t, sigInts.Cardinality(), sigEvens.UnionCardinality(sigOdds))

	log.Println("Cardinality of Ints:", sigInts.Cardinality())
	log.Println("Cardinality of Evens:", sigEvens.Cardinality())
	log.Println("Cardinality of Odds:", sigOdds.Cardinality())
	log.Println("Cardinality of union:", sigEvens.UnionCardinality(sigOdds))
	log.Println("Cardinality of Ints && Evens:", sigInts.IntersectionCardinality(sigEvens))
	log.Println("Cardinality of Evens && Odds:", sigEvens.IntersectionCardinality(sigOdds))
	log.Println("Cardinality of Ints - Evens:", sigInts.LessCardinality(sigEvens))

}

func TestCopy(t *testing.T) {
	c := makeSigOfEvens().Copy()
	odds := makeSigOfOdds()
	c.Merge(odds)
	log.Println("Cardinality of Evens Copy && Odds:", c.Cardinality())
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
