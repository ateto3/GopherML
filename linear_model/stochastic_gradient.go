package linear_model

import (
	"math"

	"github.com/ateto3/gopherml/utils"
)

// SGDClassifier implements stochastic gradient descent for classification.
type SGDClassifier struct {
	Weights      []float64
	Bias         float64
	LearningRate float64
	MaxIter      int
	Tol          float64
	IsFitted     bool
	RandomState  int64
}

// NewSGDClassifier creates a new SGDClassifier.
func NewSGDClassifier() *SGDClassifier {
	return &SGDClassifier{
		LearningRate: 0.01,
		MaxIter:      1000,
		Tol:          1e-3,
		RandomState:  -1,
	}
}

// Fit trains the model using SGD.
func (s *SGDClassifier) Fit(X [][]float64, y []float64) error {
	if err := utils.CheckXy(X, y); err != nil {
		return err
	}

	nSamples := len(X)
	nFeatures := len(X[0])
	rng := utils.NewRNG(s.RandomState)

	s.Weights = make([]float64, nFeatures)
	s.Bias = 0.0

	for iter := 0; iter < s.MaxIter; iter++ {
		indices := utils.Shuffle(nSamples, rng)

		for _, i := range indices {
			// Linear prediction
			z := s.Bias
			for j := 0; j < nFeatures; j++ {
				z += X[i][j] * s.Weights[j]
			}

			// Sigmoid
			yPred := 1.0 / (1.0 + math.Exp(-z))

			// Gradient update (Hinge or Log Loss)
			diff := yPred - y[i]

			for j := 0; j < nFeatures; j++ {
				s.Weights[j] -= s.LearningRate * diff * X[i][j]
			}
			s.Bias -= s.LearningRate * diff
		}
	}

	s.IsFitted = true
	return nil
}

// Predict predicts class labels.
func (s *SGDClassifier) Predict(X [][]float64) ([]int, error) {
	if err := utils.CheckIsFitted(s.IsFitted, "SGDClassifier"); err != nil {
		return nil, err
	}

	nSamples := len(X)
	preds := make([]int, nSamples)
	for i := 0; i < nSamples; i++ {
		z := s.Bias
		for j := 0; j < len(s.Weights); j++ {
			z += X[i][j] * s.Weights[j]
		}
		if z > 0 {
			preds[i] = 1
		} else {
			preds[i] = 0
		}
	}
	return preds, nil
}
