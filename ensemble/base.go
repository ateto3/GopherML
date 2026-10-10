package ensemble

// BaseForest contains common attributes for random forest models.
type BaseForest struct {
	NEstimators      int
	MaxDepth         int
	MinSamplesSplit  int
	MinSamplesLeaf   int
	MaxFeatures      interface{}
	Bootstrap        bool
	RandomState      int64
	NFeatures        int
	IsFitted         bool
}

// ForestParams contains parameters for creating a random forest.
type ForestParams struct {
	NEstimators     int
	MaxDepth        int
	MinSamplesSplit int
	MinSamplesLeaf  int
	MaxFeatures     interface{}
	Bootstrap       bool
	RandomState     int64
}

// DefaultForestParams returns default parameters for random forest.
func DefaultForestParams() ForestParams {
	return ForestParams{
		NEstimators:     100,
		MaxDepth:        -1,
		MinSamplesSplit: 2,
		MinSamplesLeaf:  1,
		MaxFeatures:     "sqrt",
		Bootstrap:       true,
		RandomState:     -1,
	}
}
