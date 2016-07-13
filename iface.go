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

// Similarity computes the similarity of two signatures represented
// as MinHash instances.  This estimates the Jaccard index of the
// two underlying sets.
func (m *MinHash) Similarity(m2 Interface) float64 {
	return Similarity(m, m2)
}

// Cardinality estimates the cardinality of the set. Both the
// signature for the empty set and the zero signature have 0 cardinality.
func (m *MinHash) Cardinality() int {
	return Cardinality(m)
}

// UnionCardinality estimates the cardinality of the union.
func (m *MinHash) UnionCardinality(m2 Interface) int {
	return UnionCardinality(m, m2)
}

// IntersectionCardinality estimates the cardinality of the intersection.
func (m *MinHash) IntersectionCardinality(m2 Interface) int {
	return IntersectionCardinality(m, m2)
}

// SymmetricDifferenceCardinality estimates the difference between
// the cardinality of the union and intersection.
func (m *MinHash) SymmetricDifferenceCardinality(m2 Interface) int {
	return SymmetricDifferenceCardinality(m, m2)
}

// LessCardinality estimates the cardinality of the left set minus
// the right set. This operator is not symmetric.
func (m *MinHash) LessCardinality(m2 Interface) int {
	return LessCardinality(m, m2)
}
