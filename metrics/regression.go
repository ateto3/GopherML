package metrics

import (
	"errors"
	"math"
)

// MeanSquaredError calculates the mean squared error.
func MeanSquaredError(yTrue, yPred []float64) (float64, error) {
	if len(yTrue) == 0 || len(yPred) == 0 {
		return 0, errors.New("empty input data")
	}
	if len(yTrue) != len(yPred) {
		return 0, errors.New("input arrays have different lengths")
	}

	sumSq := 0.0
	for i := range yTrue {
		diff := yTrue[i] - yPred[i]
		sumSq += diff * diff
	}

	return sumSq / float64(len(yTrue)), nil
}

// RootMeanSquaredError calculates the root mean squared error.
func RootMeanSquaredError(yTrue, yPred []float64) (float64, error) {
	mse, err := MeanSquaredError(yTrue, yPred)
	if err != nil {
		return 0, err
	}
	return math.Sqrt(mse), nil
}

// MeanAbsoluteError calculates the mean absolute error.
func MeanAbsoluteError(yTrue, yPred []float64) (float64, error) {
	if len(yTrue) == 0 || len(yPred) == 0 {
		return 0, errors.New("empty input data")
	}
	if len(yTrue) != len(yPred) {
		return 0, errors.New("input arrays have different lengths")
	}

	sumAbs := 0.0
	for i := range yTrue {
		sumAbs += math.Abs(yTrue[i] - yPred[i])
	}

	return sumAbs / float64(len(yTrue)), nil
}

// R2Score calculates the R-squared (coefficient of determination).
func R2Score(yTrue, yPred []float64) (float64, error) {
	if len(yTrue) == 0 || len(yPred) == 0 {
		return 0, errors.New("empty input data")
	}
	if len(yTrue) != len(yPred) {
		return 0, errors.New("input arrays have different lengths")
	}

	yMean := 0.0
	for _, v := range yTrue {
		yMean += v
	}
	yMean /= float64(len(yTrue))

	ssTot := 0.0
	ssRes := 0.0
	for i := range yTrue {
		diffTot := yTrue[i] - yMean
		diffRes := yTrue[i] - yPred[i]
		ssTot += diffTot * diffTot
		ssRes += diffRes * diffRes
	}

	if ssTot == 0 {
		return 1.0, nil
	}

	return 1.0 - (ssRes / ssTot), nil
}
