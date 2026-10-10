package linear_model

import (
	"errors"
	"fmt"
)

// RidgeParams contains parameters for Ridge regression.
type RidgeParams struct {
	Alpha        float64 // Regularization strength
	FitIntercept bool
	CopyX        bool
	MaxIter      int
	Tol          float64 // Convergence tolerance
	Solver       string  // "auto", "svd", "cholesky", "lsqr"
}

// Ridge represents Ridge regression (L2 regularization).
type Ridge struct {
	LinearModel
	Alpha   float64
	CopyX   bool
	MaxIter int
	Tol     float64
	Solver  string
}

// NewRidge creates a new Ridge regression model with default parameters.
func NewRidge() *Ridge {
	return &Ridge{
		LinearModel: LinearModel{
			FitIntercept: true,
		},
		Alpha:   1.0,
		CopyX:   true,
		MaxIter: 1000,
		Tol:     1e-4,
		Solver:  "auto",
	}
}

// NewRidgeWithParams creates a new Ridge regression with custom parameters.
func NewRidgeWithParams(params RidgeParams) *Ridge {
	r := &Ridge{
		LinearModel: LinearModel{
			FitIntercept: params.FitIntercept,
		},
		Alpha:   params.Alpha,
		CopyX:   params.CopyX,
		MaxIter: params.MaxIter,
		Tol:     params.Tol,
		Solver:  params.Solver,
	}

	if r.Alpha < 0 {
		r.Alpha = 1.0
	}
	if r.MaxIter <= 0 {
		r.MaxIter = 1000
	}
	if r.Tol <= 0 {
		r.Tol = 1e-4
	}
	if r.Solver == "" {
		r.Solver = "auto"
	}

	return r
}

// Fit trains the Ridge regression model.
func (r *Ridge) Fit(X [][]float64, y []float64) error {
	if len(X) == 0 || len(y) == 0 {
		return errors.New("empty training data")
	}
	if len(X) != len(y) {
		return errors.New("X and y have inconsistent dimensions")
	}

	nSamples := len(X)
	nFeatures := len(X[0])
	r.NFeaturesIn = nFeatures

	// Validate all rows have same number of features
	for i, row := range X {
		if len(row) != nFeatures {
			return fmt.Errorf("row %d has %d features, expected %d", i, len(row), nFeatures)
		}
	}

	// Create working copy if needed
	var XWork [][]float64
	var yWork []float64

	if r.CopyX {
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

	if r.FitIntercept {
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

	// Solve Ridge regression: (X^T X + alpha*I) w = X^T y
	coef, err := r.solveRidge(XWork, yWork)
	if err != nil {
		return fmt.Errorf("failed to solve ridge regression: %w", err)
	}

	r.Coef = coef

	// Set intercept
	if r.FitIntercept {
		r.LinearModel.setIntercept(xOffset, yOffset)
	} else {
		r.Intercept = 0.0
	}

	r.IsFitted = true
	return nil
}

// solveRidge solves the Ridge regression problem using Cholesky decomposition.
func (r *Ridge) solveRidge(X [][]float64, y []float64) ([]float64, error) {
	nSamples := len(X)
	nFeatures := len(X[0])

	// Compute X^T X (Gram matrix) with regularization
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

	// Add regularization term: X^T X + alpha * I
	for i := 0; i < nFeatures; i++ {
		gram[i][i] += r.Alpha * float64(nSamples)
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

	// Solve using Cholesky decomposition
	return choleskySolve(gram, xty)
}

// choleskySolve solves Ax = b using Cholesky decomposition.
func choleskySolve(A [][]float64, b []float64) ([]float64, error) {
	n := len(b)
	L := make([][]float64, n)
	for i := range L {
		L[i] = make([]float64, n)
	}

	for i := 0; i < n; i++ {
		for j := 0; j <= i; j++ {
			sum := 0.0
			for k := 0; k < j; k++ {
				sum += L[i][k] * L[j][k]
			}

			if i == j {
				val := A[i][i] - sum
				if val <= 0 {
					return gaussianElimination(A, b)
				}
				L[i][j] = sqrt(val)
			} else {
				L[i][j] = (A[i][j] - sum) / L[j][j]
			}
		}
	}

	y := make([]float64, n)
	for i := 0; i < n; i++ {
		sum := 0.0
		for j := 0; j < i; j++ {
			sum += L[i][j] * y[j]
		}
		y[i] = (b[i] - sum) / L[i][i]
	}

	x := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		sum := 0.0
		for j := i + 1; j < n; j++ {
			sum += L[j][i] * x[j]
		}
		x[i] = (y[i] - sum) / L[i][i]
	}

	return x, nil
}

// Predict makes predictions using the fitted Ridge model.
func (r *Ridge) Predict(X [][]float64) ([]float64, error) {
	return r.LinearModel.Predict(X)
}

// Score returns the coefficient of determination R^2 of the prediction.
func (r *Ridge) Score(X [][]float64, y []float64) (float64, error) {
	return r.LinearModel.Score(X, y)
}
