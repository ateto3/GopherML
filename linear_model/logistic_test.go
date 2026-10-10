package linear_model

import (
	"testing"
)

func TestLogisticRegression(t *testing.T) {
	// Simple linearly separable data
	X := [][]float64{
		{1, 1}, {2, 2}, {1, 2},
		{10, 10}, {11, 11}, {10, 12},
	}
	y := []float64{0, 0, 0, 1, 1, 1}

	lr := NewLogisticRegression()
	lr.LearningRate = 0.1
	lr.MaxIter = 500
	
	err := lr.Fit(X, y)
	if err != nil {
		t.Fatalf("Fit failed: %v", err)
	}

	predictions, err := lr.Predict(X)
	if err != nil {
		t.Fatalf("Predict failed: %v", err)
	}

	for i := range y {
		if float64(predictions[i]) != y[i] {
			t.Errorf("Mismatch at index %d: expected %f, got %d", i, y[i], predictions[i])
		}
	}
}
