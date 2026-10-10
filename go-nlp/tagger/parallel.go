package tagger

import "sync"

// ParallelPosTag tags tokens concurrently using the heuristic tagger.
func ParallelPosTag(tokens []string) []TaggedToken {
	return ParallelPosTagWithWorkerCount(tokens, 0)
}

// ParallelPosTagWithWorkerCount tags tokens concurrently with a worker pool.
func ParallelPosTagWithWorkerCount(tokens []string, workers int) []TaggedToken {
	if len(tokens) == 0 {
		return nil
	}
	if workers <= 0 {
		workers = 4
	}
	if workers > len(tokens) {
		workers = len(tokens)
	}

	results := make([]TaggedToken, len(tokens))
	jobs := make(chan int, len(tokens))
	var wg sync.WaitGroup
	tagger := NewHeuristicTagger()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				results[idx] = TaggedToken{Word: tokens[idx], Tag: tagger.tagToken(tokens[idx], idx, tokens)}
			}
		}()
	}

	for i := range tokens {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	return results
}
