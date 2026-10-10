package probability

import "math"

type ProbDist interface {
	Prob(sample string) float64
	Max() string
}

type MLEProbDist struct{ fd FreqDist }

func NewMLEProbDist(fd FreqDist) MLEProbDist { return MLEProbDist{fd: fd} }

func (p MLEProbDist) Prob(sample string) float64 {
	n := float64(p.fd.N())
	if n == 0 {
		return 0
	}
	return float64(p.fd[sample]) / n
}
func (p MLEProbDist) Max() string { return p.fd.Max() }

type LidstoneProbDist struct {
	fd    FreqDist
	gamma float64
	bins  int
}

func NewLidstoneProbDist(fd FreqDist, gamma float64, bins int) LidstoneProbDist {
	if bins <= 0 {
		bins = fd.B()
	}
	return LidstoneProbDist{fd: fd, gamma: gamma, bins: bins}
}

func (p LidstoneProbDist) Prob(sample string) float64 {
	return (float64(p.fd[sample]) + p.gamma) / (float64(p.fd.N()) + p.gamma*float64(p.bins))
}
func (p LidstoneProbDist) Max() string { return p.fd.Max() }

func NormalizeLogScores(scores map[string]float64) map[string]float64 {
	max := math.Inf(-1)
	for _, s := range scores {
		if s > max {
			max = s
		}
	}
	sum := 0.0
	out := map[string]float64{}
	for k, s := range scores {
		out[k] = math.Exp(s - max)
		sum += out[k]
	}
	for k := range out {
		out[k] /= sum
	}
	return out
}
