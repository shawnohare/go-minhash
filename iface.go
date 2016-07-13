package minhash

type HashFunc func([]byte) uint64

// Interface is an a probabilistic data structure used to
// compute a similarity preserving signature for a set.  It ingests
// a stream of the set's elements and continuously updates the signature.
type Interface interface {
	// Signature returns the signature itself.
	Signature() []uint64
}

// MinHashSimilarity computes an estimate for the
// Jaccard similarity of two sets given their MinHash signatures.
func Similarity(x, y Interface) float64 {
	return similarity(x.Signature(), y.Signature())
}

func IsEmpty(m Interface) bool {
	var empty = true

	if m != nil {
		// Check whether each minimum is infinite.
		for _, v := range m.Signature() {
			if v < infinity {
				empty = false
				break
			}
		}
	}
	return empty
}
