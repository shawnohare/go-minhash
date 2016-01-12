package minhash

// UnionCardinality estimates the cardinality of the union.
func UnionCardinality(m1, m2 Interface) int {
	u := union(m1.Signature(), m2.Signature())
	return cardinality(u)
}

// IntersectionCardinality estimates the cardinality of the intersection.
func IntersectionCardinality(m1, m2 Interface) int {
	// |A & B| + |A || B| = |A| +|B|
	c1 := Cardinality(m1)
	c2 := Cardinality(m2)
	u := UnionCardinality(m1, m2)
	est := c1 + c2 - u
	if est < 0 {
		est = 0
	}
	if est > c1 {
		est = c1
	}
	if est > c2 {
		est = c2
	}

	return est
}

// SymmetricDifferenceCardinality estimates the difference between
// the cardinality of the union and intersection.
func SymmetricDifferenceCardinality(m1, m2 Interface) int {
	est := UnionCardinality(m1, m2) - IntersectionCardinality(m1, m2)

	if est < 0 {
		est = 0
	}

	return est
}

// LessCardinality estimates the cardinality of the left set minus
// the right set. This operator is not symmetric.
func LessCardinality(m1, m2 Interface) int {
	est := Cardinality(m1) - IntersectionCardinality(m1, m2)
	if est < 0 {
		est = 0
	}

	return est
}

// Cardinality estimates the cardinality of the set. Both the
// signature for the empty set and the zero signature have cardinality 0.
func Cardinality(m Interface) int {
	return cardinality(m.Signature())
}
