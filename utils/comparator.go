package utils

const (
	Lt = -1
	Et = 0
	Gt = 1
)

// Comparator Comparator
type Comparator interface {
	// k1 > k2, return 1
	// k2 = k2, return 0
	// k1 < k2, return -1
	Compare(k1, k2 interface{}) int
}
