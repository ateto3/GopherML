package cluster

import "math"

// EM is a tiny diagonal-covariance Gaussian mixture.
type EM struct{ Means []Vector; Priors []float64; Vars []Vector }

func NewEM(means []Vector) *EM {
	k := len(means)
	pr := make([]float64, k)
	for i := range pr { pr[i] = 1 / float64(k) }
	vars := make([]Vector, k)
	for i := range vars { vars[i] = make(Vector, len(means[0])); for j := range vars[i] { vars[i][j] = 1 } }
	return &EM{Means: means, Priors: pr, Vars: vars}
}

func (e *EM) Cluster(vs []Vector) []int {
	assign := make([]int, len(vs))
	for it := 0; it < 10; it++ {
		for i, v := range vs { assign[i] = e.Classify(v) }
		k := len(e.Means)
		groups := make([][]Vector, k)
		for i, c := range assign { groups[c] = append(groups[c], vs[i]) }
		for j := 0; j < k; j++ {
			if len(groups[j]) == 0 { continue }
			e.Means[j] = centroid(groups[j])
			for d := range e.Vars[j] { e.Vars[j][d] = 1 }
			e.Priors[j] = float64(len(groups[j])) / float64(len(vs))
		}
	}
	return assign
}

func (e *EM) Classify(v Vector) int {
	best, bi := math.Inf(-1), 0
	for j := range e.Means {
		lp := math.Log(e.Priors[j])
		for d := range v {
			diff := v[d] - e.Means[j][d]
			lp += -0.5 * diff * diff / e.Vars[j][d]
		}
		if lp > best { best, bi = lp, j }
	}
	return bi
}
