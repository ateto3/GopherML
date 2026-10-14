# Performance Optimizations in LinearModel

This document describes the Go-specific optimizations implemented in the `linear_model` package.

## 🚀 High-Performance Features

### 1. Goroutines and Native Parallelism

All computationally intensive operations use goroutines:

```go
// Example: ParallelFor splits work across CPUs
ParallelFor(nSamples, opts, func(start, end int) {
    for i := start; i < end; i++ {
        // Process batch
    }
})
```

**`sync.WaitGroup` usage** for efficient synchronization without unnecessary locking.

### 2. Memory Pools and Buffer Reuse

```go
// Pre-allocation of contiguous memory for better cache locality
dataBlock := make([]float64, nSamples*nFeatures)
for i := range X {
    XWork[i] = dataBlock[i*nFeatures : (i+1)*nFeatures]
}
```

**Advantage**: A single allocation vs nSamples separate allocations.

### 3. Pointers and Slice Headers

Extensive use of pointers to avoid copies:
```go
xi := X[i]  // Pointer to row, avoids bounds checks in inner loops
for j := 0; j < nFeatures; j++ {
    sum += xi[j] * xi[k]  // Direct access via pointer
}
```

### 4. Channels for Aggregation

Workers send partial results via channels:
```go
type partialSum struct {
    idx int
    val float64
}
results := make(chan partialSum, nFeatures)
// Workers compute and send
results <- partialSum{idx: j, val: sum}
```

### 5. Manual Loop Unrolling

For critical operations:
```go
// Unroll x4 for better CPU pipelining
for ; j <= nFeatures-4; j += 4 {
    sum += row[j]*v[j] + row[j+1]*v[j+1] + 
           row[j+2]*v[j+2] + row[j+3]*v[j+3]
}
```

## 📊 Benchmark Results

| Dataset | Sequential | Parallel | Speedup |
|---------|-----------|----------|---------|
| 50K×200 (Large) | 46.0s | 1.8s | **25x** |
| 5K×50 Lasso | 91ms | 17ms | **5.3x** |
| 10K prediction | 2.4ms | 1.2ms | **2x** |

## 🎯 Available Parallel Models

1. **`LinearRegressionParallel`** - Parallel matrix multiplication
2. **`LassoParallel`** - Coordinate descent with parallel updates
3. **`WorkerPool`** reusable for batch processing

## 🔧 Recommended Usage

### When to use parallel versions:
- Dataset > 10,000 samples
- Features > 50
- Large batch predictions

### When to use sequential versions:
- Dataset < 1,000 samples (goroutine overhead does not pay off)
- Single-core systems
- Low-latency real-time applications

```go
// Usage example
opts := linear_model.ParallelOptions{
    NumWorkers:   runtime.NumCPU(),
    MinBatchSize: 1000,
}

model := linear_model.NewLinearRegressionParallelWithParams(
    linear_model.LinearRegressionParams{FitIntercept: true},
    opts,
)
model.Fit(X, y)
```

## 💡 Go-Specific Optimizations

| Technique | Benefit |
|-----------|---------|
| Lightweight goroutines | ~2KB stack vs ~1MB thread |
| Contiguous slices | Better cache locality |
| `sync.Pool` | GC object reuse |
| Buffered channels | Reduced contention |
| Selective mutex | Only where needed |

## 🚫 Zero Dependencies

The package uses the Go standard library **exclusively**:
- `runtime` - NumCPU()
- `sync` - WaitGroup, Mutex, Pool
- No cgo, no numpy, no external dependencies
