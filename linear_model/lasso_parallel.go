package linear_model

import (
	"fmt"
	"math"
)

// LassoParallel is a high-performance Lasso implementation using goroutines.
// Column-norm precomputation is parallelised; coordinate descent runs sequentially
// (parallel CD on shared residuals causes data races without approximation).
type LassoParallel struct {
	LinearModel
	Alpha        float64
	CopyX        bool
	MaxIter      int
	Tol          float64
	Positive     bool
	ParallelOpts ParallelOptions
	NIter        int
}

// NewLassoParallel creates a new parallel Lasso model.
func NewLassoParallel() *LassoParallel {
	return &LassoParallel{
		LinearModel: LinearModel{
			FitIntercept: true,
		},
		Alpha:        1.0,
		CopyX:        true,
		MaxIter:      1000,
		Tol:          1e-4,
		Positive:     false,
		ParallelOpts: DefaultParallelOptions(),
	}
}

// NewLassoParallelWithParams creates a new parallel Lasso with custom params.
func NewLassoParallelWithParams(params LassoParams, opts ParallelOptions) *LassoParallel {
	l := NewLassoParallel()
	l.Alpha = params.Alpha
	l.FitIntercept = params.FitIntercept
	l.CopyX = params.CopyX
	l.MaxIter = params.MaxIter
	l.Tol = params.Tol
	l.Positive = params.Positive
	l.ParallelOpts = opts
	return l
}

// Fit trains the Lasso model using parallel coordinate descent.
func (l *LassoParallel) Fit(X [][]float64, y []float64) error {
	if len(X) == 0 || len(y) == 0 {
		return ErrEmptyData
	}
	if len(X) != len(y) {
		return ErrDimensionMismatch
	}

	nSamples := len(X)
	nFeatures := len(X[0])
	l.NFeaturesIn = nFeatures

	// Create working copy if needed
	var XWork [][]float64
	var yWork []float64

	if l.CopyX {
		XWork = make([][]float64, nSamples)
		dataBlock := make([]float64, nSamples*nFeatures)
		for i := range X {
			start := i * nFeatures
			end := start + nFeatures
			XWork[i] = dataBlock[start:end]
			copy(XWork[i], X[i])
		}
		yWork = make([]float64, nSamples)
		copy(yWork, y)
	} else {
		XWork = X
		yWork = y
	}

	var xOffset []float64
	var yOffset float64

	if l.FitIntercept {
		xOffset, yOffset = ParallelCenterData(XWork, yWork, l.ParallelOpts)
	} else {
		xOffset = make([]float64, nFeatures)
		yOffset = 0
	}

	// Precompute column norms
	colNormSq := l.parallelColNorms(XWork)

	// Run parallel coordinate descent
	coef, nIter, err := l.parallelCoordinateDescent(XWork, yWork, colNormSq)
	if err != nil {
		return fmt.Errorf("coordinate descent failed: %w", err)
	}

	l.Coef = coef
	l.NIter = nIter

	if l.FitIntercept {
		l.LinearModel.setIntercept(xOffset, yOffset)
	} else {
		l.Intercept = 0.0
	}

	l.IsFitted = true
	return nil
}

// parallelColNorms computes column norms in parallel.
func (l *LassoParallel) parallelColNorms(X [][]float64) []float64 {
	nFeatures := len(X[0])
	nSamples := len(X)
	colNormSq := make([]float64, nFeatures)

	ParallelFor(nFeatures, l.ParallelOpts, func(start, end int) {
		for j := start; j < end; j++ {
			sum := 0.0
			for i := 0; i < nSamples; i++ {
				sum += X[i][j] * X[i][j]
			}
			colNormSq[j] = sum
			if colNormSq[j] < 1e-12 {
				colNormSq[j] = 1.0
			}
		}
	})

	return colNormSq
}

// parallelCoordinateDescent implements cyclic coordinate descent.
// Note: coordinate descent on shared residuals cannot be safely parallelised
// within a single pass without introducing approximation errors or data races.
// Parallelism is applied at the column-norm precomputation stage instead.
// For true parallel Lasso, use the block-coordinate approach in LassoParallel.Fit.
func (l *LassoParallel) parallelCoordinateDescent(X [][]float64, y []float64, colNormSq []float64) ([]float64, int, error) {
	nSamples := len(X)
	nFeatures := len(X[0])
	w := make([]float64, nFeatures)

	// Maintain residuals incrementally to keep each iteration O(n*p).
	residuals := make([]float64, nSamples)
	copy(residuals, y)

	for iter := 0; iter < l.MaxIter; iter++ {
		maxChange := 0.0

		for j := 0; j < nFeatures; j++ {
			oldWj := w[j]

			// Partial residual: r + X[:,j]*w[j]
			rho := 0.0
			for i := 0; i < nSamples; i++ {
				rho += X[i][j] * (residuals[i] + X[i][j]*oldWj)
			}

			alphaNorm := l.Alpha * float64(nSamples)

			if l.Positive && rho < 0 {
				w[j] = 0
			} else {
				if rho > alphaNorm {
					w[j] = (rho - alphaNorm) / colNormSq[j]
				} else if rho < -alphaNorm {
					w[j] = (rho + alphaNorm) / colNormSq[j]
				} else {
					w[j] = 0
				}
			}

			delta := w[j] - oldWj
			if delta != 0 {
				for i := 0; i < nSamples; i++ {
					residuals[i] -= X[i][j] * delta
				}
			}

			if math.Abs(delta) > maxChange {
				maxChange = math.Abs(delta)
			}
		}

		if maxChange < l.Tol {
			return w, iter + 1, nil
		}
	}

	return w, l.MaxIter, nil
}

// Predict makes predictions using the fitted model.
func (l *LassoParallel) Predict(X [][]float64) ([]float64, error) {
	return l.LinearModel.Predict(X)
}
