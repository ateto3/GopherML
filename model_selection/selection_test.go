package model_selection

import (
	"testing"
)

// MockEstimator for testing
type MockEstimator struct {
	Fitted bool
}

func (m *MockEstimator) Fit(X [][]float64, y []float64) error {
	m.Fitted = true
	return nil
}

func (m *MockEstimator) Score(X [][]float64, y []float64) (float64, error) {
	return 1.0, nil
}

func TestKFold(t *testing.T) {
	kf := KFold{NSplits: 3, Shuffle: false}
	records := kf.Split(6)

	if len(records) != 3 {
		t.Errorf("Expected 3 splits, got %d", len(records))
	}

	// First fold: test indices [0, 1], train [2, 3, 4, 5]
	if len(records[0].TestIndices) != 2 || records[0].TestIndices[0] != 0 || records[0].TestIndices[1] != 1 {
		t.Errorf("Fold 0 test indices mismatch: %v", records[0].TestIndices)
	}
}

func TestCrossValScore(t *testing.T) {
	X := [][]float64{{1}, {2}, {3}, {4}, {5}, {6}}
	y := []float64{0, 0, 1, 1, 0, 0}
	
	model := &MockEstimator{}
	kf := KFold{NSplits: 3, Shuffle: false}
	
	scores, err := CrossValScore(model, X, y, kf)
	if err != nil {
		t.Fatalf("CrossValScore failed: %v", err)
	}
	
	if len(scores) != 3 {
		t.Errorf("Expected 3 scores, got %d", len(scores))
	}
	
	for _, s := range scores {
		if s != 1.0 {
			t.Errorf("Expected score 1.0, got %f", s)
		}
	}
}
