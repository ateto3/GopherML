package utils

import (
	"math/rand"
	"time"
)

// NewRNG creates a new random number generator. If seed is -1, it uses current time.
func NewRNG(seed int64) *rand.Rand {
	if seed == -1 {
		return rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return rand.New(rand.NewSource(seed))
}

// Shuffle shuffles the indices of a slice.
func Shuffle(n int, rng *rand.Rand) []int {
	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}
	rng.Shuffle(n, func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})
	return indices
}
