package utils

import (
	"errors"
	"fmt"
	"math"
)

// CheckXy ensures that X and y have compatible dimensions and contain finite values.
func CheckXy(X [][]float64, y []float64) error {
	if len(X) == 0 {
		return errors.New("X is empty")
	}
	if len(y) == 0 {
		return errors.New("y is empty")
	}
	if len(X) != len(y) {
		return fmt.Errorf("found input variables with inconsistent numbers of samples: [%d, %d]", len(X), len(y))
	}

	nFeatures := len(X[0])
	for i, row := range X {
		if len(row) != nFeatures {
			return fmt.Errorf("inconsistent number of features at row %d: expected %d, got %d", i, nFeatures, len(row))
		}
	}

	if err := CheckFinite(X); err != nil {
		return err
	}
	if err := CheckFiniteVec(y); err != nil {
		return fmt.Errorf("y contains invalid value: %w", err)
	}

	return nil
}

// CheckArray ensures that X is a valid 2D array with finite values.
func CheckArray(X [][]float64) error {
	if len(X) == 0 {
		return errors.New("array is empty")
	}
	nFeatures := len(X[0])
	for i, row := range X {
		if len(row) != nFeatures {
			return fmt.Errorf("inconsistent number of features at row %d: expected %d, got %d", i, nFeatures, len(row))
		}
	}
	return CheckFinite(X)
}

// CheckIsFitted checks if the model has been fitted.
func CheckIsFitted(fitted bool, modelName string) error {
	if !fitted {
		return fmt.Errorf("this %s instance is not fitted yet. Call 'Fit' with appropriate arguments before using this estimator", modelName)
	}
	return nil
}

// CheckFinite checks that all values in a 2D float64 slice are finite (not NaN or Inf).
// This is critical for crypto market data which may contain zero prices or missing candles.
func CheckFinite(X [][]float64) error {
	for i, row := range X {
		for j, v := range row {
			if math.IsNaN(v) {
				return fmt.Errorf("X[%d][%d] is NaN — ensure feature data is clean before fitting", i, j)
			}
			if math.IsInf(v, 0) {
				return fmt.Errorf("X[%d][%d] is Inf — ensure feature data is clean before fitting", i, j)
			}
		}
	}
	return nil
}

// CheckFiniteVec checks that all values in a 1D float64 slice are finite (not NaN or Inf).
func CheckFiniteVec(y []float64) error {
	for i, v := range y {
		if math.IsNaN(v) {
			return fmt.Errorf("value at index %d is NaN", i)
		}
		if math.IsInf(v, 0) {
			return fmt.Errorf("value at index %d is Inf", i)
		}
	}
	return nil
}

// CheckFiniteSingle checks that a single sample vector is finite.
// Used in hot-path single-sample prediction.
func CheckFiniteSingle(x []float64) error {
	for i, v := range x {
		if math.IsNaN(v) {
			return fmt.Errorf("feature[%d] is NaN", i)
		}
		if math.IsInf(v, 0) {
			return fmt.Errorf("feature[%d] is Inf", i)
		}
	}
	return nil
}
