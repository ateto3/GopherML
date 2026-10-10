package metrics

import (
	"fmt"
	"strings"
)

// ClassificationReport generates a detailed text report of classification metrics.
func ClassificationReport(yTrue, yPred []float64) (string, error) {
	cm, labels, err := ConfusionMatrix(yTrue, yPred)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-12s %-10s %-10s %-10s %-10s\n", "Label", "Precision", "Recall", "F1-Score", "Support"))
	sb.WriteString(strings.Repeat("-", 55) + "\n")

	totalSupport := 0
	for _, row := range cm {
		for _, val := range row {
			totalSupport += val
		}
	}

	for i, label := range labels {
		tp := float64(cm[i][i])
		fp := 0.0
		fn := 0.0
		rowSum := 0.0
		for j := range labels {
			if i != j {
				fp += float64(cm[j][i])
				fn += float64(cm[i][j])
			}
			rowSum += float64(cm[i][j])
		}

		precision := 0.0
		if tp+fp > 0 {
			precision = tp / (tp + fp)
		}
		recall := 0.0
		if tp+fn > 0 {
			recall = tp / (tp + fn)
		}
		f1 := 0.0
		if precision+recall > 0 {
			f1 = 2 * precision * recall / (precision + recall)
		}

		sb.WriteString(fmt.Sprintf("%-12.2f %-10.4f %-10.4f %-10.4f %-10d\n", 
			label, precision, recall, f1, int(rowSum)))
	}

	sb.WriteString(strings.Repeat("-", 55) + "\n")
	
	// Add macro and weighted averages
	acc, _ := AccuracyScore(yTrue, yPred)
	sb.WriteString(fmt.Sprintf("%-34s %-10.4f %-10d\n", "Accuracy", acc, totalSupport))

	macroF1, _ := F1Score(yTrue, yPred, "macro")
	sb.WriteString(fmt.Sprintf("%-34s %-10.4f %-10d\n", "Macro Avg (F1)", macroF1, totalSupport))

	weightedF1, _ := F1Score(yTrue, yPred, "weighted")
	sb.WriteString(fmt.Sprintf("%-34s %-10.4f %-10d\n", "Weighted Avg (F1)", weightedF1, totalSupport))

	return sb.String(), nil
}
