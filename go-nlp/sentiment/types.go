package sentiment

// Score is a sentiment result.
type Score struct {
	Neg  float64
	Neu  float64
	Pos  float64
	Comp float64
}

// Analyzer analyzes sentiment.
type Analyzer interface {
	Analyze(text string) Score
}
