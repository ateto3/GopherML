package model_selection

import (
	"fmt"
	"runtime"
	"sync"
)

// GridSearchCV performs an exhaustive search over specified parameter values for an estimator.
type GridSearchCV struct {
	EstimatorFactory func(params map[string]interface{}) Estimator
	ParamGrid        []map[string]interface{}
	CV               KFold
	BestParams       map[string]interface{}
	BestScore        float64
	CVResults        []CVResult // All results, including failed combinations
	NJobs            int        // -1 means all CPUs
}

// CVResult holds the result for a single parameter combination.
type CVResult struct {
	Params    map[string]interface{}
	MeanScore float64
	Err       error
}

// Fit runs cross-validation for all combinations in ParamGrid.
// Unlike the previous implementation, it does not abort on the first error —
// it collects all results and returns an error only if every combination failed.
func (g *GridSearchCV) Fit(X [][]float64, y []float64) error {
	nJobs := g.NJobs
	if nJobs <= 0 {
		nJobs = runtime.NumCPU()
	}
	if nJobs > len(g.ParamGrid) {
		nJobs = len(g.ParamGrid)
	}

	type result struct {
		params map[string]interface{}
		score  float64
		err    error
	}

	resultsChan := make(chan result, len(g.ParamGrid))
	paramsChan := make(chan map[string]interface{}, len(g.ParamGrid))

	// Send all param combinations upfront so workers can drain without blocking.
	for _, p := range g.ParamGrid {
		paramsChan <- p
	}
	close(paramsChan)

	var wg sync.WaitGroup
	for i := 0; i < nJobs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for params := range paramsChan {
				model := g.EstimatorFactory(params)
				scores, err := CrossValScore(model, X, y, g.CV)
				if err != nil {
					resultsChan <- result{params: params, err: err}
					continue
				}

				meanScore := 0.0
				for _, s := range scores {
					meanScore += s
				}
				meanScore /= float64(len(scores))
				resultsChan <- result{params: params, score: meanScore}
			}
		}()
	}

	// Close results channel once all workers finish.
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	g.BestScore = -1e15
	g.CVResults = make([]CVResult, 0, len(g.ParamGrid))
	failCount := 0

	for res := range resultsChan {
		g.CVResults = append(g.CVResults, CVResult{
			Params:    res.params,
			MeanScore: res.score,
			Err:       res.err,
		})
		if res.err != nil {
			failCount++
			continue
		}
		if res.score > g.BestScore {
			g.BestScore = res.score
			g.BestParams = res.params
		}
	}

	if failCount == len(g.ParamGrid) {
		return fmt.Errorf("grid search: all %d parameter combinations failed", len(g.ParamGrid))
	}
	if failCount > 0 {
		return fmt.Errorf("grid search: %d of %d parameter combinations failed (best score: %.4f)", failCount, len(g.ParamGrid), g.BestScore)
	}

	return nil
}
