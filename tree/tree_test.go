package tree

import (
	"math"
	"testing"
)

// Tolerance for floating point comparisons
const epsilon = 1e-6

func TestDecisionTreeClassifierSimple(t *testing.T) {
	// Simple classification problem
	// Class 0: x < 5, Class 1: x >= 5
	X := [][]float64{
		{1}, {2}, {3}, {4}, // Class 0
		{6}, {7}, {8}, {9}, // Class 1
	}
	y := []float64{0, 0, 0, 0, 1, 1, 1, 1}

	clf := NewDecisionTreeClassifier()
	err := clf.Fit(X, y, nil)
	if err != nil {
		t.Fatalf("Fit failed: %v", err)
	}

	// Test predictions
	predictions, err := clf.Predict([][]float64{{0}, {7}, {10}})
	if err != nil {
		t.Fatalf("Predict failed: %v", err)
	}

	// Point at 0 should be class 0
	if predictions[0] != 0 {
		t.Errorf("Expected class 0 for x=0, got %d", predictions[0])
	}
	// Point at 7 should be class 1
	if predictions[1] != 1 {
		t.Errorf("Expected class 1 for x=7, got %d", predictions[1])
	}
	// Point at 10 should be class 1
	if predictions[2] != 1 {
		t.Errorf("Expected class 1 for x=10, got %d", predictions[2])
	}
}

func TestDecisionTreeClassifierIrisLike(t *testing.T) {
	// Simulate a simple 2D classification problem
	X := [][]float64{
		{1, 1}, {1, 2}, {2, 1}, {2, 2}, // Class 0
		{5, 5}, {5, 6}, {6, 5}, {6, 6}, // Class 1
	}
	y := []float64{0, 0, 0, 0, 1, 1, 1, 1}

	clf := NewDecisionTreeClassifierWithParams(DecisionTreeClassifierParams{
		Criterion: "gini",
		MaxDepth:  3,
	})

	err := clf.Fit(X, y, nil)
	if err != nil {
		t.Fatalf("Fit failed: %v", err)
	}

	// Test on training data - should be perfect
	predictions, err := clf.Predict(X)
	if err != nil {
		t.Fatalf("Predict failed: %v", err)
	}

	correct := 0
	for i := 0; i < len(y); i++ {
		if predictions[i] == int(y[i]) {
			correct++
		}
	}

	accuracy := float64(correct) / float64(len(y))
	if accuracy < 1.0 {
		t.Errorf("Expected perfect accuracy on training data, got %.2f", accuracy)
	}
}

func TestDecisionTreeClassifierEntropy(t *testing.T) {
	X := [][]float64{
		{1}, {2}, {3},
		{7}, {8}, {9},
	}
	y := []float64{0, 0, 0, 1, 1, 1}

	clf := NewDecisionTreeClassifierWithParams(DecisionTreeClassifierParams{
		Criterion: "entropy",
		MaxDepth:  2,
	})

	err := clf.Fit(X, y, nil)
	if err != nil {
		t.Fatalf("Fit failed: %v", err)
	}

	predictions, err := clf.Predict([][]float64{{2}, {8}})
	if err != nil {
		t.Fatalf("Predict failed: %v", err)
	}

	if predictions[0] != 0 {
		t.Errorf("Expected class 0, got %d", predictions[0])
	}
	if predictions[1] != 1 {
		t.Errorf("Expected class 1, got %d", predictions[1])
	}
}

func TestDecisionTreeRegressorSimple(t *testing.T) {
	// Simple regression: y = 2x
	X := [][]float64{
		{1}, {2}, {3}, {4}, {5},
	}
	y := []float64{2, 4, 6, 8, 10}

	reg := NewDecisionTreeRegressor()
	err := reg.Fit(X, y, nil)
	if err != nil {
		t.Fatalf("Fit failed: %v", err)
	}

	predictions, err := reg.Predict([][]float64{{1.5}, {3.5}, {6}})
	if err != nil {
		t.Fatalf("Predict failed: %v", err)
	}

	// Predictions should be close to actual values
	if math.Abs(predictions[0]-3) > 1 {
		t.Errorf("Expected prediction ~3 for x=1.5, got %f", predictions[0])
	}
	if math.Abs(predictions[1]-7) > 1 {
		t.Errorf("Expected prediction ~7 for x=3.5, got %f", predictions[1])
	}
}

func TestDecisionTreeRegressorScore(t *testing.T) {
	// Perfect fit should give R^2 = 1.0
	X := [][]float64{
		{1}, {2}, {3}, {4}, {5},
	}
	y := []float64{2, 4, 6, 8, 10}

	reg := NewDecisionTreeRegressor()
	err := reg.Fit(X, y, nil)
	if err != nil {
		t.Fatalf("Fit failed: %v", err)
	}

	score, err := reg.Score(X, y)
	if err != nil {
		t.Fatalf("Score failed: %v", err)
	}

	if score < 0.99 {
		t.Errorf("Expected R^2 close to 1.0 for perfect fit, got %f", score)
	}
}

func TestTreeStructure(t *testing.T) {
	// Test that tree structure is built correctly
	X := [][]float64{
		{1}, {2}, {3}, {4},
	}
	y := []float64{0, 0, 1, 1}

	clf := NewDecisionTreeClassifier()
	err := clf.Fit(X, y, nil)
	if err != nil {
		t.Fatalf("Fit failed: %v", err)
	}

	// Check tree has nodes
	if len(clf.Tree.Nodes) == 0 {
		t.Error("Tree has no nodes")
	}

	// Check at least one leaf
	leaves := clf.Tree.NLeaves()
	if leaves == 0 {
		t.Error("Tree has no leaves")
	}

	t.Logf("Tree has %d nodes and %d leaves, depth=%d",
		len(clf.Tree.Nodes), leaves, clf.GetDepth())
}

func TestGiniCriterion(t *testing.T) {
	// Test Gini impurity calculation
	y := [][]float64{
		{0}, {0}, {0}, {1}, {1}, {1},
	}
	gini := NewGini(2)
	indices := []int{0, 1, 2, 3, 4, 5}
	gini.Init(y, nil, indices, 0, 6)

	// Perfectly balanced: impurity should be 0.5
	impurity := gini.NodeImpurity()
	expected := 0.5
	if math.Abs(impurity-expected) > epsilon {
		t.Errorf("Expected Gini impurity %.4f, got %.4f", expected, impurity)
	}

	// Node value should be [0.5, 0.5]
	value := gini.NodeValue()
	if len(value) != 2 {
		t.Errorf("Expected 2 classes, got %d", len(value))
	}
	if math.Abs(value[0]-0.5) > epsilon || math.Abs(value[1]-0.5) > epsilon {
		t.Errorf("Expected equal distribution, got %v", value)
	}
}

func TestMSECriterion(t *testing.T) {
	// Test MSE calculation
	y := [][]float64{
		{1}, {2}, {3}, {4}, {5},
	}
	mse := NewMSE()
	indices := []int{0, 1, 2, 3, 4}
	mse.Init(y, nil, indices, 0, 5)

	// Mean is 3, variance is 2
	// MSE should be E[(y - mean)^2] = 2
	impurity := mse.NodeImpurity()
	expected := 2.0
	if math.Abs(impurity-expected) > epsilon {
		t.Errorf("Expected MSE %.4f, got %.4f", expected, impurity)
	}

	// Node value should be mean = 3
	value := mse.NodeValue()
	if math.Abs(value[0]-3.0) > epsilon {
		t.Errorf("Expected mean 3.0, got %f", value[0])
	}
}

func TestTreePredict(t *testing.T) {
	// Build a simple tree manually
	tree := NewTree(1, 1, []int{0, 1})

	// Root: feature 0, threshold 5
	tree.AddNode(-1, false, 0, 5.0, 4, 0.5, []float64{})
	// Left leaf (x <= 5): class 0
	tree.AddLeaf(0, true, 2, 0.0, []float64{1.0, 0.0})
	// Right leaf (x > 5): class 1
	tree.AddLeaf(0, false, 2, 0.0, []float64{0.0, 1.0})

	// Test prediction
	result := tree.Predict([]float64{3})
	if len(result) != 2 || math.Abs(result[0]-1.0) > epsilon {
		t.Errorf("Expected [1.0, 0.0] for x=3, got %v", result)
	}

	result = tree.Predict([]float64{7})
	if len(result) != 2 || math.Abs(result[1]-1.0) > epsilon {
		t.Errorf("Expected [0.0, 1.0] for x=7, got %v", result)
	}
}

func TestEmptyData(t *testing.T) {
	clf := NewDecisionTreeClassifier()
	err := clf.Fit([][]float64{}, []float64{}, nil)
	if err == nil {
		t.Error("Expected error for empty data")
	}

	reg := NewDecisionTreeRegressor()
	err = reg.Fit([][]float64{}, []float64{}, nil)
	if err == nil {
		t.Error("Expected error for empty data")
	}
}

func TestNotFitted(t *testing.T) {
	clf := NewDecisionTreeClassifier()
	_, err := clf.Predict([][]float64{{1}})
	if err == nil {
		t.Error("Expected error when predicting before fitting")
	}

	reg := NewDecisionTreeRegressor()
	_, err = reg.Predict([][]float64{{1}})
	if err == nil {
		t.Error("Expected error when predicting before fitting")
	}
}
