package cluster

// GAAC is a simple agglomerative clusterer using centroid distance.
type GAAC struct{ K int }

func NewGAAC(k int) *GAAC { return &GAAC{K: k} }

func (g *GAAC) Cluster(vs []Vector) []int {
	assign := make([]int, len(vs))
	for i := range assign { assign[i] = i }
	for len(unique(assign)) > g.K {
		a, b := 0, 1
		best := Euclidean(vs[0], vs[1])
		for i := 0; i < len(vs); i++ {
			for j := i + 1; j < len(vs); j++ {
				if d := Euclidean(vs[i], vs[j]); d < best { best, a, b = d, i, j }
			}
		}
		for i := range assign { if assign[i] == b { assign[i] = a } }
	}
	return assign
}

func (g *GAAC) Classify(v Vector) int { return 0 }

func unique(xs []int) []int {
	m := map[int]struct{}{}
	for _, x := range xs { m[x] = struct{}{} }
	out := make([]int, 0, len(m))
	for x := range m { out = append(out, x) }
	return out
}
