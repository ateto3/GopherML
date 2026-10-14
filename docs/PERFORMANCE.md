# ⚡ Performance Guide

GopherML is designed for low-latency inference in production (trading bots, microservices). This document collects available optimizations and **how to use them correctly**.

---

## 1. Golden rule: the hot path

In a high-frequency loop (one tick per candle, thousands of predictions/second):

```go
// ❌ Avoid: reserves memory for an entire batch each call
preds, _ := rf.Predict([][]float64{sample})

// ✅ Correct: single sample API, no extra allocations
label, _ := rf.PredictSingle(sample)
proba, _ := rf.PredictProbaSingle(sample)

// ✅ In the crypto pipeline
signal, conf, err := p.PredictSignal(lastCandles)
```

| API | Recommended use |
| :--- | :--- |
| `PredictSingle` / `PredictProbaSingle` | Trading loops, one sample per tick |
| `Predict` / `PredictProba` | Batch evaluation, backtesting, dataset scoring |
| `utils.CheckFiniteSingle` | Single sample validation in hot path |

---

## 2. Parallelism: when yes and when no

| Situation | Recommendation |
| :--- | :--- |
| Dataset > 10,000 samples or > 50 features | Parallel versions (`LinearRegressionParallel`, `LassoParallel`) |
| `GridSearchCV` with many combinations | `NJobs: -1` (all cores) |
| Random Forest training | Already parallel by default (one goroutine per tree) |
| Dataset < 1,000 samples | Sequential: goroutine overhead doesn't pay off |
| Single-core (e.g., container with 1 CPU) | Sequential |
| Single sample prediction | **Never** parallelize: `PredictSingle` |

```go
opts := linear_model.ParallelOptions{
    NumWorkers:   runtime.NumCPU(),
    MinBatchSize: 1000, // don't split small batches
}
model := linear_model.NewLinearRegressionParallel()
model.Fit(X, y)
```

---

## 3. GopherML internal optimizations

| Technique | Where | Benefit |
| :--- | :--- | :--- |
| Goroutines + `sync.WaitGroup` | `ensemble`, `model_selection`, `linear_model` | Leverages all cores |
| Pre-allocated contiguous memory | `linear_model` (flat blocks) | Single allocation vs N; better cache locality |
| Row pointers in inner loops | `linear_model` | Avoids *bounds checks* in hot loops |
| `sync.Pool` / reusable buffers | `linear_model`, `crypto` | Less GC pressure |
| *Loop unrolling* (x4) | critical multiplications | Better CPU *pipelining* |
| Indexed nodes (`[]Node` + indices) | `tree` | Fewer allocations than pointer graphs |
| Iterative DFS with explicit stack | `tree.DepthFirstTreeBuilder` | Zero risk of stack overflow |
| `[]byte` / `strings.Builder` | `go-nlp` | *Zero-copy* and less GC |

Detailed `linear_model` benchmarks: **[../linear_model/PERFORMANCE.md](../linear_model/PERFORMANCE.md)** (e.g., 50K×200: 46.0s → 1.8s sequential vs parallel).

---

## 4. How to measure yourself

```bash
# Benchmarks for a package
go test ./linear_model -bench=. -benchmem

# Tests + benchmark with CPU profiling
go test ./ensemble -bench=. -cpuprofile=cpu.out
go tool pprof cpu.out

# Race condition detector (essential before production)
go test ./... -race
```

Indicators to watch:

- **`allocs/op`** in `PredictSingle`: should be 0 or almost 0.
- **`B/op`**: allocations per operation; if it rises after a change, there's a regression.
- **p99** time of `PredictSignal` on your target hardware (Raspberry Pi ARM64, x86_64).

---

## 5. Production checklist for bots

| Practice | Reason |
| :--- | :--- |
| Load model **at startup**, not on every prediction | `Load` is expensive: do it once |
| Use `MinConfidence >= 0.60` | Reduces false signals in sideways markets |
| Validate candles before predicting | Exchanges can send candles with gaps or volume 0 |
| Save scaler next to model | `CryptoPredictor.SaveModel` does this automatically |
| Fixed `RandomState` in training | Reproducible results |
| Train on ARM64 if deploying on ARM | Avoids FP differences between architectures |

---

## 6. Recommended compilation

```bash
# Optimized binary for production (no VCS info if no git yet)
go build -ldflags="-s -w" -trimpath -o bot ./cmd
```

GopherML doesn't use Cgo, so the binary is static and cross-platform by default (`CGO_ENABLED=0`).
