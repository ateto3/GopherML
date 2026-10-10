package probability

import "sync"

type FreqDist map[string]int

func NewFreqDist(samples ...string) FreqDist {
	fd := FreqDist{}
	for _, s := range samples {
		fd[s]++
	}
	return fd
}

func (fd FreqDist) Inc(sample string) { fd[sample]++ }
func (fd FreqDist) N() int            { n := 0; for _, c := range fd { n += c }; return n }
func (fd FreqDist) B() int            { return len(fd) }
func (fd FreqDist) Max() string {
	var best string
	bestN := -1
	for s, c := range fd {
		if c > bestN {
			best, bestN = s, c
		}
	}
	return best
}

func (fd FreqDist) Hapaxes() []string {
	out := make([]string, 0)
	for s, c := range fd {
		if c == 1 {
			out = append(out, s)
		}
	}
	return out
}

func NewFreqDistConcurrent(samples []string) FreqDist {
	fd := FreqDist{}
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, s := range samples {
		s := s
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			fd[s]++
			mu.Unlock()
		}()
	}
	wg.Wait()
	return fd
}
