package cluster

import "math"

func Euclidean(a, b Vector) float64 {
	s := 0.0
	for i := range a {
		d := a[i] - b[i]
		s += d * d
	}
	return math.Sqrt(s)
}

func centroid(vs []Vector) Vector {
	if len(vs) == 0 { return nil }
	out := make(Vector, len(vs[0]))
	for _, v := range vs { for i := range v { out[i] += v[i] } }
	for i := range out { out[i] /= float64(len(vs)) }
	return out
}
