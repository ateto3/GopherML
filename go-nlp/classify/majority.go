package classify

// MajorityClassifier always predicts the most common label.
type MajorityClassifier struct {
	label string
}

// TrainMajority creates a majority classifier from labeled data.
func TrainMajority(data []LabeledFeatures) *MajorityClassifier {
	counts := map[string]int{}
	best := ""
	for _, ex := range data {
		counts[ex.Label]++
		if counts[ex.Label] > counts[best] {
			best = ex.Label
		}
	}
	return &MajorityClassifier{label: best}
}

func (c *MajorityClassifier) Labels() []string { return []string{c.label} }

func (c *MajorityClassifier) Classify(Features) string { return c.label }

func (c *MajorityClassifier) ClassifyMany(featuresets []Features) []string {
	out := make([]string, len(featuresets))
	for i := range out {
		out[i] = c.label
	}
	return out
}
