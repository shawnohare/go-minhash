package mhsig

import (
	"math"
	"sync"
)

// MinWise is a data structure for generating a parametric family of
// hash functions of the form h1 + i*h2 for i=1, ..., k to compute
// a MinHash signatures.
type MinWise struct {
	// size is the signature length the instance will compute
	size int
	h1 HashFunc
	h2 HashFunc
}

func (m *MinWise) Sketch(xs Set) Signature {
	mins := defaultSignature(m.size)

	// map m.h1 and m.h2 over input set xs.
	// These are temporarily stored so that h1 and h2 do not
	// have to be repeatedly computed for each x in xs.
	h1s := make(Signature, len(xs))
	h2s := make(Signature, len(xs))
	for i, x := range xs {
		h1s[i] = m.h1(x)
		h2s[i] = m.h2(x)
	}
	// Determine minim values for the hash functions in parallel. 
	var wg sync.WaitGroup()
	wg.Add(m.size)
	for i, m := range mins {
		go func(i int) {
			defer wg.Done()
			for j := range xs {
				// Compute ith hash function on jth data point
				hx := h1s[j] + SignatureElement(i) * h2s[j]
				if hx < m {
					mins[i] = hx
				}
			}
		}
	}
	wg.Wait()
	
	return mins
}




func (m *MinWise) Similarity(sig1, sig2 Signature) float64 {
	if len(sig1) != len(sig2) {
		panic("signature size mismatch")
	}

	intersect := 0 // counter for number of elements in intersection
	
	for i  := range sig1 {
		if sig1[i] = sig2[i] {
			intersect ++
		}
	}

	return float64(intersect) / float64(len(sig1))
}
