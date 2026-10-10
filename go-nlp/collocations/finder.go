package collocations

import (
	"sort"
	"sync"

	"github.com/ateto3/gopherml/go-nlp/probability"
)

type BigramFinder struct{ W, B probability.FreqDist }
type TrigramFinder struct{ W, B, T probability.FreqDist }

func NewBigramFinder(words []string) BigramFinder {
	w, b := probability.FreqDist{}, probability.FreqDist{}
	for i := 0; i+1 < len(words); i++ {
		w[words[i]]++
		b[words[i]+"|"+words[i+1]]++
	}
	if len(words) > 0 {
		w[words[len(words)-1]]++
	}
	return BigramFinder{W: w, B: b}
}

func NewTrigramFinder(words []string) TrigramFinder {
	w, b, t := probability.FreqDist{}, probability.FreqDist{}, probability.FreqDist{}
	for i := 0; i+1 < len(words); i++ {
		w[words[i]]++
		b[words[i]+"|"+words[i+1]]++
		if i+2 < len(words) {
			t[words[i]+"|"+words[i+1]+"|"+words[i+2]]++
		}
	}
	if len(words) > 0 {
		w[words[len(words)-1]]++
	}
	return TrigramFinder{W: w, B: b, T: t}
}

func (f BigramFinder) Filter(min int) BigramFinder {
	out := probability.FreqDist{}
	for k, v := range f.B {
		if v >= min {
			out[k] = v
		}
	}
	f.B = out
	return f
}

func (f BigramFinder) PMI() []ScoredNGram {
	n := float64(f.W.N())
	out := make([]ScoredNGram, 0, len(f.B))
	for k, c := range f.B {
		parts := split2(k)
		w1, w2 := parts[0], parts[1]
		pxy := float64(c) / n
		px := float64(f.W[w1]) / n
		py := float64(f.W[w2]) / n
		if pxy > 0 && px > 0 && py > 0 {
			out = append(out, ScoredNGram{NGram: []string{w1, w2}, Score: log2(pxy / (px * py))})
		}
	}
	sortScored(out)
	return out
}

func (f TrigramFinder) Top(n int) []ScoredNGram {
	out := make([]ScoredNGram, 0, len(f.T))
	for k, c := range f.T {
		out = append(out, ScoredNGram{NGram: split3(k), Score: float64(c)})
	}
	sortScored(out)
	if n > len(out) {
		n = len(out)
	}
	return out[:n]
}

type ScoredNGram struct {
	NGram []string
	Score float64
}

func sortScored(s []ScoredNGram) {
	sort.Slice(s, func(i, j int) bool {
		if s[i].Score == s[j].Score {
			return join(s[i].NGram) < join(s[j].NGram)
		}
		return s[i].Score > s[j].Score
	})
}

func BigramFinderConcurrent(docs [][]string) BigramFinder {
	var mu sync.Mutex
	w, b := probability.FreqDist{}, probability.FreqDist{}
	var wg sync.WaitGroup
	for _, doc := range docs {
		doc := doc
		wg.Add(1)
		go func() {
			defer wg.Done()
			local := NewBigramFinder(doc)
			mu.Lock()
			for k, v := range local.W {
				w[k] += v
			}
			for k, v := range local.B {
				b[k] += v
			}
			mu.Unlock()
		}()
	}
	wg.Wait()
	return BigramFinder{W: w, B: b}
}
