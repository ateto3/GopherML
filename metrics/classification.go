package metrics

import (
	"errors"
	"fmt"
	"sort"
)

// AccuracyScore calculates the accuracy of predictions.
func AccuracyScore(yTrue, yPred []float64) (float64, error) {
	if len(yTrue) == 0 || len(yPred) == 0 {
		return 0, errors.New("empty input data")
	}
	if len(yTrue) != len(yPred) {
		return 0, errors.New("input arrays have different lengths")
	}

	correct := 0
	for i := range yTrue {
		if yTrue[i] == yPred[i] {
			correct++
		}
	}

	return float64(correct) / float64(len(yTrue)), nil
}

// ConfusionMatrix calculates the confusion matrix.
func ConfusionMatrix(yTrue, yPred []float64) ([][]int, []float64, error) {
	if len(yTrue) == 0 || len(yPred) == 0 {
		return nil, nil, errors.New("empty input data")
	}
	if len(yTrue) != len(yPred) {
		return nil, nil, errors.New("input arrays have different lengths")
	}

	// Find unique labels
	labelMap := make(map[float64]int)
	for _, v := range yTrue {
		labelMap[v] = 0
	}
	for _, v := range yPred {
		labelMap[v] = 0
	}

	labels := make([]float64, 0, len(labelMap))
	for l := range labelMap {
		labels = append(labels, l)
	}
	sort.Float64s(labels)

	for i, l := range labels {
		labelMap[l] = i
	}

	n := len(labels)
	matrix := make([][]int, n)
	for i := range matrix {
		matrix[i] = make([]int, n)
	}

	for i := range yTrue {
		trueIdx := labelMap[yTrue[i]]
		predIdx := labelMap[yPred[i]]
		matrix[trueIdx][predIdx]++
	}

	return matrix, labels, nil
}

// PrecisionScore calculates precision. average can be "macro", "weighted" or "".
func PrecisionScore(yTrue, yPred []float64, average string) (float64, error) {
	cm, _, err := ConfusionMatrix(yTrue, yPred)
	if err != nil {
		return 0, err
	}

	n := len(cm)
	precisions := make([]float64, n)
	weights := make([]float64, n)
	totalSamples := float64(len(yTrue))

	for i := 0; i < n; i++ {
		tp := float64(cm[i][i])
		fp := 0.0
		rowSum := 0.0
		for j := 0; j < n; j++ {
			if i != j {
				fp += float64(cm[j][i])
			}
			rowSum += float64(cm[i][j])
		}

		if tp+fp == 0 {
			precisions[i] = 0
		} else {
			precisions[i] = tp / (tp + fp)
		}
		weights[i] = rowSum / totalSamples
	}

	return averageScores(precisions, weights, average)
}

// RecallScore calculates recall. average can be "macro", "weighted" or "".
func RecallScore(yTrue, yPred []float64, average string) (float64, error) {
	cm, _, err := ConfusionMatrix(yTrue, yPred)
	if err != nil {
		return 0, err
	}

	n := len(cm)
	recalls := make([]float64, n)
	weights := make([]float64, n)
	totalSamples := float64(len(yTrue))

	for i := 0; i < n; i++ {
		tp := float64(cm[i][i])
		fn := 0.0
		rowSum := 0.0
		for j := 0; j < n; j++ {
			if i != j {
				fn += float64(cm[i][j])
			}
			rowSum += float64(cm[i][j])
		}

		if tp+fn == 0 {
			recalls[i] = 0
		} else {
			recalls[i] = tp / (tp + fn)
		}
		weights[i] = rowSum / totalSamples
	}

	return averageScores(recalls, weights, average)
}

// F1Score calculates F1 score. average can be "macro", "weighted" or "".
func F1Score(yTrue, yPred []float64, average string) (float64, error) {
	cm, _, err := ConfusionMatrix(yTrue, yPred)
	if err != nil {
		return 0, err
	}

	n := len(cm)
	f1s := make([]float64, n)
	weights := make([]float64, n)
	totalSamples := float64(len(yTrue))

	for i := 0; i < n; i++ {
		tp := float64(cm[i][i])
		fp := 0.0
		fn := 0.0
		rowSum := 0.0
		for j := 0; j < n; j++ {
			if i != j {
				fp += float64(cm[j][i])
				fn += float64(cm[i][j])
			}
			rowSum += float64(cm[i][j])
		}

		if tp+fp+fn == 0 {
			f1s[i] = 0
		} else {
			f1s[i] = 2 * tp / (2*tp + fp + fn)
		}
		weights[i] = rowSum / totalSamples
	}

	return averageScores(f1s, weights, average)
}

func averageScores(scores []float64, weights []float64, average string) (float64, error) {
	switch average {
	case "macro":
		sum := 0.0
		for _, s := range scores {
			sum += s
		}
		return sum / float64(len(scores)), nil
	case "weighted":
		sum := 0.0
		for i, s := range scores {
			sum += s * weights[i]
		}
		return sum, nil
	case "", "binary":
		if len(scores) > 2 {
			return 0, fmt.Errorf("average='%s' is not supported for multiclass data", average)
		}
		// For binary, return the score for the positive class (usually label 1)
		if len(scores) == 2 {
			return scores[1], nil
		}
		return scores[0], nil
	default:
		return 0, fmt.Errorf("unknown average type: %s", average)
	}
}
