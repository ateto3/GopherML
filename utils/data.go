package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// LoadCSV loads a CSV file and separates features and target (assumes target is in the last column).
func LoadCSV(path string, hasHeader bool) ([][]float64, []float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	startRow := 0
	if hasHeader {
		startRow = 1
	}

	nSamples := len(records) - startRow
	if nSamples <= 0 {
		return nil, nil, os.ErrInvalid
	}

	nCols := len(records[0])
	X := make([][]float64, nSamples)
	y := make([]float64, nSamples)

	for i := 0; i < nSamples; i++ {
		row := records[i+startRow]
		X[i] = make([]float64, nCols-1)
		for j := 0; j < nCols-1; j++ {
			val, err := strconv.ParseFloat(row[j], 64)
			if err != nil {
				return nil, nil, err
			}
			X[i][j] = val
		}
		targetVal, err := strconv.ParseFloat(row[nCols-1], 64)
		if err != nil {
			return nil, nil, err
		}
		y[i] = targetVal
	}

	return X, y, nil
}

// TrainTestSplit splits the data into training and testing sets.
// testSize must be in the range (0, 1).
func TrainTestSplit(X [][]float64, y []float64, testSize float64, randomState int64) ([][]float64, [][]float64, []float64, []float64, error) {
	if testSize <= 0 || testSize >= 1 {
		return nil, nil, nil, nil, fmt.Errorf("TrainTestSplit: testSize must be in (0, 1), got %f", testSize)
	}
	if len(X) == 0 {
		return nil, nil, nil, nil, fmt.Errorf("TrainTestSplit: empty input")
	}
	if len(X) != len(y) {
		return nil, nil, nil, nil, fmt.Errorf("TrainTestSplit: X and y have different lengths (%d vs %d)", len(X), len(y))
	}

	rng := NewRNG(randomState)
	n := len(X)
	indices := Shuffle(n, rng)

	nTest := int(float64(n) * testSize)
	if nTest == 0 {
		nTest = 1
	}
	nTrain := n - nTest

	XTrain := make([][]float64, nTrain)
	XTest := make([][]float64, nTest)
	yTrain := make([]float64, nTrain)
	yTest := make([]float64, nTest)

	for i := 0; i < nTrain; i++ {
		XTrain[i] = X[indices[i]]
		yTrain[i] = y[indices[i]]
	}
	for i := 0; i < nTest; i++ {
		XTest[i] = X[indices[i+nTrain]]
		yTest[i] = y[indices[i+nTrain]]
	}

	return XTrain, XTest, yTrain, yTest, nil
}
