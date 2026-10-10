package metrics

import (
	"fmt"
	"sort"
)

// ROCAUCScore calculates the Area Under the Receiver Operating Characteristic Curve.
func ROCAUCScore(yTrue []float64, yScores []float64) (float64, error) {
	if len(yTrue) != len(yScores) {
		return 0, fmt.Errorf("ROCAUCScore: yTrue and yScores have different lengths (%d vs %d)", len(yTrue), len(yScores))
	}
	if len(yTrue) == 0 {
		return 0, fmt.Errorf("ROCAUCScore: empty input")
	}

	// Simple AUC implementation based on rank sum
	n := len(yTrue)
	type pair struct {
		score float64
		label float64
	}
	pairs := make([]pair, n)
	nPos := 0.0
	nNeg := 0.0
	for i := range yTrue {
		pairs[i] = pair{yScores[i], yTrue[i]}
		if yTrue[i] == 1 {
			nPos++
		} else {
			nNeg++
		}
	}

	if nPos == 0 || nNeg == 0 {
		return 0.5, nil // Undefined AUC — return chance level
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].score < pairs[j].score
	})

	rankSum := 0.0
	for i := 0; i < n; i++ {
		if pairs[i].label == 1 {
			rankSum += float64(i + 1)
		}
	}

	auc := (rankSum - (nPos*(nPos+1))/2.0) / (nPos * nNeg)
	return auc, nil
}
