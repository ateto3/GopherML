package linear_model

import (
	"math"

	"github.com/ateto3/gopherml/utils"
)

// LogisticRegression implements a binary logistic regression classifier.
type LogisticRegression struct {
	Weights      []float64
	Bias         float64
	LearningRate float64
	MaxIter      int
	Tol          float64
	IsFitted     bool
}

// NewLogisticRegression creates a new LogisticRegression with default parameters.
func NewLogisticRegression() *LogisticRegression {
	return &LogisticRegression{
		LearningRate: 0.01,
		MaxIter:      1000,
		Tol:          1e-4,
	}
}

// Fit trains the logistic regression model using Gradient Descent.
func (lr *LogisticRegression) Fit(X [][]float64, y []float64) error {
	if err := utils.CheckXy(X, y); err != nil {
		return err
	}

	nSamples := len(X)
	nFeatures := len(X[0])

	// Initialize weights and bias
	lr.Weights = make([]float64, nFeatures)
	lr.Bias = 0.0

	for iter := 0; iter < lr.MaxIter; iter++ {
		totalLoss := 0.0
		dw := make([]float64, nFeatures)
		db := 0.0

		for i := 0; i < nSamples; i++ {
			// Linear combination
			z := lr.Bias
			for j := 0; j < nFeatures; j++ {
				z += X[i][j] * lr.Weights[j]
			}

			// Sigmoid
			yPred := 1.0 / (1.0 + math.Exp(-z))

			// Gradient calculation
			diff := yPred - y[i]
			for j := 0; j < nFeatures; j++ {
				dw[j] += diff * X[i][j]
			}
			db += diff

			// Binary Cross-Entropy Loss
			if y[i] == 1 {
				totalLoss += -math.Log(yPred + 1e-15)
			} else {
				totalLoss += -math.Log(1.0 - yPred + 1e-15)
			}
		}

		// Update parameters
		for j := 0; j < nFeatures; j++ {
			lr.Weights[j] -= lr.LearningRate * (dw[j] / float64(nSamples))
		}
		lr.Bias -= lr.LearningRate * (db / float64(nSamples))

		// Check convergence (simplified)
		if iter > 0 && totalLoss/float64(nSamples) < lr.Tol {
			break
		}
	}

	lr.IsFitted = true
	return nil
}

// Predict predicts class labels for samples in X.
func (lr *LogisticRegression) Predict(X [][]float64) ([]int, error) {
	if err := utils.CheckIsFitted(lr.IsFitted, "LogisticRegression"); err != nil {
		return nil, err
	}

	nSamples := len(X)
	predictions := make([]int, nSamples)

	for i := 0; i < nSamples; i++ {
		z := lr.Bias
		for j := 0; j < len(lr.Weights); j++ {
			z += X[i][j] * lr.Weights[j]
		}
		yPred := 1.0 / (1.0 + math.Exp(-z))

		if yPred >= 0.5 {
			predictions[i] = 1
		} else {
			predictions[i] = 0
		}
	}

	return predictions, nil
}

// PredictProba predicts probability estimates.
func (lr *LogisticRegression) PredictProba(X [][]float64) ([][]float64, error) {
	if err := utils.CheckIsFitted(lr.IsFitted, "LogisticRegression"); err != nil {
		return nil, err
	}

	nSamples := len(X)
	probabilities := make([][]float64, nSamples)

	for i := 0; i < nSamples; i++ {
		z := lr.Bias
		for j := 0; j < len(lr.Weights); j++ {
			z += X[i][j] * lr.Weights[j]
		}
		p1 := 1.0 / (1.0 + math.Exp(-z))
		probabilities[i] = []float64{1.0 - p1, p1}
	}

	return probabilities, nil
}

// Score returns the mean accuracy on the given test data and labels.
func (lr *LogisticRegression) Score(X [][]float64, y []float64) (float64, error) {
	preds, err := lr.Predict(X)
	if err != nil {
		return 0, err
	}

	correct := 0
	for i := range y {
		if float64(preds[i]) == y[i] {
			correct++
		}
	}
	return float64(correct) / float64(len(y)), nil
}
