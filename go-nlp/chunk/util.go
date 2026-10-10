package chunk

import "math"

// ChunkScore tracks chunking quality in a simple way.
type ChunkScore struct {
	Correct int
	Guess   int
	Total   int
}

// Score compares gold and predicted trees.
func (cs *ChunkScore) Score(gold, guess Tree) {
	g := len(gold.Leaves())
	p := len(guess.Leaves())
	cs.Total += g
	cs.Guess += p
	cs.Correct += min(g, p)
}

// Precision returns precision.
func (cs ChunkScore) Precision() float64 {
	if cs.Guess == 0 {
		return 0
	}
	return float64(cs.Correct) / float64(cs.Guess)
}

// Recall returns recall.
func (cs ChunkScore) Recall() float64 {
	if cs.Total == 0 {
		return 0
	}
	return float64(cs.Correct) / float64(cs.Total)
}

// FMeasure returns F1 with alpha=0.5 semantics.
func (cs ChunkScore) FMeasure() float64 {
	p, r := cs.Precision(), cs.Recall()
	if p == 0 || r == 0 {
		return 0
	}
	return 2 * p * r / (p + r)
}

func min(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}
