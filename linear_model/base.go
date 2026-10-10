// Package linear_model provides linear models for regression and classification.
// Based on scikit-learn's linear_model module, optimized for Go with native concurrency.
package linear_model

import (
	"errors"
	"fmt"
	"math"
)

// ErrNotFitted is returned when a model is used before fitting.
var ErrNotFitted = errors.New("model not fitted")

// ErrEmptyData is returned when input data is empty.
var ErrEmptyData = errors.New("empty input data")

// ErrDimensionMismatch is returned when input dimensions are inconsistent.
var ErrDimensionMismatch = errors.New("dimension mismatch")

// LinearModel is the base struct for all linear models.
// It provides common attributes and methods for linear regression models.
type LinearModel struct {
	// Coefficients of the linear model
	Coef []float64
	// Intercept (bias) term
	Intercept float64
	// Whether to calculate the intercept
	FitIntercept bool
	// Number of features seen during fit
	NFeaturesIn int
	// Indicates if the model has been fitted
	IsFitted bool
}

// setIntercept computes the intercept from the centered data.
func (lm *LinearModel) setIntercept(xOffset []float64, yOffset float64) {
	// intercept = y_mean - x_mean^T * coef
	dotProduct := 0.0
	for j := 0; j < len(lm.Coef); j++ {
		dotProduct += xOffset[j] * lm.Coef[j]
	}
	lm.Intercept = yOffset - dotProduct
}

// Predict makes predictions using the fitted linear model.
func (lm *LinearModel) Predict(X [][]float64) ([]float64, error) {
	if !lm.IsFitted {
		return nil, ErrNotFitted
	}

	nSamples := len(X)
	if nSamples == 0 {
		return nil, errors.New("empty input data")
	}

	// Validate features
	for i, row := range X {
		if len(row) != lm.NFeaturesIn {
			return nil, fmt.Errorf("row %d has %d features, expected %d", i, len(row), lm.NFeaturesIn)
		}
	}

	predictions := make([]float64, nSamples)
	for i := 0; i < nSamples; i++ {
		pred := lm.Intercept
		for j := 0; j < lm.NFeaturesIn; j++ {
			pred += X[i][j] * lm.Coef[j]
		}
		predictions[i] = pred
	}

	return predictions, nil
}

// LinearRegressionParams contains parameters for LinearRegression.
type LinearRegressionParams struct {
	FitIntercept bool
	CopyX        bool
}

// NewLinearRegression creates a new LinearRegression model with default parameters.
func NewLinearRegression() *LinearRegression {
	return &LinearRegression{
		LinearModel: LinearModel{
			FitIntercept: true,
		},
		CopyX: true,
	}
}

// NewLinearRegressionWithParams creates a new LinearRegression with custom parameters.
func NewLinearRegressionWithParams(params LinearRegressionParams) *LinearRegression {
	return &LinearRegression{
		LinearModel: LinearModel{
			FitIntercept: params.FitIntercept,
		},
		CopyX: params.CopyX,
	}
}

// LinearRegression represents ordinary least squares linear regression.
type LinearRegression struct {
	LinearModel
	CopyX bool
	// Rank of matrix X (only for dense data)
	Rank int
	// Singular values of X
	Singular []float64
}

// Fit trains the linear regression model.
// X: matrix of shape (nSamples, nFeatures)
// y: target vector of shape (nSamples,)
func (lr *LinearRegression) Fit(X [][]float64, y []float64) error {
	if len(X) == 0 || len(y) == 0 {
		return errors.New("empty training data")
	}
	if len(X) != len(y) {
		return errors.New("X and y have inconsistent dimensions")
	}

	nSamples := len(X)
	nFeatures := len(X[0])
	lr.NFeaturesIn = nFeatures

	// Validate all rows have same number of features
	for i, row := range X {
		if len(row) != nFeatures {
			return fmt.Errorf("row %d has %d features, expected %d", i, len(row), nFeatures)
		}
	}

	// Create working copy if needed
	var XWork [][]float64
	var yWork []float64

	if lr.CopyX {
		XWork = make([][]float64, nSamples)
		for i := range X {
			XWork[i] = make([]float64, nFeatures)
			copy(XWork[i], X[i])
		}
		yWork = make([]float64, nSamples)
		copy(yWork, y)
	} else {
		XWork = X
		yWork = y
	}

	// Compute means for centering if fit_intercept is true
	var xOffset []float64
	var yOffset float64

	if lr.FitIntercept {
		xOffset = make([]float64, nFeatures)
		for j := 0; j < nFeatures; j++ {
			sum := 0.0
			for i := 0; i < nSamples; i++ {
				sum += XWork[i][j]
			}
			xOffset[j] = sum / float64(nSamples)
		}

		// Center X
		for i := 0; i < nSamples; i++ {
			for j := 0; j < nFeatures; j++ {
				XWork[i][j] -= xOffset[j]
			}
		}

		// Center y
		ySum := 0.0
		for i := 0; i < nSamples; i++ {
			ySum += yWork[i]
		}
		yOffset = ySum / float64(nSamples)
		for i := 0; i < nSamples; i++ {
			yWork[i] -= yOffset
		}
	}

	// Solve using normal equations: (X^T X) w = X^T y
	// Using Gaussian elimination for now (can be optimized with QR/SVD)
	coef, err := solveNormalEquations(XWork, yWork)
	if err != nil {
		return fmt.Errorf("failed to solve linear system: %w", err)
	}

	lr.Coef = coef

	// Set intercept
	if lr.FitIntercept {
		lr.LinearModel.setIntercept(xOffset, yOffset)
	} else {
		lr.Intercept = 0.0
	}

	lr.IsFitted = true
	return nil
}

// Predict makes predictions using the fitted linear regression model.
func (lr *LinearRegression) Predict(X [][]float64) ([]float64, error) {
	return lr.LinearModel.Predict(X)
}

// Score returns the coefficient of determination R^2 of the prediction.
func (lm *LinearModel) Score(X [][]float64, y []float64) (float64, error) {
	predictions, err := lm.Predict(X)
	if err != nil {
		return 0, err
	}

	return r2Score(y, predictions)
}

// solveNormalEquations solves (X^T X) w = X^T y using Gaussian elimination.
func solveNormalEquations(X [][]float64, y []float64) ([]float64, error) {
	nSamples := len(X)
	nFeatures := len(X[0])

	// Compute X^T X (Gram matrix)
	gram := make([][]float64, nFeatures)
	for i := range gram {
		gram[i] = make([]float64, nFeatures)
	}

	for i := 0; i < nFeatures; i++ {
		for j := i; j < nFeatures; j++ {
			sum := 0.0
			for k := 0; k < nSamples; k++ {
				sum += X[k][i] * X[k][j]
			}
			gram[i][j] = sum
			if i != j {
				gram[j][i] = sum
			}
		}
	}

	// Compute X^T y
	xty := make([]float64, nFeatures)
	for j := 0; j < nFeatures; j++ {
		sum := 0.0
		for i := 0; i < nSamples; i++ {
			sum += X[i][j] * y[i]
		}
		xty[j] = sum
	}

	// Solve using Gaussian elimination with partial pivoting
	return gaussianElimination(gram, xty)
}

// gaussianElimination solves Ax = b using Gaussian elimination with partial pivoting.
func gaussianElimination(A [][]float64, b []float64) ([]float64, error) {
	n := len(b)

	// Create augmented matrix
	aug := make([][]float64, n)
	for i := range aug {
		aug[i] = make([]float64, n+1)
		copy(aug[i], A[i])
		aug[i][n] = b[i]
	}

	// Forward elimination with partial pivoting
	for k := 0; k < n; k++ {
		// Find pivot
		maxRow := k
		maxVal := abs(aug[k][k])
		for i := k + 1; i < n; i++ {
			if abs(aug[i][k]) > maxVal {
				maxVal = abs(aug[i][k])
				maxRow = i
			}
		}

		if maxVal < 1e-12 {
			return nil, errors.New("matrix is singular or near-singular")
		}

		// Swap rows
		if maxRow != k {
			aug[k], aug[maxRow] = aug[maxRow], aug[k]
		}

		// Eliminate column
		for i := k + 1; i < n; i++ {
			factor := aug[i][k] / aug[k][k]
			for j := k; j <= n; j++ {
				aug[i][j] -= factor * aug[k][j]
			}
		}
	}

	// Back substitution
	x := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		x[i] = aug[i][n]
		for j := i + 1; j < n; j++ {
			x[i] -= aug[i][j] * x[j]
		}
		x[i] /= aug[i][i]
	}

	return x, nil
}

// r2Score computes the coefficient of determination.
func r2Score(yTrue, yPred []float64) (float64, error) {
	if len(yTrue) != len(yPred) {
		return 0, errors.New("yTrue and yPred have different lengths")
	}

	n := len(yTrue)
	if n == 0 {
		return 0, errors.New("empty arrays")
	}

	// Compute mean of yTrue
	mean := 0.0
	for _, v := range yTrue {
		mean += v
	}
	mean /= float64(n)

	// Compute total sum of squares and residual sum of squares
	tss := 0.0
	rss := 0.0
	for i := 0; i < n; i++ {
		tss += (yTrue[i] - mean) * (yTrue[i] - mean)
		diff := yTrue[i] - yPred[i]
		rss += diff * diff
	}

	if tss == 0 {
		return 1.0, nil // Perfect score if variance is 0
	}

	return 1.0 - rss/tss, nil
}

// abs returns the absolute value of a float64.
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// Mean computes the mean of a slice.
func Mean(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

// DotProduct computes the dot product of two vectors.
func DotProduct(a, b []float64) (float64, error) {
	if len(a) != len(b) {
		return 0, ErrDimensionMismatch
	}
	sum := 0.0
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum, nil
}

// MeanSquaredError computes the MSE.
func MeanSquaredError(yTrue, yPred []float64) (float64, error) {
	if len(yTrue) != len(yPred) {
		return 0, ErrDimensionMismatch
	}
	if len(yTrue) == 0 {
		return 0, ErrEmptyData
	}
	sum := 0.0
	for i := range yTrue {
		diff := yTrue[i] - yPred[i]
		sum += diff * diff
	}
	return sum / float64(len(yTrue)), nil
}

// sqrt returns the square root of a float64.
func sqrt(x float64) float64 {
	return math.Sqrt(x)
}
