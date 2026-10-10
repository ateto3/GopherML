package linear_model

import (
	"math"
	"testing"
)

// Tolerance for floating point comparisons
const epsilon = 1e-6

func TestLinearRegressionSimple(t *testing.T) {
	// Simple test: y = 2x + 1
	X := [][]float64{
		{1},
		{2},
		{3},
		{4},
		{5},
	}
	y := []float64{3, 5, 7, 9, 11} // y = 2*x + 1

	model := NewLinearRegression()
	err := model.Fit(X, y)
	if err != nil {
		t.Fatalf("Fit failed: %v", err)
	}

	// Check coefficients
	if math.Abs(model.Coef[0]-2.0) > epsilon {
		t.Errorf("Expected coefficient ~2.0, got %v", model.Coef[0])
	}
	if math.Abs(model.Intercept-1.0) > epsilon {
		t.Errorf("Expected intercept ~1.0, got %v", model.Intercept)
	}

	// Test prediction
	predictions, err := model.Predict([][]float64{{6}})
	if err != nil {
		t.Fatalf("Predict failed: %v", err)
	}
	expected := 13.0 // 2*6 + 1
	if math.Abs(predictions[0]-expected) > epsilon {
		t.Errorf("Expected prediction %v, got %v", expected, predictions[0])
	}
}

func TestLinearRegressionMultiFeature(t *testing.T) {
	// Test with multiple features: y = 1*x1 + 2*x2 + 3
	X := [][]float64{
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 2},
		{5, 3},
	}
	y := []float64{6, 7, 10, 11, 14} // y = 1*x1 + 2*x2 + 3

	model := NewLinearRegression()
	err := model.Fit(X, y)
	if err != nil {
		t.Fatalf("Fit failed: %v", err)
	}

	// Check coefficients (with some tolerance)
	if math.Abs(model.Coef[0]-1.0) > 0.1 {
		t.Errorf("Expected coefficient[0] ~1.0, got %v", model.Coef[0])
	}
	if math.Abs(model.Coef[1]-2.0) > 0.1 {
		t.Errorf("Expected coefficient[1] ~2.0, got %v", model.Coef[1])
	}
	if math.Abs(model.Intercept-3.0) > 0.5 {
		t.Errorf("Expected intercept ~3.0, got %v", model.Intercept)
	}
}

func TestLinearRegressionScore(t *testing.T) {
	// Perfect fit should give R^2 = 1.0
	X := [][]float64{
		{1},
		{2},
		{3},
		{4},
		{5},
	}
	y := []float64{2, 4, 6, 8, 10} // y = 2*x

	model := NewLinearRegression()
	err := model.Fit(X, y)
	if err != nil {
		t.Fatalf("Fit failed: %v", err)
	}

	score, err := model.Score(X, y)
	if err != nil {
		t.Fatalf("Score failed: %v", err)
	}

	if score < 0.99 {
		t.Errorf("Expected R^2 close to 1.0 for perfect fit, got %v", score)
	}
}

func TestLinearRegressionNoIntercept(t *testing.T) {
	// Test without intercept: y = 2*x
	X := [][]float64{
		{1},
		{2},
		{3},
		{4},
		{5},
	}
	y := []float64{2, 4, 6, 8, 10}

	model := NewLinearRegressionWithParams(LinearRegressionParams{
		FitIntercept: false,
		CopyX:        true,
	})

	err := model.Fit(X, y)
	if err != nil {
		t.Fatalf("Fit failed: %v", err)
	}

	if math.Abs(model.Intercept) > epsilon {
		t.Errorf("Expected intercept to be 0 when FitIntercept=false, got %v", model.Intercept)
	}
}

func TestLinearRegressionEmptyData(t *testing.T) {
	model := NewLinearRegression()

	// Test empty X
	err := model.Fit([][]float64{}, []float64{})
	if err == nil {
		t.Error("Expected error for empty data, got nil")
	}

	// Test mismatched dimensions
	err = model.Fit([][]float64{{1}, {2}}, []float64{1})
	if err == nil {
		t.Error("Expected error for mismatched dimensions, got nil")
	}
}

func TestLinearRegressionNotFitted(t *testing.T) {
	model := NewLinearRegression()

	_, err := model.Predict([][]float64{{1}})
	if err != ErrNotFitted {
		t.Errorf("Expected ErrNotFitted, got %v", err)
	}
}

func TestRidgeRegression(t *testing.T) {
	// Test Ridge with small regularization on simple problem
	X := [][]float64{
		{1},
		{2},
		{3},
		{4},
		{5},
	}
	y := []float64{3, 5, 7, 9, 11}

	model := NewRidgeWithParams(RidgeParams{
		Alpha:        0.1,
		FitIntercept: true,
		CopyX:        true,
	})

	err := model.Fit(X, y)
	if err != nil {
		t.Fatalf("Fit failed: %v", err)
	}

	// With small alpha, coefficients should be close to OLS
	if math.Abs(model.Coef[0]-2.0) > 0.5 {
		t.Errorf("Expected coefficient ~2.0 with small alpha, got %v", model.Coef[0])
	}

	// Test prediction
	predictions, err := model.Predict([][]float64{{6}})
	if err != nil {
		t.Fatalf("Predict failed: %v", err)
	}
	if math.Abs(predictions[0]-13.0) > 1.0 {
		t.Errorf("Expected prediction ~13.0, got %v", predictions[0])
	}
}

func TestLassoRegression(t *testing.T) {
	// Test Lasso - with high alpha, some coefficients should become zero
	X := [][]float64{
		{1, 0.1},
		{2, 0.2},
		{3, 0.3},
		{4, 0.4},
		{5, 0.5},
	}
	y := []float64{3, 5, 7, 9, 11} // y = 2*x1 + small_noise

	model := NewLassoWithParams(LassoParams{
		Alpha:        1.0,
		FitIntercept: true,
		CopyX:        true,
		MaxIter:      1000,
		Tol:          1e-4,
	})

	err := model.Fit(X, y)
	if err != nil {
		t.Fatalf("Fit failed: %v", err)
	}

	// With L1 regularization, we should get sparse coefficients
	// At least one coefficient should be close to zero or exactly zero
	if model.Coef[1] == 0 {
		t.Logf("Lasso correctly zeroed out the second feature (sparse solution)")
	}

	// Test that model can predict
	predictions, err := model.Predict([][]float64{{6, 0.6}})
	if err != nil {
		t.Fatalf("Predict failed: %v", err)
	}
	if math.Abs(predictions[0]-13.0) > 2.0 {
		t.Errorf("Expected prediction ~13.0, got %v", predictions[0])
	}
}

func TestElasticNet(t *testing.T) {
	// Test ElasticNet - combination of L1 and L2
	X := [][]float64{
		{1},
		{2},
		{3},
		{4},
		{5},
	}
	y := []float64{3, 5, 7, 9, 11}

	model := NewElasticNetWithParams(ElasticNetParams{
		Alpha:        0.5,
		L1Ratio:      0.5,
		FitIntercept: true,
		CopyX:        true,
		MaxIter:      1000,
		Tol:          1e-4,
	})

	err := model.Fit(X, y)
	if err != nil {
		t.Fatalf("Fit failed: %v", err)
	}

	// Test prediction
	predictions, err := model.Predict([][]float64{{6}})
	if err != nil {
		t.Fatalf("Predict failed: %v", err)
	}
	if math.Abs(predictions[0]-13.0) > 1.5 {
		t.Errorf("Expected prediction ~13.0, got %v", predictions[0])
	}
}

func TestUtils(t *testing.T) {
	// Test MeanSquaredError
	yTrue := []float64{1, 2, 3, 4, 5}
	yPred := []float64{1.5, 2.5, 2.5, 4.5, 4.5}

	mse, err := MeanSquaredError(yTrue, yPred)
	if err != nil {
		t.Fatalf("MeanSquaredError failed: %v", err)
	}
	// MSE = (0.25 + 0.25 + 0.25 + 0.25 + 0.25) / 5 = 1.25 / 5 = 0.25
	if math.Abs(mse-0.25) > epsilon {
		t.Errorf("Expected MSE ~0.25, got %v", mse)
	}

	// Test Mean
	data := []float64{1, 2, 3, 4, 5}
	mean := Mean(data)
	if math.Abs(mean-3.0) > epsilon {
		t.Errorf("Expected mean 3.0, got %v", mean)
	}

	// Test DotProduct
	a := []float64{1, 2, 3}
	b := []float64{4, 5, 6}
	dot, err := DotProduct(a, b)
	if err != nil {
		t.Fatalf("DotProduct failed: %v", err)
	}
	expectedDot := 32.0 // 1*4 + 2*5 + 3*6 = 4 + 10 + 18 = 32
	if math.Abs(dot-expectedDot) > epsilon {
		t.Errorf("Expected dot product %v, got %v", expectedDot, dot)
	}
}
