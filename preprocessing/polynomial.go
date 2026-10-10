package preprocessing

import (
	"fmt"

	"github.com/ateto3/gopherml/utils"
)

// PolynomialFeatures generates polynomial and interaction features up to degree 2.
// Degrees higher than 2 are not supported and will return an error.
type PolynomialFeatures struct {
	Degree      int
	IncludeBias bool
}

// NewPolynomialFeatures creates a new PolynomialFeatures transformer.
// degree must be 1 or 2.
func NewPolynomialFeatures(degree int) *PolynomialFeatures {
	return &PolynomialFeatures{
		Degree:      degree,
		IncludeBias: true,
	}
}

// Fit validates the input array (maintained for API compatibility with sklearn).
func (p *PolynomialFeatures) Fit(X [][]float64) error {
	if p.Degree < 1 || p.Degree > 2 {
		return fmt.Errorf("PolynomialFeatures: degree %d is not supported; only degrees 1 and 2 are implemented", p.Degree)
	}
	return utils.CheckArray(X)
}

// Transform generates polynomial features.
// Returns an error if degree > 2 or if the input is invalid.
func (p *PolynomialFeatures) Transform(X [][]float64) ([][]float64, error) {
	if p.Degree < 1 || p.Degree > 2 {
		return nil, fmt.Errorf("PolynomialFeatures: degree %d is not supported; only degrees 1 and 2 are implemented", p.Degree)
	}
	if err := utils.CheckArray(X); err != nil {
		return nil, err
	}

	nSamples := len(X)
	nFeatures := len(X[0])

	// Pre-calculate output width: bias + linear + (degree-2 interactions+squares)
	nOutput := 0
	if p.IncludeBias {
		nOutput++
	}
	nOutput += nFeatures
	if p.Degree >= 2 {
		nOutput += (nFeatures * (nFeatures + 1)) / 2
	}

	result := make([][]float64, nSamples)
	for i := 0; i < nSamples; i++ {
		row := make([]float64, 0, nOutput)

		if p.IncludeBias {
			row = append(row, 1.0)
		}

		// Degree-1 terms
		for j := 0; j < nFeatures; j++ {
			row = append(row, X[i][j])
		}

		// Degree-2 terms: interactions and squares
		if p.Degree >= 2 {
			for j := 0; j < nFeatures; j++ {
				for k := j; k < nFeatures; k++ {
					row = append(row, X[i][j]*X[i][k])
				}
			}
		}

		result[i] = row
	}

	return result, nil
}

// FitTransform fits and transforms in one step.
func (p *PolynomialFeatures) FitTransform(X [][]float64) ([][]float64, error) {
	if err := p.Fit(X); err != nil {
		return nil, err
	}
	return p.Transform(X)
}
