// Package tree provides criteria for evaluating splits in decision trees.
package tree

import (
	"math"
)

// Criterion evaluates the quality of a split.
type Criterion interface {
	// Init initializes the criterion with targets, sample weights, and indices.
	Init(y [][]float64, sampleWeight []float64, indices []int, start, end int)
	// Reset resets the criterion for a new node.
	Reset(start, end int)
	// Update updates the criterion with a new sample.
	Update(idx int)
	// ProxyImpurityImprovement calculates a proxy for impurity improvement.
	ProxyImpurityImprovement() float64
	// ImpurityImprovement calculates the actual impurity improvement.
	ImpurityImprovement(parentImpurity float64) float64
	// NodeValue returns the value for the current node.
	NodeValue() []float64
	// NodeImpurity returns the impurity for the current node.
	NodeImpurity() float64
	// ChildrenImpurity returns the impurities for left and right children.
	ChildrenImpurity() (left, right float64)
}

// Gini implements the Gini impurity criterion for classification.
type Gini struct {
	y            [][]float64
	sampleWeight []float64
	indices      []int
	nClasses     int
	start, end   int

	totalCounts []float64
	totalN      float64
	leftCounts  []float64
	leftN       float64
}

// NewGini creates a new Gini criterion.
func NewGini(nClasses int) *Gini {
	return &Gini{
		nClasses:    nClasses,
		totalCounts: make([]float64, nClasses),
		leftCounts:  make([]float64, nClasses),
	}
}

// Init initializes the criterion.
func (g *Gini) Init(y [][]float64, sampleWeight []float64, indices []int, start, end int) {
	g.y = y
	g.sampleWeight = sampleWeight
	g.indices = indices
	g.Reset(start, end)
}

// Reset resets for a new node.
func (g *Gini) Reset(start, end int) {
	g.start = start
	g.end = end

	// Grow count slices if any class label in this node exceeds current allocation.
	// This can happen when a bootstrap sample drops some classes, causing nClasses
	// to be smaller than the maximum class label in the full dataset.
	for i := start; i < end; i++ {
		dataIdx := g.indices[i]
		classIdx := int(g.y[dataIdx][0])
		if classIdx >= g.nClasses {
			g.nClasses = classIdx + 1
			g.totalCounts = append(g.totalCounts, make([]float64, g.nClasses-len(g.totalCounts))...)
			g.leftCounts = append(g.leftCounts, make([]float64, g.nClasses-len(g.leftCounts))...)
		}
	}

	for i := 0; i < g.nClasses; i++ {
		g.totalCounts[i] = 0
		g.leftCounts[i] = 0
	}
	g.totalN = 0
	g.leftN = 0

	for i := start; i < end; i++ {
		dataIdx := g.indices[i]
		classIdx := int(g.y[dataIdx][0])
		weight := 1.0
		if g.sampleWeight != nil {
			weight = g.sampleWeight[dataIdx]
		}
		g.totalCounts[classIdx] += weight
		g.totalN += weight
	}
}

// Update updates with a new sample (moves it to the left).
func (g *Gini) Update(idx int) {
	classIdx := int(g.y[idx][0])
	weight := 1.0
	if g.sampleWeight != nil {
		weight = g.sampleWeight[idx]
	}
	g.leftCounts[classIdx] += weight
	g.leftN += weight
}

// NodeValue returns the class distribution for the total node.
func (g *Gini) NodeValue() []float64 {
	if g.totalN == 0 {
		return make([]float64, g.nClasses)
	}
	result := make([]float64, g.nClasses)
	for i := 0; i < g.nClasses; i++ {
		result[i] = g.totalCounts[i] / g.totalN
	}
	return result
}

// NodeImpurity calculates Gini impurity for the total node.
func (g *Gini) NodeImpurity() float64 {
	return g.calculateGini(g.totalCounts, g.totalN)
}

func (g *Gini) calculateGini(counts []float64, n float64) float64 {
	if n == 0 {
		return 0
	}
	impurity := 1.0
	for i := 0; i < g.nClasses; i++ {
		p := counts[i] / n
		impurity -= p * p
	}
	return impurity
}

// ProxyImpurityImprovement returns weighted impurity of children (to be minimized).
func (g *Gini) ProxyImpurityImprovement() float64 {
	leftImpurity := g.calculateGini(g.leftCounts, g.leftN)
	
	rightN := g.totalN - g.leftN
	rightCounts := make([]float64, g.nClasses)
	for i := 0; i < g.nClasses; i++ {
		rightCounts[i] = g.totalCounts[i] - g.leftCounts[i]
	}
	rightImpurity := g.calculateGini(rightCounts, rightN)

	return (g.leftN*leftImpurity + rightN*rightImpurity) / g.totalN
}

// ImpurityImprovement calculates actual improvement (to be maximized).
func (g *Gini) ImpurityImprovement(parentImpurity float64) float64 {
	return parentImpurity - g.ProxyImpurityImprovement()
}

// ChildrenImpurity returns impurities for children.
func (g *Gini) ChildrenImpurity() (left, right float64) {
	left = g.calculateGini(g.leftCounts, g.leftN)
	
	rightN := g.totalN - g.leftN
	rightCounts := make([]float64, g.nClasses)
	for i := 0; i < g.nClasses; i++ {
		rightCounts[i] = g.totalCounts[i] - g.leftCounts[i]
	}
	right = g.calculateGini(rightCounts, rightN)
	return
}

// Entropy implements the Entropy criterion for classification.
type Entropy struct {
	y            [][]float64
	sampleWeight []float64
	indices      []int
	nClasses     int
	start, end   int

	totalCounts []float64
	totalN      float64
	leftCounts  []float64
	leftN       float64
}

// NewEntropy creates a new Entropy criterion.
func NewEntropy(nClasses int) *Entropy {
	return &Entropy{
		nClasses:    nClasses,
		totalCounts: make([]float64, nClasses),
		leftCounts:  make([]float64, nClasses),
	}
}

// Init initializes the criterion.
func (e *Entropy) Init(y [][]float64, sampleWeight []float64, indices []int, start, end int) {
	e.y = y
	e.sampleWeight = sampleWeight
	e.indices = indices
	e.Reset(start, end)
}

// Reset resets for a new node.
func (e *Entropy) Reset(start, end int) {
	e.start = start
	e.end = end

	// Grow count slices if any class label in this node exceeds current allocation.
	for i := start; i < end; i++ {
		dataIdx := e.indices[i]
		classIdx := int(e.y[dataIdx][0])
		if classIdx >= e.nClasses {
			e.nClasses = classIdx + 1
			e.totalCounts = append(e.totalCounts, make([]float64, e.nClasses-len(e.totalCounts))...)
			e.leftCounts = append(e.leftCounts, make([]float64, e.nClasses-len(e.leftCounts))...)
		}
	}

	for i := 0; i < e.nClasses; i++ {
		e.totalCounts[i] = 0
		e.leftCounts[i] = 0
	}
	e.totalN = 0
	e.leftN = 0

	for i := start; i < end; i++ {
		dataIdx := e.indices[i]
		classIdx := int(e.y[dataIdx][0])
		weight := 1.0
		if e.sampleWeight != nil {
			weight = e.sampleWeight[dataIdx]
		}
		e.totalCounts[classIdx] += weight
		e.totalN += weight
	}
}

// Update updates with a new sample (moves it to the left).
func (e *Entropy) Update(idx int) {
	classIdx := int(e.y[idx][0])
	weight := 1.0
	if e.sampleWeight != nil {
		weight = e.sampleWeight[idx]
	}
	e.leftCounts[classIdx] += weight
	e.leftN += weight
}

// NodeValue returns the class distribution for the total node.
func (e *Entropy) NodeValue() []float64 {
	if e.totalN == 0 {
		return make([]float64, e.nClasses)
	}
	result := make([]float64, e.nClasses)
	for i := 0; i < e.nClasses; i++ {
		result[i] = e.totalCounts[i] / e.totalN
	}
	return result
}

// NodeImpurity calculates Entropy for the total node.
func (e *Entropy) NodeImpurity() float64 {
	return e.calculateEntropy(e.totalCounts, e.totalN)
}

func (e *Entropy) calculateEntropy(counts []float64, n float64) float64 {
	if n == 0 {
		return 0
	}
	entropy := 0.0
	for i := 0; i < e.nClasses; i++ {
		p := counts[i] / n
		if p > 0 {
			entropy -= p * math.Log2(p)
		}
	}
	return entropy
}

// ProxyImpurityImprovement returns weighted impurity of children.
func (e *Entropy) ProxyImpurityImprovement() float64 {
	leftImpurity := e.calculateEntropy(e.leftCounts, e.leftN)
	
	rightN := e.totalN - e.leftN
	rightCounts := make([]float64, e.nClasses)
	for i := 0; i < e.nClasses; i++ {
		rightCounts[i] = e.totalCounts[i] - e.leftCounts[i]
	}
	rightImpurity := e.calculateEntropy(rightCounts, rightN)

	return (e.leftN*leftImpurity + rightN*rightImpurity) / e.totalN
}

// ImpurityImprovement calculates actual improvement.
func (e *Entropy) ImpurityImprovement(parentImpurity float64) float64 {
	return parentImpurity - e.ProxyImpurityImprovement()
}

// ChildrenImpurity returns impurities for children.
func (e *Entropy) ChildrenImpurity() (left, right float64) {
	left = e.calculateEntropy(e.leftCounts, e.leftN)
	
	rightN := e.totalN - e.leftN
	rightCounts := make([]float64, e.nClasses)
	for i := 0; i < e.nClasses; i++ {
		rightCounts[i] = e.totalCounts[i] - e.leftCounts[i]
	}
	right = e.calculateEntropy(rightCounts, rightN)
	return
}

// MSE implements Mean Squared Error criterion for regression.
type MSE struct {
	y            [][]float64
	sampleWeight []float64
	indices      []int
	start, end   int

	totalSum   float64
	totalSumSq float64
	totalN     float64
	leftSum    float64
	leftSumSq  float64
	leftN      float64
}

// NewMSE creates a new MSE criterion.
func NewMSE() *MSE {
	return &MSE{}
}

// Init initializes the criterion.
func (m *MSE) Init(y [][]float64, sampleWeight []float64, indices []int, start, end int) {
	m.y = y
	m.sampleWeight = sampleWeight
	m.indices = indices
	m.Reset(start, end)
}

// Reset resets for a new node.
func (m *MSE) Reset(start, end int) {
	m.start = start
	m.end = end

	m.totalSum = 0
	m.totalSumSq = 0
	m.totalN = 0
	m.leftSum = 0
	m.leftSumSq = 0
	m.leftN = 0

	for i := start; i < end; i++ {
		dataIdx := m.indices[i]
		weight := 1.0
		if m.sampleWeight != nil {
			weight = m.sampleWeight[dataIdx]
		}
		val := m.y[dataIdx][0]
		m.totalSum += weight * val
		m.totalSumSq += weight * val * val
		m.totalN += weight
	}
}

// Update updates with a new sample (moves it to the left).
func (m *MSE) Update(idx int) {
	weight := 1.0
	if m.sampleWeight != nil {
		weight = m.sampleWeight[idx]
	}
	val := m.y[idx][0]
	m.leftSum += weight * val
	m.leftSumSq += weight * val * val
	m.leftN += weight
}

// NodeValue returns the mean for the total node.
func (m *MSE) NodeValue() []float64 {
	if m.totalN == 0 {
		return []float64{0}
	}
	return []float64{m.totalSum / m.totalN}
}

// NodeImpurity calculates MSE for the total node.
func (m *MSE) NodeImpurity() float64 {
	return m.calculateMSE(m.totalSum, m.totalSumSq, m.totalN)
}

func (m *MSE) calculateMSE(sum, sumSq, n float64) float64 {
	if n == 0 {
		return 0
	}
	mean := sum / n
	return sumSq/n - mean*mean
}

// ProxyImpurityImprovement returns weighted impurity of children.
func (m *MSE) ProxyImpurityImprovement() float64 {
	leftImpurity := m.calculateMSE(m.leftSum, m.leftSumSq, m.leftN)
	
	rightN := m.totalN - m.leftN
	rightSum := m.totalSum - m.leftSum
	rightSumSq := m.totalSumSq - m.leftSumSq
	rightImpurity := m.calculateMSE(rightSum, rightSumSq, rightN)

	return (m.leftN*leftImpurity + rightN*rightImpurity) / m.totalN
}

// ImpurityImprovement calculates actual improvement.
func (m *MSE) ImpurityImprovement(parentImpurity float64) float64 {
	return parentImpurity - m.ProxyImpurityImprovement()
}

// ChildrenImpurity returns impurities for children.
func (m *MSE) ChildrenImpurity() (left, right float64) {
	left = m.calculateMSE(m.leftSum, m.leftSumSq, m.leftN)
	
	rightN := m.totalN - m.leftN
	rightSum := m.totalSum - m.leftSum
	rightSumSq := m.totalSumSq - m.leftSumSq
	right = m.calculateMSE(rightSum, rightSumSq, rightN)
	return
}
