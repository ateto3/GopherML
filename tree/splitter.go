// Package tree provides split finding algorithms for decision trees.
package tree

import (
	"math"
	"math/rand"
	"sort"
)

// SplitRecord stores information about a potential split.
type SplitRecord struct {
	Feature   int
	Threshold float64
	Impurity  float64
	Pos       int // Position in sorted samples where split occurs
}

// Splitter finds the best split for a node.
type Splitter interface {
	// Init initializes the splitter with data and indices.
	Init(X [][]float64, y [][]float64, sampleWeight []float64, indices []int)
	// NodeSplit finds the best split for the current node.
	NodeSplit(start, end int, criterion Criterion) (split SplitRecord, hasSplit bool)
}

// BestSplitter finds the best split by evaluating all features.
type BestSplitter struct {
	X            [][]float64
	y            [][]float64
	sampleWeight []float64
	indices      []int
	nFeatures    int
	rng          *rand.Rand
}

// NewBestSplitter creates a new BestSplitter.
func NewBestSplitter(randomState int) *BestSplitter {
	var rng *rand.Rand
	if randomState >= 0 {
		rng = rand.New(rand.NewSource(int64(randomState)))
	} else {
		rng = rand.New(rand.NewSource(rand.Int63()))
	}
	return &BestSplitter{rng: rng}
}

// Init initializes the splitter.
func (s *BestSplitter) Init(X [][]float64, y [][]float64, sampleWeight []float64, indices []int) {
	s.X = X
	s.y = y
	s.sampleWeight = sampleWeight
	s.indices = indices
	s.nFeatures = len(X[0])
}

// NodeSplit finds the best split.
func (s *BestSplitter) NodeSplit(start, end int, criterion Criterion) (SplitRecord, bool) {
	if end-start <= 1 {
		return SplitRecord{}, false
	}

	bestImpurity := math.Inf(1)
	bestSplit := SplitRecord{}
	hasSplit := false

	// Try each feature
	for feature := 0; feature < s.nFeatures; feature++ {
		split, found := s.findBestSplitForFeature(start, end, feature, criterion)
		if found && split.Impurity < bestImpurity {
			bestImpurity = split.Impurity
			bestSplit = split
			hasSplit = true
		}
	}

	return bestSplit, hasSplit
}

// findBestSplitForFeature finds the best split for a specific feature.
func (s *BestSplitter) findBestSplitForFeature(start, end, feature int, criterion Criterion) (SplitRecord, bool) {
	nSamples := end - start
	type indexedValue struct {
		idx   int // The index into s.indices
		value float64
	}
	values := make([]indexedValue, nSamples)
	for i := 0; i < nSamples; i++ {
		dataIdx := s.indices[start+i]
		values[i] = indexedValue{idx: start + i, value: s.X[dataIdx][feature]}
	}

	// Sort by feature value
	sort.Slice(values, func(i, j int) bool {
		return values[i].value < values[j].value
	})

	bestImpurity := math.Inf(1)
	bestSplit := SplitRecord{}
	hasSplit := false

	criterion.Reset(start, end)

	for i := 0; i < nSamples-1; i++ {
		// Move sample to left side using its absolute index
		criterion.Update(s.indices[values[i].idx])

		// Check if we can split here (distinct values)
		if values[i].value < values[i+1].value {
			impurity := criterion.ProxyImpurityImprovement()
			if impurity < bestImpurity {
				bestImpurity = impurity
				threshold := (values[i].value + values[i+1].value) / 2
				bestSplit = SplitRecord{
					Feature:   feature,
					Threshold: threshold,
					Impurity:  impurity,
					Pos:       start + i + 1,
				}
				hasSplit = true
			}
		}
	}

	return bestSplit, hasSplit
}

// RandomSplitter finds random splits among the best ones.
type RandomSplitter struct {
	X            [][]float64
	y            [][]float64
	sampleWeight []float64
	indices      []int
	nFeatures    int
	maxFeatures  int
	rng          *rand.Rand
}

// NewRandomSplitter creates a new RandomSplitter.
func NewRandomSplitter(maxFeatures int, randomState int) *RandomSplitter {
	var rng *rand.Rand
	if randomState >= 0 {
		rng = rand.New(rand.NewSource(int64(randomState)))
	} else {
		rng = rand.New(rand.NewSource(rand.Int63()))
	}
	return &RandomSplitter{
		maxFeatures: maxFeatures,
		rng:         rng,
	}
}

// Init initializes the splitter.
func (s *RandomSplitter) Init(X [][]float64, y [][]float64, sampleWeight []float64, indices []int) {
	s.X = X
	s.y = y
	s.sampleWeight = sampleWeight
	s.indices = indices
	s.nFeatures = len(X[0])
}

// NodeSplit finds a random split.
func (s *RandomSplitter) NodeSplit(start, end int, criterion Criterion) (SplitRecord, bool) {
	if end-start <= 1 {
		return SplitRecord{}, false
	}

	// Select random subset of features
	features := s.rng.Perm(s.nFeatures)
	if len(features) > s.maxFeatures {
		features = features[:s.maxFeatures]
	}

	bestImpurity := math.Inf(1)
	bestSplit := SplitRecord{}
	hasSplit := false

	// Try random features
	for _, feature := range features {
		split, found := s.findBestSplitForFeature(start, end, feature, criterion)
		if found && split.Impurity < bestImpurity {
			bestImpurity = split.Impurity
			bestSplit = split
			hasSplit = true
		}
	}

	return bestSplit, hasSplit
}

// findBestSplitForFeature finds the best split for a feature.
func (s *RandomSplitter) findBestSplitForFeature(start, end, feature int, criterion Criterion) (SplitRecord, bool) {
	nSamples := end - start
	type indexedValue struct {
		idx   int
		value float64
	}
	values := make([]indexedValue, nSamples)
	for i := 0; i < nSamples; i++ {
		dataIdx := s.indices[start+i]
		values[i] = indexedValue{idx: start + i, value: s.X[dataIdx][feature]}
	}

	sort.Slice(values, func(i, j int) bool {
		return values[i].value < values[j].value
	})

	bestImpurity := math.Inf(1)
	bestSplit := SplitRecord{}
	hasSplit := false

	criterion.Reset(start, end)

	for i := 0; i < nSamples-1; i++ {
		criterion.Update(s.indices[values[i].idx])

		if values[i].value < values[i+1].value {
			impurity := criterion.ProxyImpurityImprovement()
			if impurity < bestImpurity {
				bestImpurity = impurity
				threshold := (values[i].value + values[i+1].value) / 2
				bestSplit = SplitRecord{
					Feature:   feature,
					Threshold: threshold,
					Impurity:  impurity,
					Pos:       start + i + 1,
				}
				hasSplit = true
			}
		}
	}

	return bestSplit, hasSplit
}
