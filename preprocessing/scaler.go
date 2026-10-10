package preprocessing

import (
	"errors"
	"fmt"
	"math"
)

// StandardScaler standardizes features by removing the mean and scaling to unit variance.
// Equivalent to sklearn.preprocessing.StandardScaler.
type StandardScaler struct {
	Mean     []float64
	Std      []float64
	IsFitted bool
}

// NewStandardScaler creates a new StandardScaler.
func NewStandardScaler() *StandardScaler {
	return &StandardScaler{}
}

// Fit computes the mean and std to be used for later scaling.
func (s *StandardScaler) Fit(X [][]float64) error {
	if len(X) == 0 {
		return errors.New("cannot fit StandardScaler on empty data")
	}
	nSamples := float64(len(X))
	nFeatures := len(X[0])
	s.Mean = make([]float64, nFeatures)
	s.Std = make([]float64, nFeatures)

	// Calculate Mean
	for _, row := range X {
		for j, val := range row {
			s.Mean[j] += val
		}
	}
	for j := range s.Mean {
		s.Mean[j] /= nSamples
	}

	// Calculate Std
	for _, row := range X {
		for j, val := range row {
			diff := val - s.Mean[j]
			s.Std[j] += diff * diff
		}
	}
	for j := range s.Std {
		s.Std[j] = math.Sqrt(s.Std[j] / nSamples)
		if s.Std[j] == 0 {
			s.Std[j] = 1.0 // Avoid division by zero (constant feature)
		}
	}

	s.IsFitted = true
	return nil
}

// Transform scales the data using the previously fitted parameters.
// Returns an error if Fit has not been called first or if any row has the wrong number of features.
func (s *StandardScaler) Transform(X [][]float64) ([][]float64, error) {
	if !s.IsFitted {
		return nil, errors.New("StandardScaler is not fitted — call Fit before Transform")
	}
	nFeatures := len(s.Mean)
	result := make([][]float64, len(X))
	for i, row := range X {
		if len(row) != nFeatures {
			return nil, fmt.Errorf("StandardScaler: row %d has %d features, expected %d", i, len(row), nFeatures)
		}
		result[i] = make([]float64, nFeatures)
		for j, val := range row {
			result[i][j] = (val - s.Mean[j]) / s.Std[j]
		}
	}
	return result, nil
}

// FitTransform fits the scaler and transforms X in one step.
// This is the preferred method to avoid accidentally fitting on test data.
func (s *StandardScaler) FitTransform(X [][]float64) ([][]float64, error) {
	if err := s.Fit(X); err != nil {
		return nil, err
	}
	return s.Transform(X)
}

// InverseTransform reverses the scaling transformation.
func (s *StandardScaler) InverseTransform(X [][]float64) ([][]float64, error) {
	if !s.IsFitted {
		return nil, errors.New("StandardScaler is not fitted")
	}
	result := make([][]float64, len(X))
	for i, row := range X {
		result[i] = make([]float64, len(row))
		for j, val := range row {
			result[i][j] = val*s.Std[j] + s.Mean[j]
		}
	}
	return result, nil
}

// MinMaxScaler scales features to the range [0, 1].
// Equivalent to sklearn.preprocessing.MinMaxScaler.
type MinMaxScaler struct {
	Min      []float64
	Max      []float64
	IsFitted bool
}

// NewMinMaxScaler creates a new MinMaxScaler.
func NewMinMaxScaler() *MinMaxScaler {
	return &MinMaxScaler{}
}

// Fit computes the min and max to be used for later scaling.
func (s *MinMaxScaler) Fit(X [][]float64) error {
	if len(X) == 0 {
		return errors.New("cannot fit MinMaxScaler on empty data")
	}
	nFeatures := len(X[0])
	s.Min = make([]float64, nFeatures)
	s.Max = make([]float64, nFeatures)

	for j := 0; j < nFeatures; j++ {
		s.Min[j] = math.Inf(1)
		s.Max[j] = math.Inf(-1)
	}

	for _, row := range X {
		for j, val := range row {
			if val < s.Min[j] {
				s.Min[j] = val
			}
			if val > s.Max[j] {
				s.Max[j] = val
			}
		}
	}

	s.IsFitted = true
	return nil
}

// Transform scales the data using the previously fitted parameters.
// Returns an error if Fit has not been called first or if any row has the wrong number of features.
func (s *MinMaxScaler) Transform(X [][]float64) ([][]float64, error) {
	if !s.IsFitted {
		return nil, errors.New("MinMaxScaler is not fitted — call Fit before Transform")
	}
	nFeatures := len(s.Min)
	result := make([][]float64, len(X))
	for i, row := range X {
		if len(row) != nFeatures {
			return nil, fmt.Errorf("MinMaxScaler: row %d has %d features, expected %d", i, len(row), nFeatures)
		}
		result[i] = make([]float64, nFeatures)
		for j, val := range row {
			denom := s.Max[j] - s.Min[j]
			if denom == 0 {
				result[i][j] = 0
			} else {
				result[i][j] = (val - s.Min[j]) / denom
			}
		}
	}
	return result, nil
}

// FitTransform fits the scaler and transforms X in one step.
// This is the preferred method to avoid accidentally fitting on test data.
func (s *MinMaxScaler) FitTransform(X [][]float64) ([][]float64, error) {
	if err := s.Fit(X); err != nil {
		return nil, err
	}
	return s.Transform(X)
}

// InverseTransform reverses the MinMax scaling.
func (s *MinMaxScaler) InverseTransform(X [][]float64) ([][]float64, error) {
	if !s.IsFitted {
		return nil, errors.New("MinMaxScaler is not fitted")
	}
	result := make([][]float64, len(X))
	for i, row := range X {
		result[i] = make([]float64, len(row))
		for j, val := range row {
			result[i][j] = val*(s.Max[j]-s.Min[j]) + s.Min[j]
		}
	}
	return result, nil
}
