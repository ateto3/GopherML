package ensemble

import (
	"math"
	"testing"
)

func TestRandomForestClassifier(t *testing.T) {
	X := [][]float64{
		{0, 0}, {0, 1}, {1, 0}, {1, 1},
		{0, 0}, {0, 1}, {1, 0}, {1, 1},
	}
	y := []float64{0, 1, 1, 0, 0, 1, 1, 0}

	rf := NewRandomForestClassifier()
	rf.NEstimators = 10
	rf.MaxDepth = 3
	rf.RandomState = 42

	if err := rf.Fit(X, y); err != nil {
		t.Fatalf("Fit failed: %v", err)
	}

	predictions, err := rf.Predict(X)
	if err != nil {
		t.Fatalf("Predict failed: %v", err)
	}
	if len(predictions) != len(y) {
		t.Errorf("Expected %d predictions, got %d", len(y), len(predictions))
	}

	score, err := rf.Score(X, y)
	if err != nil {
		t.Fatalf("Score failed: %v", err)
	}
	if score < 0.9 {
		t.Errorf("Expected score >= 0.9, got %.2f", score)
	}
}

func TestRandomForestClassifierPredictProba(t *testing.T) {
	X := [][]float64{
		{1, 0}, {1, 0}, {1, 0},
		{0, 1}, {0, 1}, {0, 1},
	}
	y := []float64{0, 0, 0, 1, 1, 1}

	rf := NewRandomForestClassifier()
	rf.NEstimators = 20
	rf.MaxDepth = 3
	rf.RandomState = 7

	if err := rf.Fit(X, y); err != nil {
		t.Fatalf("Fit failed: %v", err)
	}

	proba, err := rf.PredictProba(X)
	if err != nil {
		t.Fatalf("PredictProba failed: %v", err)
	}
	if len(proba) != len(X) {
		t.Fatalf("Expected %d proba rows, got %d", len(X), len(proba))
	}

	// Each row must sum to 1.0 (within float tolerance)
	for i, row := range proba {
		sum := 0.0
		for _, p := range row {
			if p < 0 || p > 1 {
				t.Errorf("proba[%d] has out-of-range value %f", i, p)
			}
			sum += p
		}
		if math.Abs(sum-1.0) > 1e-9 {
			t.Errorf("proba[%d] sums to %f, expected 1.0", i, sum)
		}
	}
}

func TestRandomForestClassifierPredictSingle(t *testing.T) {
	X := [][]float64{
		{1, 0}, {1, 0}, {1, 0},
		{0, 1}, {0, 1}, {0, 1},
	}
	y := []float64{0, 0, 0, 1, 1, 1}

	rf := NewRandomForestClassifier()
	rf.NEstimators = 20
	rf.MaxDepth = 3
	rf.RandomState = 99

	if err := rf.Fit(X, y); err != nil {
		t.Fatalf("Fit failed: %v", err)
	}

	// PredictSingle should match batch Predict for each sample
	batchPreds, _ := rf.Predict(X)
	for i, row := range X {
		single, err := rf.PredictSingle(row)
		if err != nil {
			t.Fatalf("PredictSingle[%d] failed: %v", i, err)
		}
		if single != batchPreds[i] {
			t.Errorf("PredictSingle[%d] = %d, batch = %d", i, single, batchPreds[i])
		}
	}
}

func TestRandomForestClassifierPredictProbaSingle(t *testing.T) {
	X := [][]float64{{1, 0}, {1, 0}, {0, 1}, {0, 1}}
	y := []float64{0, 0, 1, 1}

	rf := NewRandomForestClassifier()
	rf.NEstimators = 20
	rf.MaxDepth = 3
	rf.RandomState = 5

	if err := rf.Fit(X, y); err != nil {
		t.Fatalf("Fit failed: %v", err)
	}

	proba, err := rf.PredictProbaSingle(X[0])
	if err != nil {
		t.Fatalf("PredictProbaSingle failed: %v", err)
	}

	sum := 0.0
	for _, p := range proba {
		sum += p
	}
	if math.Abs(sum-1.0) > 1e-9 {
		t.Errorf("PredictProbaSingle sum = %f, expected 1.0", sum)
	}
}

func TestRandomForestClassifierRejectsNaN(t *testing.T) {
	X := [][]float64{{1, 0}, {math.NaN(), 1}}
	y := []float64{0, 1}

	rf := NewRandomForestClassifier()
	rf.NEstimators = 5
	if err := rf.Fit(X, y); err == nil {
		t.Error("Expected Fit to reject NaN in X, got nil error")
	}
}

func TestRandomForestRegressorPredictSingle(t *testing.T) {
	X := [][]float64{
		{1, 1}, {1, 2}, {2, 1}, {2, 2},
		{3, 3}, {3, 4}, {4, 3}, {4, 4},
	}
	y := []float64{2, 3, 3, 4, 6, 7, 7, 8}

	rf := NewRandomForestRegressor()
	rf.NEstimators = 20
	rf.MaxDepth = 4
	rf.RandomState = 42

	if err := rf.Fit(X, y); err != nil {
		t.Fatalf("Fit failed: %v", err)
	}

	batchPreds, _ := rf.Predict(X)
	for i, row := range X {
		single, err := rf.PredictSingle(row)
		if err != nil {
			t.Fatalf("PredictSingle[%d] failed: %v", i, err)
		}
		if math.Abs(single-batchPreds[i]) > 1e-9 {
			t.Errorf("PredictSingle[%d] = %f, batch = %f", i, single, batchPreds[i])
		}
	}
}

func TestRandomForestRegressor(t *testing.T) {
	X := [][]float64{
		{1, 1}, {1, 2}, {2, 1}, {2, 2},
		{3, 3}, {3, 4}, {4, 3}, {4, 4},
	}
	y := []float64{2, 3, 3, 4, 6, 7, 7, 8}

	rf := NewRandomForestRegressor()
	rf.NEstimators = 10
	rf.MaxDepth = 4
	rf.RandomState = 42

	if err := rf.Fit(X, y); err != nil {
		t.Fatalf("Fit failed: %v", err)
	}

	predictions, err := rf.Predict(X)
	if err != nil {
		t.Fatalf("Predict failed: %v", err)
	}

	for i, pred := range predictions {
		if math.Abs(pred-y[i]) > 1.5 {
			t.Errorf("Prediction %d too far: got %.2f, expected %.2f", i, pred, y[i])
		}
	}

	score, err := rf.Score(X, y)
	if err != nil {
		t.Fatalf("Score failed: %v", err)
	}
	if score < 0.8 {
		t.Errorf("Expected score >= 0.8, got %.2f", score)
	}
}
