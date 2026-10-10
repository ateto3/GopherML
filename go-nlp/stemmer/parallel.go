package stemmer

import "sync"

// StemAll stems words concurrently.
func StemAll(words []string, s Stemmer) []string {
	if len(words) == 0 {
		return nil
	}
	out := make([]string, len(words))
	var wg sync.WaitGroup
	for i := range words {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			out[i] = s.Stem(words[i])
		}()
	}
	wg.Wait()
	return out
}
