package model_selection

import (
	"github.com/ateto3/gopherml/utils"
)

// KFold provides train/test indices to split data in train/test sets.
type KFold struct {
	NSplits     int
	Shuffle     bool
	RandomState int64
}

// SplitRecord stores indices for one fold.
type SplitRecord struct {
	TrainIndices []int
	TestIndices  []int
}

// Split generates indices to split data into training and test set.
func (k *KFold) Split(nSamples int) []SplitRecord {
	indices := make([]int, nSamples)
	for i := range indices {
		indices[i] = i
	}

	if k.Shuffle {
		rng := utils.NewRNG(k.RandomState)
		rng.Shuffle(nSamples, func(i, j int) {
			indices[i], indices[j] = indices[j], indices[i]
		})
	}

	foldSize := nSamples / k.NSplits
	records := make([]SplitRecord, k.NSplits)

	for i := 0; i < k.NSplits; i++ {
		start := i * foldSize
		end := start + foldSize
		if i == k.NSplits-1 {
			end = nSamples
		}

		testIndices := make([]int, end-start)
		copy(testIndices, indices[start:end])

		trainIndices := make([]int, 0, nSamples-(end-start))
		trainIndices = append(trainIndices, indices[:start]...)
		trainIndices = append(trainIndices, indices[end:]...)

		records[i] = SplitRecord{
			TrainIndices: trainIndices,
			TestIndices:  testIndices,
		}
	}

	return records
}
