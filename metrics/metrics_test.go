package metrics

import (
	"math"
	"testing"
)

func TestAccuracyScore(t *testing.T) {
	yTrue := []float64{0, 1, 2, 3}
	yPred := []float64{0, 2, 1, 3}
	expected := 0.5

	score, err := AccuracyScore(yTrue, yPred)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if score != expected {
		t.Errorf("Expected score %f, got %f", expected, score)
	}
}

func TestConfusionMatrix(t *testing.T) {
	yTrue := []float64{2, 0, 2, 2, 0, 1}
	yPred := []float64{0, 0, 2, 2, 0, 2}
	expected := [][]int{
		{2, 0, 0},
		{0, 0, 1},
		{1, 0, 2},
	}

	cm, labels, err := ConfusionMatrix(yTrue, yPred)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(labels) != 3 {
		t.Errorf("Expected 3 labels, got %d", len(labels))
	}

	for i := range expected {
		for j := range expected[i] {
			if cm[i][j] != expected[i][j] {
				t.Errorf("Matrix mismatch at [%d][%d]: expected %d, got %d", i, j, expected[i][j], cm[i][j])
			}
		}
	}
}

func TestClassificationMetrics(t *testing.T) {
	yTrue := []float64{0, 1, 1, 0, 1, 1}
	yPred := []float64{0, 0, 1, 0, 0, 1}

	// Precision = TP / (TP + FP) = 2 / (2 + 0) = 1.0
	// Recall = TP / (TP + FN) = 2 / (2 + 2) = 0.5
	// F1 = 2 * (1.0 * 0.5) / (1.0 + 0.5) = 1 / 1.5 = 0.666...

	p, _ := PrecisionScore(yTrue, yPred, "binary")
	r, _ := RecallScore(yTrue, yPred, "binary")
	f1, _ := F1Score(yTrue, yPred, "binary")

	if p != 1.0 {
		t.Errorf("Precision mismatch: expected 1.0, got %f", p)
	}
	if r != 0.5 {
		t.Errorf("Recall mismatch: expected 0.5, got %f", r)
	}
	if math.Abs(f1-0.666666666) > 1e-6 {
		t.Errorf("F1 mismatch: expected ~0.666, got %f", f1)
	}
}

func TestRegressionMetrics(t *testing.T) {
	yTrue := []float64{3, -0.5, 2, 7}
	yPred := []float64{2.5, 0.0, 2, 8}

	// MSE = (0.5^2 + 0.5^2 + 0^2 + 1^2) / 4 = (0.25 + 0.25 + 0 + 1) / 4 = 1.5 / 4 = 0.375
	// MAE = (0.5 + 0.5 + 0 + 1) / 4 = 2 / 4 = 0.5

	mse, _ := MeanSquaredError(yTrue, yPred)
	mae, _ := MeanAbsoluteError(yTrue, yPred)
	rmse, _ := RootMeanSquaredError(yTrue, yPred)

	if mse != 0.375 {
		t.Errorf("MSE mismatch: expected 0.375, got %f", mse)
	}
	if mae != 0.5 {
		t.Errorf("MAE mismatch: expected 0.5, got %f", mae)
	}
	if math.Abs(rmse-math.Sqrt(0.375)) > 1e-6 {
		t.Errorf("RMSE mismatch: expected %f, got %f", math.Sqrt(0.375), rmse)
	}

	yTrueR2 := []float64{1, 2, 3}
	yPredR2 := []float64{1, 2, 3}
	r2, _ := R2Score(yTrueR2, yPredR2)
	if r2 != 1.0 {
		t.Errorf("R2 mismatch: expected 1.0 for perfect fit, got %f", r2)
	}
}
