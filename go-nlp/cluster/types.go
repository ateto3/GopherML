package cluster

type Vector []float64

type Clusterer interface {
	Cluster([]Vector) []int
	Classify(Vector) int
}

func copyVec(v Vector) Vector { out := make(Vector, len(v)); copy(out, v); return out }
