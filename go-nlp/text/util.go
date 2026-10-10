package text

import "sort"

func min(a, b int) int { if a < b { return a }; return b }
func max(a, b int) int { if a > b { return a }; return b }

func topN(scores map[string]int, n int) []string {
	type kv struct{ k string; v int }
	arr := make([]kv, 0, len(scores))
	for k, v := range scores { arr = append(arr, kv{k, v}) }
	sort.Slice(arr, func(i, j int) bool { if arr[i].v == arr[j].v { return arr[i].k < arr[j].k }; return arr[i].v > arr[j].v })
	if n > len(arr) { n = len(arr) }
	out := make([]string, 0, n)
	for _, x := range arr[:n] { out = append(out, x.k) }
	return out
}
