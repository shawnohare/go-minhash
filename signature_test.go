package minhash

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnion(t *testing.T) {
	xs := []uint64{1, 2, 3}
	ys := []uint64{0, 3, 1}
	expected := []uint64{0, 2, 1}
	assert.Equal(t, expected, union(xs, ys))
}
