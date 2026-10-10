package classify

// Features is a sparse feature set.
type Features map[string]any

// Classifier labels feature sets with a single class.
type Classifier interface {
	Labels() []string
	Classify(features Features) string
	ClassifyMany(featuresets []Features) []string
}
