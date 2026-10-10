package lm

import "sync"

type NgramCounter struct{ Unigrams map[string]int; Bigrams map[string]map[string]int }

func NewCounter() *NgramCounter { return &NgramCounter{Unigrams: map[string]int{}, Bigrams: map[string]map[string]int{}} }

func (c *NgramCounter) Update(sentences []Sentence) {
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, s := range sentences {
		s := s
		wg.Add(1)
		go func() {
			defer wg.Done()
			localU := map[string]int{}
			localB := map[string]map[string]int{}
			for i, w := range s {
				localU[w]++
				if i+1 < len(s) {
					if localB[w] == nil { localB[w] = map[string]int{} }
					localB[w][s[i+1]]++
				}
			}
			mu.Lock()
			for k, v := range localU { c.Unigrams[k] += v }
			for k, m2 := range localB {
				if c.Bigrams[k] == nil { c.Bigrams[k] = map[string]int{} }
				for k2, v := range m2 { c.Bigrams[k][k2] += v }
			}
			mu.Unlock()
		}()
	}
	wg.Wait()
}
