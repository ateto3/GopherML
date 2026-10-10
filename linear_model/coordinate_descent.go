package linear_model

import (
	"math"

	"github.com/ateto3/gopherml/utils"
)

// LassoParams contains parameters for Lasso regression.
type LassoParams struct {
	Alpha        float64
	FitIntercept bool
	CopyX        bool
	MaxIter      int
	Tol          float64
	Positive     bool
}

// Lasso represents Lasso regression (L1 regularization).
type Lasso struct {
	LinearModel
	Alpha    float64
	CopyX    bool
	MaxIter  int
	Tol      float64
	Positive bool
	NIter    int
}

// NewLasso creates a new Lasso regression model.
func NewLasso() *Lasso {
	return &Lasso{
		LinearModel: LinearModel{FitIntercept: true},
		Alpha:       1.0,
		CopyX:       true,
		MaxIter:     1000,
		Tol:         1e-4,
	}
}

// NewLassoWithParams creates a new Lasso regression with custom parameters.
func NewLassoWithParams(params LassoParams) *Lasso {
	l := NewLasso()
	l.Alpha = params.Alpha
	l.FitIntercept = params.FitIntercept
	l.CopyX = params.CopyX
	l.MaxIter = params.MaxIter
	l.Tol = params.Tol
	l.Positive = params.Positive
	return l
}

// Fit trains the Lasso model.
func (l *Lasso) Fit(X [][]float64, y []float64) error {
	en := NewElasticNet()
	en.Alpha = l.Alpha
	en.L1Ratio = 1.0 // Lasso is ElasticNet with L1Ratio = 1
	en.FitIntercept = l.FitIntercept
	en.MaxIter = l.MaxIter
	en.Tol = l.Tol
	en.Positive = l.Positive

	err := en.Fit(X, y)
	if err != nil {
		return err
	}

	l.Coef = en.Coef
	l.Intercept = en.Intercept
	l.NFeaturesIn = en.NFeaturesIn
	l.IsFitted = true
	l.NIter = en.NIter
	return nil
}

// Predict makes predictions.
func (l *Lasso) Predict(X [][]float64) ([]float64, error) {
	return l.LinearModel.Predict(X)
}

// ElasticNetParams contains parameters for ElasticNet regression.
type ElasticNetParams struct {
	Alpha        float64
	L1Ratio      float64
	FitIntercept bool
	CopyX        bool
	MaxIter      int
	Tol          float64
	Positive     bool
}

// ElasticNet represents ElasticNet regression (L1 + L2).
type ElasticNet struct {
	LinearModel
	Alpha    float64
	L1Ratio  float64
	CopyX    bool
	MaxIter  int
	Tol      float64
	Positive bool
	NIter    int
}

// NewElasticNet creates a new ElasticNet model.
func NewElasticNet() *ElasticNet {
	return &ElasticNet{
		LinearModel: LinearModel{FitIntercept: true},
		Alpha:       1.0,
		L1Ratio:     0.5,
		CopyX:       true,
		MaxIter:     1000,
		Tol:         1e-4,
	}
}

// NewElasticNetWithParams creates a new ElasticNet regression with custom parameters.
func NewElasticNetWithParams(params ElasticNetParams) *ElasticNet {
	en := NewElasticNet()
	en.Alpha = params.Alpha
	en.L1Ratio = params.L1Ratio
	en.FitIntercept = params.FitIntercept
	en.CopyX = params.CopyX
	en.MaxIter = params.MaxIter
	en.Tol = params.Tol
	en.Positive = params.Positive
	return en
}

// Fit trains the ElasticNet model using coordinate descent.
func (en *ElasticNet) Fit(X [][]float64, y []float64) error {
	if err := utils.CheckXy(X, y); err != nil {
		return err
	}

	nSamples := len(X)
	nFeatures := len(X[0])
	en.NFeaturesIn = nFeatures

	// Work on copies so we can center without modifying caller's data.
	XWork := make([][]float64, nSamples)
	for i := range X {
		XWork[i] = make([]float64, nFeatures)
		copy(XWork[i], X[i])
	}
	yWork := make([]float64, nSamples)
	copy(yWork, y)

	// Center data when FitIntercept is true.
	xOffset := make([]float64, nFeatures)
	yOffset := 0.0
	if en.FitIntercept {
		for j := 0; j < nFeatures; j++ {
			for i := 0; i < nSamples; i++ {
				xOffset[j] += XWork[i][j]
			}
			xOffset[j] /= float64(nSamples)
		}
		for i := 0; i < nSamples; i++ {
			for j := 0; j < nFeatures; j++ {
				XWork[i][j] -= xOffset[j]
			}
		}
		for i := 0; i < nSamples; i++ {
			yOffset += yWork[i]
		}
		yOffset /= float64(nSamples)
		for i := 0; i < nSamples; i++ {
			yWork[i] -= yOffset
		}
	}

	en.Coef = make([]float64, nFeatures)

	// Precompute column norms squared.
	colNormSq := make([]float64, nFeatures)
	for j := 0; j < nFeatures; j++ {
		for i := 0; i < nSamples; i++ {
			colNormSq[j] += XWork[i][j] * XWork[i][j]
		}
		if colNormSq[j] < 1e-12 {
			colNormSq[j] = 1.0
		}
	}

	alphaL1 := en.Alpha * en.L1Ratio * float64(nSamples)
	alphaL2 := en.Alpha * (1 - en.L1Ratio) * float64(nSamples)

	// Maintain residuals to avoid O(n*p²) inner loop.
	residuals := make([]float64, nSamples)
	copy(residuals, yWork)

	en.NIter = en.MaxIter
	for iter := 0; iter < en.MaxIter; iter++ {
		maxChange := 0.0
		for j := 0; j < nFeatures; j++ {
			oldCoef := en.Coef[j]

			// Compute partial residual: r + X[:,j]*w[j]
			rho := 0.0
			for i := 0; i < nSamples; i++ {
				rho += XWork[i][j] * (residuals[i] + XWork[i][j]*oldCoef)
			}

			var newCoef float64
			if en.Positive && rho < alphaL1 {
				newCoef = 0
			} else {
				if rho > alphaL1 {
					newCoef = (rho - alphaL1) / (colNormSq[j] + alphaL2)
				} else if rho < -alphaL1 {
					newCoef = (rho + alphaL1) / (colNormSq[j] + alphaL2)
				} else {
					newCoef = 0
				}
			}
			en.Coef[j] = newCoef

			// Update residuals incrementally.
			delta := newCoef - oldCoef
			if delta != 0 {
				for i := 0; i < nSamples; i++ {
					residuals[i] -= XWork[i][j] * delta
				}
			}

			if math.Abs(delta) > maxChange {
				maxChange = math.Abs(delta)
			}
		}

		if maxChange < en.Tol {
			en.NIter = iter + 1
			break
		}
	}

	// Compute intercept from centered data.
	if en.FitIntercept {
		en.LinearModel.setIntercept(xOffset, yOffset)
	} else {
		en.Intercept = 0.0
	}

	en.IsFitted = true
	return nil
}

// Predict makes predictions.
func (en *ElasticNet) Predict(X [][]float64) ([]float64, error) {
	return en.LinearModel.Predict(X)
}
