package ensemble

import (
	"math/rand"
)

// GenerateBootstrapIndices generates n samples with replacement from 0..n-1.
func GenerateBootstrapIndices(n int, rng *rand.Rand) []int {
	indices := make([]int, n)
	for i := 0; i < n; i++ {
		indices[i] = rng.Intn(n)
	}
	return indices
}

// GetBootstrapSample returns a bootstrap sample of X and y.
func GetBootstrapSample(X [][]float64, y []float64, rng *rand.Rand) ([][]float64, []float64) {
	n := len(X)
	indices := GenerateBootstrapIndices(n, rng)

	XSample := make([][]float64, n)
	ySample := make([]float64, n)

	for i, idx := range indices {
		XSample[i] = X[idx]
		ySample[i] = y[idx]
	}

	return XSample, ySample
}
