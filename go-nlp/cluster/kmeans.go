package cluster

// KMeans is a compact k-means clusterer.
type KMeans struct{ K int; Means []Vector; MaxIter int }

func NewKMeans(k int) *KMeans { return &KMeans{K: k, MaxIter: 50} }

func (k *KMeans) Cluster(vs []Vector) []int {
	if len(vs) == 0 { return nil }
	if len(k.Means) == 0 { k.Means = append([]Vector(nil), vs[:k.K]...) }
	assign := make([]int, len(vs))
	for it := 0; it < k.MaxIter; it++ {
		changed := false
		clusters := make([][]Vector, k.K)
		for i, v := range vs {
			best, bi := Euclidean(v, k.Means[0]), 0
			for j := 1; j < k.K; j++ { if d := Euclidean(v, k.Means[j]); d < best { best, bi = d, j } }
			if assign[i] != bi { changed = true }
			assign[i] = bi
			clusters[bi] = append(clusters[bi], v)
		}
		for i := 0; i < k.K; i++ { if len(clusters[i]) > 0 { k.Means[i] = centroid(clusters[i]) } }
		if !changed { break }
	}
	return assign
}

func (k *KMeans) Classify(v Vector) int {
	best, bi := Euclidean(v, k.Means[0]), 0
	for j := 1; j < len(k.Means); j++ { if d := Euclidean(v, k.Means[j]); d < best { best, bi = d, j } }
	return bi
}
