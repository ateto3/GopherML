// Package linear_model provides parallel implementations optimized for Go's concurrency model.
// This file contains low-level parallel primitives for matrix operations.
package linear_model

import (
	"fmt"
	"runtime"
	"sync"
)

// ParallelOptions configures parallel execution.
type ParallelOptions struct {
	// NumWorkers is the number of goroutines to use. If 0, uses runtime.NumCPU().
	NumWorkers int
	// MinBatchSize is the minimum number of elements per worker to avoid overhead.
	MinBatchSize int
}

// DefaultParallelOptions returns sensible defaults for parallel execution.
func DefaultParallelOptions() ParallelOptions {
	return ParallelOptions{
		NumWorkers:   runtime.NumCPU(),
		MinBatchSize: 1000,
	}
}

// ParallelFor executes a loop in parallel using goroutines.
// The work is divided evenly among workers.
func ParallelFor(n int, opts ParallelOptions, fn func(start, end int)) {
	if opts.NumWorkers <= 0 {
		opts.NumWorkers = runtime.NumCPU()
	}
	if opts.MinBatchSize <= 0 {
		opts.MinBatchSize = 1000
	}

	// Sequential execution for small workloads
	if n < opts.MinBatchSize || opts.NumWorkers == 1 {
		fn(0, n)
		return
	}

	workers := opts.NumWorkers
	if n < workers*opts.MinBatchSize {
		workers = n / opts.MinBatchSize
		if workers < 1 {
			workers = 1
		}
	}

	chunkSize := (n + workers - 1) / workers // Ceiling division

	var wg sync.WaitGroup
	wg.Add(workers)

	for w := 0; w < workers; w++ {
		start := w * chunkSize
		end := start + chunkSize
		if end > n {
			end = n
		}
		if start >= n {
			wg.Done()
			continue
		}

		go func(s, e int) {
			defer wg.Done()
			fn(s, e)
		}(start, end)
	}

	wg.Wait()
}

// ParallelMatrixMultiply computes X^T * X (Gram matrix) in parallel.
// Uses pointer arithmetic to minimize allocations.
func ParallelMatrixMultiply(X [][]float64, opts ParallelOptions) [][]float64 {
	nSamples := len(X)
	if nSamples == 0 {
		return [][]float64{}
	}
	nFeatures := len(X[0])

	// Result matrix initialized to zero
	result := make([][]float64, nFeatures)
	for i := range result {
		result[i] = make([]float64, nFeatures)
	}

	// Use mutex for thread-safe accumulation
	var mu sync.Mutex

	ParallelFor(nSamples, opts, func(start, end int) {
		// Local buffer to reduce mutex contention
		localResult := make([][]float64, nFeatures)
		for i := range localResult {
			localResult[i] = make([]float64, nFeatures)
		}

		// Compute local contribution
		for i := start; i < end; i++ {
			xi := X[i] // Pointer to row, avoids bounds checks in inner loop
			for j := 0; j < nFeatures; j++ {
				for k := j; k < nFeatures; k++ {
					localResult[j][k] += xi[j] * xi[k]
				}
			}
		}

		// Accumulate with global lock
		mu.Lock()
		for j := 0; j < nFeatures; j++ {
			for k := j; k < nFeatures; k++ {
				result[j][k] += localResult[j][k]
				if j != k {
					result[k][j] += localResult[j][k]
				}
			}
		}
		mu.Unlock()
	})

	return result
}

// ParallelMatVecMul computes X * v (matrix-vector product) in parallel.
func ParallelMatVecMul(X [][]float64, v []float64, opts ParallelOptions) []float64 {
	nSamples := len(X)
	if nSamples == 0 {
		return []float64{}
	}
	nFeatures := len(X[0])

	result := make([]float64, nSamples)

	ParallelFor(nSamples, opts, func(start, end int) {
		for i := start; i < end; i++ {
			sum := 0.0
			row := X[i]
			// Manual loop unrolling for better performance
			j := 0
			for ; j <= nFeatures-4; j += 4 {
				sum += row[j]*v[j] + row[j+1]*v[j+1] + row[j+2]*v[j+2] + row[j+3]*v[j+3]
			}
			for ; j < nFeatures; j++ {
				sum += row[j] * v[j]
			}
			result[i] = sum
		}
	})

	return result
}

// ParallelCenterData centers X and y by subtracting means, in parallel.
// Modifies data in-place using pointers to avoid allocations.
func ParallelCenterData(X [][]float64, y []float64, opts ParallelOptions) (xOffset []float64, yOffset float64) {
	nSamples := len(X)
	if nSamples == 0 {
		return []float64{}, 0
	}
	nFeatures := len(X[0])

	// Compute column means in parallel
	xOffset = make([]float64, nFeatures)

	// Parallel column mean computation
	type meanResult struct {
		idx  int
		mean float64
	}
	meanChan := make(chan meanResult, nFeatures)

	// Worker pool for computing means
	var meanWg sync.WaitGroup
	numWorkers := opts.NumWorkers
	if numWorkers > nFeatures {
		numWorkers = nFeatures
	}

	colsPerWorker := (nFeatures + numWorkers - 1) / numWorkers

	for w := 0; w < numWorkers; w++ {
		startCol := w * colsPerWorker
		endCol := startCol + colsPerWorker
		if endCol > nFeatures {
			endCol = nFeatures
		}
		if startCol >= nFeatures {
			break
		}

		meanWg.Add(1)
		go func(start, end int) {
			defer meanWg.Done()
			for j := start; j < end; j++ {
				sum := 0.0
				for i := 0; i < nSamples; i++ {
					sum += X[i][j]
				}
				meanChan <- meanResult{idx: j, mean: sum / float64(nSamples)}
			}
		}(startCol, endCol)
	}

	// Close channel when all workers done
	go func() {
		meanWg.Wait()
		close(meanChan)
	}()

	// Collect results
	for res := range meanChan {
		xOffset[res.idx] = res.mean
	}

	// Compute y mean
	ySum := 0.0
	for i := 0; i < nSamples; i++ {
		ySum += y[i]
	}
	yOffset = ySum / float64(nSamples)

	// Center X in parallel using workers that modify data in-place
	ParallelFor(nSamples, opts, func(start, end int) {
		for i := start; i < end; i++ {
			row := X[i] // Pointer to row
			for j := 0; j < nFeatures; j++ {
				row[j] -= xOffset[j]
			}
		}
	})

	// Center y in parallel
	ParallelFor(nSamples, opts, func(start, end int) {
		for i := start; i < end; i++ {
			y[i] -= yOffset
		}
	})

	return xOffset, yOffset
}

// WorkerPool provides a reusable worker pool for batch processing.
type WorkerPool struct {
	workers int
	tasks   chan func()
	wg      sync.WaitGroup
}

// NewWorkerPool creates a new worker pool with the specified number of workers.
func NewWorkerPool(workers int) *WorkerPool {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	pool := &WorkerPool{
		workers: workers,
		tasks:   make(chan func(), workers*2),
	}

	// Start worker goroutines
	for i := 0; i < workers; i++ {
		go pool.worker()
	}

	return pool
}

// worker is the goroutine that processes tasks.
func (p *WorkerPool) worker() {
	for task := range p.tasks {
		task()
		p.wg.Done()
	}
}

// Submit adds a task to the pool.
func (p *WorkerPool) Submit(task func()) {
	p.wg.Add(1)
	p.tasks <- task
}

// Wait waits for all submitted tasks to complete.
func (p *WorkerPool) Wait() {
	p.wg.Wait()
}

// Close shuts down the worker pool.
func (p *WorkerPool) Close() {
	close(p.tasks)
}

// LinearRegressionParallel is a parallel implementation of Ordinary Least Squares.
type LinearRegressionParallel struct {
	LinearRegression
	ParallelOpts ParallelOptions
}

// NewLinearRegressionParallel creates a new parallel LinearRegression model.
func NewLinearRegressionParallel() *LinearRegressionParallel {
	return &LinearRegressionParallel{
		LinearRegression: *NewLinearRegression(),
		ParallelOpts:     DefaultParallelOptions(),
	}
}

// Fit trains the linear regression model using parallel operations.
func (lr *LinearRegressionParallel) Fit(X [][]float64, y []float64) error {
	if len(X) == 0 || len(y) == 0 {
		return ErrEmptyData
	}
	if len(X) != len(y) {
		return ErrDimensionMismatch
	}

	nSamples := len(X)
	nFeatures := len(X[0])
	lr.NFeaturesIn = nFeatures

	// Create working copy
	XWork := make([][]float64, nSamples)
	dataBlock := make([]float64, nSamples*nFeatures)
	for i := range X {
		start := i * nFeatures
		end := start + nFeatures
		XWork[i] = dataBlock[start:end]
		copy(XWork[i], X[i])
	}
	yWork := make([]float64, nSamples)
	copy(yWork, y)

	var xOffset []float64
	var yOffset float64

	if lr.FitIntercept {
		xOffset, yOffset = ParallelCenterData(XWork, yWork, lr.ParallelOpts)
	}

	// Solve using parallel normal equations
	// (X^T X) w = X^T y
	gram := ParallelMatrixMultiply(XWork, lr.ParallelOpts)
	xty := ParallelMatVecMul(Transpose(XWork), yWork, lr.ParallelOpts)

	coef, err := gaussianElimination(gram, xty)
	if err != nil {
		return err
	}

	lr.Coef = coef

	if lr.FitIntercept {
		lr.LinearModel.setIntercept(xOffset, yOffset)
	}

	lr.IsFitted = true
	return nil
}

// Predict makes predictions in parallel.
func (lr *LinearRegressionParallel) Predict(X [][]float64) ([]float64, error) {
	if !lr.IsFitted {
		return nil, ErrNotFitted
	}
	if len(X) == 0 {
		return nil, ErrEmptyData
	}
	// Validate feature count
	for i, row := range X {
		if len(row) != lr.NFeaturesIn {
			return nil, fmt.Errorf("row %d has %d features, expected %d", i, len(row), lr.NFeaturesIn)
		}
	}
	result := ParallelMatVecMul(X, lr.Coef, lr.ParallelOpts)
	// Add intercept (omitted in the original implementation — correctness bug fix)
	if lr.Intercept != 0 {
		for i := range result {
			result[i] += lr.Intercept
		}
	}
	return result, nil
}

// Transpose returns the transpose of a matrix.
func Transpose(X [][]float64) [][]float64 {
	if len(X) == 0 {
		return [][]float64{}
	}
	nRows := len(X)
	nCols := len(X[0])
	XT := make([][]float64, nCols)
	for i := range XT {
		XT[i] = make([]float64, nRows)
		for j := 0; j < nRows; j++ {
			XT[i][j] = X[j][i]
		}
	}
	return XT
}
