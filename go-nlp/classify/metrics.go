package classify

// Accuracy computes classification accuracy.
func Accuracy(gold, pred []string) float64 {
	if len(gold) == 0 || len(gold) != len(pred) {
		return 0
	}
	correct := 0
	for i := range gold {
		if gold[i] == pred[i] {
			correct++
		}
	}
	return float64(correct) / float64(len(gold))
}
