package linear_model

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// generateData creates synthetic regression data.
func generateData(nSamples, nFeatures int, noise float64) ([][]float64, []float64) {
	rng := rand.New(rand.NewSource(42))

	X := make([][]float64, nSamples)
	for i := 0; i < nSamples; i++ {
		X[i] = make([]float64, nFeatures)
		for j := 0; j < nFeatures; j++ {
			X[i][j] = rng.Float64()*2 - 1 // Random values in [-1, 1]
		}
	}

	// True coefficients
	trueCoef := make([]float64, nFeatures)
	for j := 0; j < nFeatures; j++ {
		trueCoef[j] = rng.Float64()*2 - 1
	}

	// Generate y = X * coef + noise
	y := make([]float64, nSamples)
	for i := 0; i < nSamples; i++ {
		sum := 0.0
		for j := 0; j < nFeatures; j++ {
			sum += X[i][j] * trueCoef[j]
		}
		y[i] = sum + rng.NormFloat64()*noise
	}

	return X, y
}

// BenchmarkLinearRegressionVsParallel compara rendimiento secuencial vs paralelo.
func BenchmarkLinearRegressionVsParallel(b *testing.B) {
	sizes := []struct {
		samples  int
		features int
		desc     string
	}{
		{1000, 10, "Small"},
		{10000, 100, "Medium"},
		{50000, 200, "Large"},
	}

	for _, size := range sizes {
		X, y := generateData(size.samples, size.features, 0.1)

		b.Run(fmt.Sprintf("Sequential_%s_%dx%d", size.desc, size.samples, size.features), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				model := NewLinearRegression()
				model.Fit(X, y)
			}
		})

		b.Run(fmt.Sprintf("Parallel_%s_%dx%d", size.desc, size.samples, size.features), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				model := NewLinearRegressionParallel()
				model.Fit(X, y)
			}
		})
	}
}

// BenchmarkLassoVsParallel compara Lasso secuencial vs paralelo.
func BenchmarkLassoVsParallel(b *testing.B) {
	sizes := []struct {
		samples  int
		features int
		desc     string
	}{
		{1000, 20, "Small"},
		{5000, 50, "Medium"},
	}

	for _, size := range sizes {
		X, y := generateData(size.samples, size.features, 0.1)

		b.Run(fmt.Sprintf("Sequential_%s_%dx%d", size.desc, size.samples, size.features), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				model := NewLassoWithParams(LassoParams{Alpha: 0.1, MaxIter: 100})
				model.Fit(X, y)
			}
		})

		b.Run(fmt.Sprintf("Parallel_%s_%dx%d", size.desc, size.samples, size.features), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				model := NewLassoParallelWithParams(LassoParams{Alpha: 0.1, MaxIter: 100}, DefaultParallelOptions())
				model.Fit(X, y)
			}
		})
	}
}

// BenchmarkPredictionParallel compara predicción secuencial vs paralela.
func BenchmarkPredictionParallel(b *testing.B) {
	// Generate large dataset
	nTrain := 10000
	nFeatures := 100
	nTest := 10000

	XTrain, yTrain := generateData(nTrain, nFeatures, 0.1)
	XTest, _ := generateData(nTest, nFeatures, 0)

	// Train model once
	model := NewLinearRegressionParallel()
	model.Fit(XTrain, yTrain)

	b.Run("Sequential_Predict", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Use base sequential version
			baseModel := &LinearRegression{
				LinearModel: model.LinearModel,
				CopyX:       false,
			}
			baseModel.Predict(XTest)
		}
	})

	b.Run("Parallel_Predict", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			model.Predict(XTest)
		}
	})
}

// BenchmarkGramMatrixParallel mide rendimiento de multiplicación paralela.
func BenchmarkGramMatrixParallel(b *testing.B) {
	sizes := []int{100, 500, 1000}

	for _, n := range sizes {
		X, _ := generateData(n, n, 0)

		opts := ParallelOptions{NumWorkers: 4, MinBatchSize: 100}

		b.Run(fmt.Sprintf("Gram_%dx%d", n, n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ParallelMatrixMultiply(X, opts)
			}
		})
	}
}

// DemonstrateParallelSpeedup shows real speedup between sequential and parallel.
func DemonstrateParallelSpeedup() {
	// Generate medium dataset
	nSamples := 20000
	nFeatures := 50

	X, y := generateData(nSamples, nFeatures, 0.1)

	// Sequential timing
	start := time.Now()
	seqModel := NewLinearRegression()
	seqModel.Fit(X, y)
	seqTime := time.Since(start)

	// Parallel timing
	start = time.Now()
	parModel := NewLinearRegressionParallel()
	parModel.Fit(X, y)
	parTime := time.Since(start)

	speedup := float64(seqTime) / float64(parTime)

	fmt.Printf("Dataset: %d samples x %d features\n", nSamples, nFeatures)
	fmt.Printf("Sequential time: %v\n", seqTime)
	fmt.Printf("Parallel time: %v\n", parTime)
	fmt.Printf("Speedup: %.2fx\n", speedup)
}
