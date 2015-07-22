package minhash

import "math"

const (
	MaxUint  uint   = ^uint(0)
	MaxInt   int    = int(MaxUint >> 1)
	infinity uint64 = math.MaxUint64
)

func cardinality(sig []uint64) int {
	// http://www.cohenwang.com/edith/Papers/tcest.pdf
	// Let h be a hash function.
	// With M as the the maximum uint64, the value
	// z := -ln( (M - min h(A)) / M)
	// is a sample from Z := min { X1 ~ Exp(1), ... , Xn ~ Exp(1)}.
	// where n := |A|.  But Z ~ Exp(n) has mean 1/n.
	// Thus the inverse of the mean of all the transformed signature values
	// is an estimate for |A|.
	var card int
	sum := 0.0 // running sum

	// Use inverse transform to obtain a sample from Z for each sig element.
	for _, v := range sig {

		d := infinity - v
		// Make sure the zero signature and sig for empty set return 0.
		if d == 0 || d == infinity {
			continue
		} else {
			sum += -math.Log(float64(d) / float64(infinity))
		}

	}

	if sum != 0.0 {
		card = int(float64(len(sig)) / sum)
	}

	return card
}

// union computes the signature for the union of two sets
// S and R given their signatures s and r.
func union(s, r []uint64) []uint64 {
	if len(s) != len(r) {
		panic("Signature length mismatch.")
	}

	res := make([]uint64, len(s))
	for i := range s {
		res[i] = min(s[i], r[i])
	}
	return res
}

// compute the similarity between two MinHash signatures.
func similarity(s, r []uint64) float64 {
	if len(s) != len(r) {
		panic("Signature size mismatch.")
	}

	intersect := 0 // counter for number of elements in intersection

	for i := range s {
		if s[i] == r[i] {
			intersect++
		}
	}

	return float64(intersect) / float64(len(s))
}

// emptySignature will return the signature for the empty set.
func emptySetSignature(size int) []uint64 {
	s := make([]uint64, size)
	for i := range s {
		s[i] = infinity
	}
	return s
}

func min(x, y uint64) uint64 {
	if x <= y {
		return x
	}
	return y
}
