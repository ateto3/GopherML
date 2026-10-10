package classify

import "sync"

// ApplyFeatures converts raw samples to labeled feature sets in parallel.
func ApplyFeatures[T any](samples []T, extractor func(T) Features, labeler func(T) string) []LabeledFeatures {
	if len(samples) == 0 {
		return nil
	}
	out := make([]LabeledFeatures, len(samples))
	var wg sync.WaitGroup
	for i := range samples {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			out[i] = LabeledFeatures{Features: extractor(samples[i]), Label: labeler(samples[i])}
		}()
	}
	wg.Wait()
	return out
}

// LabeledFeatures is a training example.
type LabeledFeatures struct {
	Features Features
	Label    string
}
