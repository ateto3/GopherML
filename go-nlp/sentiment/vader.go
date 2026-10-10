package sentiment

import (
	"math"
	"strings"
)

var lexicon = map[string]float64{
	"good": 2.0, "great": 3.0, "excellent": 3.5, "amazing": 3.2, "fast": 1.5,
	"bad": -2.0, "terrible": -3.5, "awful": -3.2, "slow": -1.5, "hate": -3.0,
	"love": 3.0, "happy": 2.5, "sad": -2.5, "perfect": 3.2,
}

var boosters = map[string]float64{"very": 0.3, "extremely": 0.5, "really": 0.2, "slightly": -0.2}

// VaderAnalyzer is a compact rule-based analyzer.
type VaderAnalyzer struct{}

// NewVaderAnalyzer creates a new analyzer.
func NewVaderAnalyzer() VaderAnalyzer { return VaderAnalyzer{} }

// Analyze scores a sentence.
func (VaderAnalyzer) Analyze(text string) Score {
	toks := strings.Fields(strings.ToLower(text))
	var sum float64
	for i, tok := range toks {
		val := lexicon[tok]
		if i > 0 {
			if b, ok := boosters[toks[i-1]]; ok {
				val += math.Copysign(b, val)
			}
			if toks[i-1] == "not" || strings.HasSuffix(toks[i-1], "n't") {
				val *= -0.74
			}
		}
		sum += val
	}
	comp := sum / math.Sqrt(sum*sum+15)
	pos, neg := 0.0, 0.0
	if sum > 0 {
		pos = comp
	} else if sum < 0 {
		neg = -comp
	}
	neu := 1 - math.Abs(comp)
	return Score{Neg: neg, Neu: neu, Pos: pos, Comp: comp}
}
