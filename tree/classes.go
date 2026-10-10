// Package tree provides DecisionTreeClassifier and DecisionTreeRegressor.
package tree

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

// DecisionTreeClassifier is a decision tree classifier.
type DecisionTreeClassifier struct {
	// Parameters
	Criterion             string
	Splitter              string
	MaxDepth              int
	MinSamplesSplit       int
	MinSamplesLeaf        int
	MinWeightFractionLeaf float64
	MaxFeatures           interface{} // int, float64, or string ("sqrt", "log2")
	RandomState           int
	MaxLeafNodes          int
	MinImpurityDecrease   float64
	ClassWeight           string
	CCPAlpha              float64

	// Fitted attributes
	Tree      *Tree
	NFeatures int
	NOutputs  int
	NClasses  int
	Classes   []int
	IsFitted  bool
}

// NewDecisionTreeClassifier creates a new decision tree classifier with default parameters.
func NewDecisionTreeClassifier() *DecisionTreeClassifier {
	return &DecisionTreeClassifier{
		Criterion:             "gini",
		Splitter:              "best",
		MaxDepth:              -1,
		MinSamplesSplit:       2,
		MinSamplesLeaf:        1,
		MinWeightFractionLeaf: 0.0,
		MaxFeatures:           nil,
		RandomState:           -1,
		MaxLeafNodes:          -1,
		MinImpurityDecrease:   0.0,
		ClassWeight:           "",
		CCPAlpha:              0.0,
	}
}

// DecisionTreeClassifierParams contains parameters for DecisionTreeClassifier.
type DecisionTreeClassifierParams struct {
	Criterion             string
	Splitter              string
	MaxDepth              int
	MinSamplesSplit       int
	MinSamplesLeaf        int
	MinWeightFractionLeaf float64
	RandomState           int
	MinImpurityDecrease   float64
}

// NewDecisionTreeClassifierWithParams creates a classifier with custom parameters.
func NewDecisionTreeClassifierWithParams(params DecisionTreeClassifierParams) *DecisionTreeClassifier {
	splitter := params.Splitter
	if splitter == "" {
		splitter = "best"
	}
	return &DecisionTreeClassifier{
		Criterion:             params.Criterion,
		Splitter:              splitter,
		MaxDepth:              params.MaxDepth,
		MinSamplesSplit:       params.MinSamplesSplit,
		MinSamplesLeaf:        params.MinSamplesLeaf,
		MinWeightFractionLeaf: params.MinWeightFractionLeaf,
		RandomState:           params.RandomState,
		MinImpurityDecrease:   params.MinImpurityDecrease,
	}
}

// Fit trains the decision tree classifier.
func (dt *DecisionTreeClassifier) Fit(X [][]float64, y []float64, sampleWeight []float64) error {
	if len(X) == 0 || len(y) == 0 {
		return errors.New("empty training data")
	}
	if len(X) != len(y) {
		return errors.New("X and y have inconsistent dimensions")
	}

	nSamples := len(X)
	dt.NFeatures = len(X[0])

	// Convert y to multi-output format
	yMulti := make([][]float64, nSamples)
	for i := 0; i < nSamples; i++ {
		yMulti[i] = []float64{y[i]}
	}

	// Determine number of classes
	classSet := make(map[int]bool)
	for i := 0; i < nSamples; i++ {
		classSet[int(y[i])] = true
	}

	dt.NClasses = len(classSet)
	dt.Classes = make([]int, 0, dt.NClasses)
	for class := range classSet {
		dt.Classes = append(dt.Classes, class)
	}

	// Sort classes for consistent ordering
	sortInts(dt.Classes)

	// Create criterion
	var criterion Criterion
	switch dt.Criterion {
	case "gini":
		criterion = NewGini(dt.NClasses)
	case "entropy", "log_loss":
		criterion = NewEntropy(dt.NClasses)
	default:
		return fmt.Errorf("unknown criterion: %s", dt.Criterion)
	}

	// Create splitter
	var splitter Splitter
	switch dt.Splitter {
	case "best":
		splitter = NewBestSplitter(dt.RandomState)
	case "random":
		maxFeatures := dt.NFeatures
		if dt.MaxFeatures != nil {
			switch v := dt.MaxFeatures.(type) {
			case int:
				maxFeatures = v
			case float64:
				maxFeatures = int(v * float64(dt.NFeatures))
			case string:
				switch v {
				case "sqrt":
					maxFeatures = int(math.Sqrt(float64(dt.NFeatures)))
				case "log2":
					maxFeatures = int(math.Log2(float64(dt.NFeatures)))
				}
			}
		}
		splitter = NewRandomSplitter(maxFeatures, dt.RandomState)
	default:
		return fmt.Errorf("unknown splitter: %s", dt.Splitter)
	}

	// Build tree
	params := TreeParams{
		MaxDepth:            dt.MaxDepth,
		MinSamplesSplit:     dt.MinSamplesSplit,
		MinSamplesLeaf:      dt.MinSamplesLeaf,
		MinImpurityDecrease: dt.MinImpurityDecrease,
		NClasses:            dt.NClasses,
	}

	tree, err := BuildTree(X, yMulti, sampleWeight, criterion, splitter, params)
	if err != nil {
		return fmt.Errorf("failed to build tree: %w", err)
	}

	dt.Tree = tree
	dt.IsFitted = true
	dt.NOutputs = 1

	return nil
}

// Predict predicts class labels for samples in X.
func (dt *DecisionTreeClassifier) Predict(X [][]float64) ([]int, error) {
	if !dt.IsFitted {
		return nil, errors.New("tree is not fitted")
	}

	proba, err := dt.PredictProba(X)
	if err != nil {
		return nil, err
	}

	predictions := make([]int, len(X))
	for i := 0; i < len(X); i++ {
		// Find class with highest probability
		maxProba := -1.0
		maxClass := 0
		for j, p := range proba[i] {
			if p > maxProba {
				maxProba = p
				maxClass = j
			}
		}
		predictions[i] = dt.Classes[maxClass]
	}

	return predictions, nil
}

// PredictProba predicts class probabilities for samples in X.
func (dt *DecisionTreeClassifier) PredictProba(X [][]float64) ([][]float64, error) {
	if !dt.IsFitted {
		return nil, errors.New("tree is not fitted")
	}

	// Validate input
	for i, row := range X {
		if len(row) != dt.NFeatures {
			return nil, fmt.Errorf("row %d has %d features, expected %d", i, len(row), dt.NFeatures)
		}
	}

	probabilities := make([][]float64, len(X))
	for i := 0; i < len(X); i++ {
		value := dt.Tree.Predict(X[i])
		probabilities[i] = make([]float64, dt.NClasses)
		copy(probabilities[i], value)
	}

	return probabilities, nil
}

// Score returns the mean accuracy on the given test data.
func (dt *DecisionTreeClassifier) Score(X [][]float64, y []float64) (float64, error) {
	predictions, err := dt.Predict(X)
	if err != nil {
		return 0, err
	}

	if len(predictions) != len(y) {
		return 0, errors.New("predictions and y have different lengths")
	}

	correct := 0
	for i := 0; i < len(y); i++ {
		if predictions[i] == int(y[i]) {
			correct++
		}
	}

	return float64(correct) / float64(len(y)), nil
}

// GetDepth returns the depth of the tree.
func (dt *DecisionTreeClassifier) GetDepth() int {
	if dt.Tree == nil {
		return 0
	}
	return dt.Tree.MaxDepth
}

// GetNLeaves returns the number of leaves.
func (dt *DecisionTreeClassifier) GetNLeaves() int {
	if dt.Tree == nil {
		return 0
	}
	return dt.Tree.NLeaves()
}

// DecisionTreeRegressor is a decision tree regressor.
type DecisionTreeRegressor struct {
	// Parameters
	Criterion             string
	Splitter              string
	MaxDepth              int
	MinSamplesSplit       int
	MinSamplesLeaf        int
	MinWeightFractionLeaf float64
	MaxFeatures           interface{}
	RandomState           int
	MaxLeafNodes          int
	MinImpurityDecrease   float64
	CCPAlpha              float64

	// Fitted attributes
	Tree      *Tree
	NFeatures int
	NOutputs  int
	IsFitted  bool
}

// NewDecisionTreeRegressor creates a new decision tree regressor with default parameters.
func NewDecisionTreeRegressor() *DecisionTreeRegressor {
	return &DecisionTreeRegressor{
		Criterion:             "squared_error",
		Splitter:              "best",
		MaxDepth:              -1,
		MinSamplesSplit:       2,
		MinSamplesLeaf:        1,
		MinWeightFractionLeaf: 0.0,
		MaxFeatures:           nil,
		RandomState:           -1,
		MaxLeafNodes:          -1,
		MinImpurityDecrease:   0.0,
		CCPAlpha:              0.0,
	}
}

// DecisionTreeRegressorParams contains parameters for DecisionTreeRegressor.
type DecisionTreeRegressorParams struct {
	Criterion             string
	Splitter              string
	MaxDepth              int
	MinSamplesSplit       int
	MinSamplesLeaf        int
	MinWeightFractionLeaf float64
	RandomState           int
	MinImpurityDecrease   float64
}

// NewDecisionTreeRegressorWithParams creates a regressor with custom parameters.
func NewDecisionTreeRegressorWithParams(params DecisionTreeRegressorParams) *DecisionTreeRegressor {
	splitter := params.Splitter
	if splitter == "" {
		splitter = "best"
	}
	return &DecisionTreeRegressor{
		Criterion:             params.Criterion,
		Splitter:              splitter,
		MaxDepth:              params.MaxDepth,
		MinSamplesSplit:       params.MinSamplesSplit,
		MinSamplesLeaf:        params.MinSamplesLeaf,
		MinWeightFractionLeaf: params.MinWeightFractionLeaf,
		RandomState:           params.RandomState,
		MinImpurityDecrease:   params.MinImpurityDecrease,
	}
}

// Fit trains the decision tree regressor.
func (dt *DecisionTreeRegressor) Fit(X [][]float64, y []float64, sampleWeight []float64) error {
	if len(X) == 0 || len(y) == 0 {
		return errors.New("empty training data")
	}
	if len(X) != len(y) {
		return errors.New("X and y have inconsistent dimensions")
	}

	nSamples := len(X)
	dt.NFeatures = len(X[0])

	// Convert y to multi-output format
	yMulti := make([][]float64, nSamples)
	for i := 0; i < nSamples; i++ {
		yMulti[i] = []float64{y[i]}
	}

	// Create criterion
	var criterion Criterion
	switch dt.Criterion {
	case "squared_error":
		criterion = NewMSE()
	default:
		return fmt.Errorf("unknown criterion: %s", dt.Criterion)
	}

	// Create splitter
	var splitter Splitter
	switch dt.Splitter {
	case "best":
		splitter = NewBestSplitter(dt.RandomState)
	case "random":
		maxFeatures := dt.NFeatures
		if dt.MaxFeatures != nil {
			switch v := dt.MaxFeatures.(type) {
			case int:
				maxFeatures = v
			case float64:
				maxFeatures = int(v * float64(dt.NFeatures))
			case string:
				switch v {
				case "sqrt":
					maxFeatures = int(math.Sqrt(float64(dt.NFeatures)))
				case "log2":
					maxFeatures = int(math.Log2(float64(dt.NFeatures)))
				}
			}
		}
		splitter = NewRandomSplitter(maxFeatures, dt.RandomState)
	default:
		return fmt.Errorf("unknown splitter: %s", dt.Splitter)
	}

	// Build tree
	params := TreeParams{
		MaxDepth:            dt.MaxDepth,
		MinSamplesSplit:     dt.MinSamplesSplit,
		MinSamplesLeaf:      dt.MinSamplesLeaf,
		MinImpurityDecrease: dt.MinImpurityDecrease,
		NClasses:            0, // Regression
	}

	tree, err := BuildTree(X, yMulti, sampleWeight, criterion, splitter, params)
	if err != nil {
		return fmt.Errorf("failed to build tree: %w", err)
	}

	dt.Tree = tree
	dt.IsFitted = true
	dt.NOutputs = 1

	return nil
}

// Predict predicts target values for samples in X.
func (dt *DecisionTreeRegressor) Predict(X [][]float64) ([]float64, error) {
	if !dt.IsFitted {
		return nil, errors.New("tree is not fitted")
	}

	// Validate input
	for i, row := range X {
		if len(row) != dt.NFeatures {
			return nil, fmt.Errorf("row %d has %d features, expected %d", i, len(row), dt.NFeatures)
		}
	}

	predictions := make([]float64, len(X))
	for i := 0; i < len(X); i++ {
		value := dt.Tree.Predict(X[i])
		if len(value) > 0 {
			predictions[i] = value[0]
		}
	}

	return predictions, nil
}

// Score returns the coefficient of determination R^2.
func (dt *DecisionTreeRegressor) Score(X [][]float64, y []float64) (float64, error) {
	predictions, err := dt.Predict(X)
	if err != nil {
		return 0, err
	}

	if len(predictions) != len(y) {
		return 0, errors.New("predictions and y have different lengths")
	}

	// Calculate R^2
	var yMean float64
	for i := 0; i < len(y); i++ {
		yMean += y[i]
	}
	yMean /= float64(len(y))

	var ssTot, ssRes float64
	for i := 0; i < len(y); i++ {
		diffTot := y[i] - yMean
		ssTot += diffTot * diffTot
		diffRes := y[i] - predictions[i]
		ssRes += diffRes * diffRes
	}

	if ssTot == 0 {
		return 1.0, nil
	}

	return 1.0 - ssRes/ssTot, nil
}

// GetDepth returns the depth of the tree.
func (dt *DecisionTreeRegressor) GetDepth() int {
	if dt.Tree == nil {
		return 0
	}
	return dt.Tree.MaxDepth
}

// GetNLeaves returns the number of leaves.
func (dt *DecisionTreeRegressor) GetNLeaves() int {
	if dt.Tree == nil {
		return 0
	}
	return dt.Tree.NLeaves()
}

// sortInts sorts a slice of ints using the standard library sort.
func sortInts(data []int) {
	sort.Ints(data)
}
