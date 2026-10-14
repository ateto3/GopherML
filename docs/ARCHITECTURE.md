# 🏗️ GopherML Architecture

This document describes the project architecture, design patterns, and technical decisions. For a more narrative explanation of the algorithms, see [EXPLAINING.md](../EXPLAINING.md).

---

## 1. Philosophy

1. **scikit-learn mirror:** packages, method names, and flows identical to `sklearn` for zero learning curve if you come from Python.
2. **Zero dependencies:** exclusively Go standard library. No Cgo, no NumPy, no bindings.
3. **Data efficiency:** flat arrays, pre-allocated memory, and direct access to improve cache locality.
4. **Concurrency where it pays:** goroutines only in operations with enough work to justify synchronization cost.
5. **Fail fast and loud:** all invalid input (`NaN`, `Inf`, inconsistent dimensions, untrained model) produces an explicit `error` — never silently wrong results.

---

## 2. System Layers

```
                ┌─────────────────────────────────────────────┐
   High layer   │  crypto/   (OHLCV → features → signal)      │  ← orchestration
                ├─────────────────────────────────────────────┤
                │  model_selection/   (KFold, CV, GridSearch) │
                │  persistence/       (Save/Load + Validator) │
                ├─────────────────────────────────────────────┤
   Models       │  ensemble/ (RandomForest) ──► tree/ (CART)   │
                │  linear_model/ (OLS, Ridge, Lasso, Logistic,│
                │                  ElasticNet, SGD)           │
                ├─────────────────────────────────────────────┤
   Support      │  preprocessing/ (scalers, encoders)         │
                │  metrics/ (accuracy, F1, R², ROC-AUC...)    │
                ├─────────────────────────────────────────────┤
   Base         │  utils/ (CheckFinite, TrainTestSplit, RNG)  │  ← everything passes through here
                └─────────────────────────────────────────────┘
```

Dependency rule: **a layer can only use those below it**. `utils` imports from no one; `crypto` can use everything else.

---

## 3. Training Data Flow

```
Raw data ([][]float64)
      │
      ▼
utils.CheckXy / CheckFinite      ← validates NaN/Inf and dimensions
      │
      ▼
preprocessing (optional)         ← Fit on train, Transform on train and test
      │
      ▼
model.Fit(X, y)                  ← each model validates internally again
      │                              (CheckFinite + CheckIsFitted)
      ▼
model.Predict(X) / PredictSingle ← batch or single sample inference
      │
      ▼
metrics.* (yTrue, yPred)         ← evaluation
      │
      ▼
persistence.Save("model.gob")    ← gob + integrity validation on load
```

---

## 4. Design Patterns

### 4.1 Struct embedding (Python-like inheritance)

```go
type LinearModel struct { Coef []float64; Intercept float64 /* ... */ }

type Lasso struct {
    LinearModel          // ← embeds base model
    Alpha float64
    // ...
}
```

`Ridge`, `Lasso`, `ElasticNet`, and `LinearRegression` embed `LinearModel` and reuse its linear system solving logic (Gaussian Elimination). Random Forest hyperparameters live in `ensemble.BaseForest`.

**Convention:** coefficients are accessed directly (`model.Coef`), without getters.

### 4.2 Small interfaces

| Interface | Package | Purpose |
| :--- | :--- | :--- |
| `Estimator` (`Fit` + `Score`) | `model_selection` | Any model evaluable by `CrossValScore`/`GridSearchCV` |
| `Validator` (`ValidateLoaded`) | `persistence` | Auto integrity validation after `Load` |
| `Criterion`, `Splitter`, `TreeBuilder` | `tree` | Interchangeable split strategies (Gini/Entropy/MSE, best split/random) |

### 4.3 Constructors with safe defaults

Always `New...()`:

```go
rf := ensemble.NewRandomForestClassifier() // NEstimators, MaxDepth, RandomState initialized
```

An empty struct literal (`&ensemble.RandomForestClassifier{}`) may train with zero hyperparameters.

### 4.4 Sentinel errors

Packages like `tree` and `linear_model` export errors comparable with `errors.Is`:

```go
var ErrNotFitted         = errors.New("tree is not fitted yet")
var ErrDimensionMismatch = errors.New("dimension mismatch")
var ErrEmptyData         = errors.New("empty input data")
```

---

## 5. Intelligent parallelism

GopherML doesn't parallelize for the sake of parallelizing:

| Operation | Parallel? | Mechanism |
| :--- | :--- | :--- |
| Random Forest training | ✅ | One goroutine per tree + `sync.WaitGroup` |
| `GridSearchCV` | ✅ | One goroutine per hyperparameter combination (`NJobs`) |
| Linear algebra (Gram matrix, `MatVec`) | ✅ | Row blocks via `ParallelFor` / `WorkerPool` |
| Single sample prediction (`PredictSingle`) | ❌ | Goroutine overhead exceeds useful work |
| Small dataset scaling | ❌ | Sequential is faster |

Details and benchmarks in [PERFORMANCE.md](PERFORMANCE.md) and [../linear_model/PERFORMANCE.md](../linear_model/PERFORMANCE.md).

---

## 6. `tree` package design (base of `ensemble`)

- **Indexed nodes:** the tree stores nodes in a flat `[]Node` with integer indices instead of pointers between structs → less garbage collector pressure and better cache locality.
- **Iterative DFS constructor:** `DepthFirstTreeBuilder` uses an explicit stack (`StackRecord`) instead of recursion → impossible to suffer stack overflow with large `MaxDepth`.
- **Splitter with interchangeable criterion:** `BestSplitter` (exhaustive search) or `RandomSplitter` (used by Random Forest for per-node randomness), combined with `Gini`, `Entropy`, or `MSE`.

---

## 7. Persistence

- **`gob` format** from the standard library: fast, Go-native, no dependencies.
- `persistence.Save` serializes the complete model; `Load` deserializes it into a model you provide (the type must match).
- If the model implements `Validator`, `Load` executes `ValidateLoaded()` automatically and fails if the model is incomplete (corrupt file, incomplete structure).
- `CryptoPredictor.SaveModel` saves model **and** scaler together internally, avoiding the classic mismatch error between train and inference.

---

## 8. Code conventions

- Doc comments in **English** starting with the symbol name (`// RandomForestClassifier implements...`).
- Package names in `snake_case` sklearn style (`linear_model`, `model_selection`).
- Tests live in the same package (`*_test.go`) and are deterministic (`utils.NewRNG(seed)`).
- No global state: models are self-contained structs; predictions don't mutate state, so they are parallel-safe **after** `Fit`.

See [CONTRIBUTING.md](../CONTRIBUTING.md) for the complete contribution workflow.
