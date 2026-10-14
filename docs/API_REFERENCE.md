# 📖 GopherML API Reference

Reference of the exported types and functions of each package, generated from the source code (Go 1.23, module `github.com/ateto3/gopherml`).

> **Standard model API convention:**
>
> ```
> NewModel()          → initialize with default parameters
> Fit(X, y)           → train (validates NaN/Inf and dimensions)
> Predict(X)          → batch prediction
> PredictSingle(x)    → single sample prediction (hot path)
> PredictProba(X)     → class probabilities
> Score(X, y)         → evaluate (R² for regression, accuracy for classification)
> ```
>
> `X` is `[][]float64` (`[nSamples][nFeatures]`), `y` is `[]float64`. Never ignore errors.

---

## `utils` — Data validation and utilities

All data input must pass through here (models do this internally).

| Function | Signature | Description |
| :--- | :--- | :--- |
| `CheckArray` | `(X [][]float64) error` | Validates that the matrix is not empty and has homogeneous rows |
| `CheckFinite` | `(X [][]float64) error` | Rejects `NaN`/`Inf` in the entire matrix |
| `CheckFiniteSingle` | `(x []float64) error` | Same for a vector (hot path, faster) |
| `CheckFiniteVec` | `(y []float64) error` | Same for the target vector |
| `CheckXy` | `(X [][]float64, y []float64) error` | Validates matrix + vector: dimensions and finite |
| `CheckIsFitted` | `(fitted bool, modelName string) error` | Returns error if the model is not trained |
| `TrainTestSplit` | `(X, y, testSize float64, randomState int64) (XTrain, XTest, yTrain, yTest [][]float64/[]float64, error)` | Train/test split; `testSize ∈ (0,1)` |
| `LoadCSV` | `(path string, hasHeader bool) (X [][]float64, y []float64, error)` | Loads CSV: last column = target |
| `NewRNG` | `(seed int64) *rand.Rand` | Deterministic RNG for reproducible tests |
| `Shuffle` | `(n int, rng *rand.Rand) []int` | Deterministic permutation |

---

## `preprocessing` — Scaling and encoding

They all follow the `Fit` → `Transform` (or `FitTransform`) pattern. `Transform` before `Fit` returns error.

### `StandardScaler`

Standardize by subtracting the mean and dividing by the standard deviation (equivalent to `sklearn.preprocessing.StandardScaler`).

```go
s := preprocessing.NewStandardScaler()   // *StandardScaler{Mean, Std, IsFitted}
s.Fit(X) error
s.Transform(X) ([][]float64, error)
s.FitTransform(X) ([][]float64, error)
s.InverseTransform(X) ([][]float64, error)
```

### `MinMaxScaler`

Scales each feature to the range `[0, 1]` (`Min`, `Max`, `IsFitted`). Same methods as `StandardScaler`.

### `LabelEncoder`

Encodes text labels to integers `0..n-1`:

```go
le := preprocessing.NewLabelEncoder()    // *LabelEncoder{Classes []string, Mapping map[string]int}
le.Fit(y []string)
le.Transform(y []string) ([]float64, error)
le.InverseTransform(y []float64) ([]string, error)
```

### `PolynomialFeatures`

Generates polynomial and interaction features **up to degree 2** (higher degree returns error):

```go
pf := preprocessing.NewPolynomialFeatures(2)  // *PolynomialFeatures{Degree, IncludeBias}
pf.Fit(X) error
pf.FitTransform(X) ([][]float64, error)
pf.Transform(X) ([][]float64, error)
```

---

## `metrics` — Evaluation

All functions compare `yTrue` vs `yPred` (both `[]float64`) and return `(float64, error)` unless otherwise indicated.

### CLASSIFICATION

| Function | Signature | Notes |
| :--- | :--- | :--- |
| `AccuracyScore` | `(yTrue, yPred []float64) (float64, error)` | Fraction of correct predictions |
| `PrecisionScore` | `(yTrue, yPred []float64, average string) (float64, error)` | `average`: `"macro"`, `"weighted"` or `""` (binary) |
| `RecallScore` | `(yTrue, yPred []float64, average string) (float64, error)` | Same |
| `F1Score` | `(yTrue, yPred []float64, average string) (float64, error)` | Same |
| `ROCAUCScore` | `(yTrue, yScores []float64) (float64, error)` | `yScores` = positive class probabilities |
| `ConfusionMatrix` | `(yTrue, yPred []float64) ([][]int, []float64, error)` | Confusion matrix + classes |
| `ClassificationReport` | `(yTrue, yPred []float64) (string, error)` | Per-class sklearn-style report |

### REGRESSION

| Function | Signature |
| :--- | :--- |
| `MeanSquaredError` | `(yTrue, yPred []float64) (float64, error)` |
| `RootMeanSquaredError` | `(yTrue, yPred []float64) (float64, error)` |
| `MeanAbsoluteError` | `(yTrue, yPred []float64) (float64, error)` |
| `R2Score` | `(yTrue, yPred []float64) (float64, error)` |

> Models also include their own `Score(X, y)`: **accuracy** for classifiers, **R²** for regressors.

---

## `linear_model` — Linear regression and classification

### Sentinel Errors

```go
linear_model.ErrEmptyData         // empty input
linear_model.ErrDimensionMismatch // incompatible dimensions
linear_model.ErrNotFitted         // Predict before Fit
```

Models

| Model | Constructor | Methods | Key Hyperparameters |
| :--- | :--- | :--- | :--- |
| `LinearRegression` (OLS) | `NewLinearRegression()`, `NewLinearRegressionWithParams(LinearRegressionParams)` | `Fit`, `Predict`, `Score` | `FitIntercept`, `CopyX` |
| `Ridge` (L2) | `NewRidge()`, `NewRidgeWithParams(RidgeParams)` | `Fit`, `Predict`, `Score` | `Alpha`, `Solver`, `MaxIter`, `Tol` |
| `Lasso` (L1) | `NewLasso()`, `NewLassoWithParams(LassoParams)` | `Fit`, `Predict`, `Score` | `Alpha`, `MaxIter`, `Tol`, `Positive` |
| `ElasticNet` (L1+L2) | `NewElasticNet()`, `NewElasticNetWithParams(ElasticNetParams)` | `Fit`, `Predict`, `Score` | `Alpha`, `L1Ratio`, `MaxIter`, `Tol` |
| `LogisticRegression` | `NewLogisticRegression()` | `Fit`, `Predict`, `PredictProba`, `Score` | `LearningRate`, `MaxIter`, `Tol` (binary) |
| `SGDClassifier` | `NewSGDClassifier()` | `Fit`, `Predict` | `LearningRate`, `MaxIter`, `Tol`, `RandomState` |

`LinearRegression`, `Ridge`, `Lasso`, and `ElasticNet` embed `LinearModel`:

```go
type LinearModel struct {
    Coef         []float64 // coefficients
    Intercept    float64   // bias
    FitIntercept bool
    NFeaturesIn  int
    IsFitted     bool
}
// methods: Predict(X) ([]float64, error), Score(X, y) (float64, error)
```

### Parallel versions and utilities

| Symbol | Signature / description |
| :--- | :--- |
| `ParallelOptions` | `{NumWorkers, MinBatchSize int}` — `DefaultParallelOptions()` uses all cores |
| `ParallelFor` | `(n int, opts ParallelOptions, fn func(start, end int))` — parallelizes a range |
| `ParallelMatrixMultiply` | `(X [][]float64, opts) [][]float64` |
| `ParallelMatVecMul` | `(X [][]float64, v []float64, opts) []float64` |
| `ParallelCenterData` | centers data in parallel, returns offsets |
| `LinearRegressionParallel` | `NewLinearRegressionParallel()` — parallel OLS |
| `LassoParallel` | `NewLassoParallel()`, `NewLassoParallelWithParams(LassoParams, ParallelOptions)` |
| `WorkerPool` | `NewWorkerPool(workers int)` reusable for batches |
| `DotProduct`, `Mean`, `Transpose`, `MeanSquaredError` | mathematical utilities |

Guide to when to parallelize in [PERFORMANCE.md](PERFORMANCE.md) and [../linear_model/PERFORMANCE.md](../linear_model/PERFORMANCE.md).

---

## `tree` — CART decision trees

> ⚠️ **Attention:** `Fit` of the trees receives a third argument `sampleWeight` (use `nil` if there are no weights). That's why they **don't** directly implement the `Estimator` interface of `model_selection`; use the `ensemble` forests for cross-validation.

### `DecisionTreeClassifier`

```go
dt := tree.NewDecisionTreeClassifier()
// or with explicit parameters:
dt := tree.NewDecisionTreeClassifierWithParams(tree.DecisionTreeClassifierParams{ /* ... */ })

dt.Fit(X [][]float64, y []float64, sampleWeight []float64) error
dt.Predict(X) ([]int, error)
dt.PredictProba(X) ([][]float64, error)
dt.Score(X, y) (float64, error)
dt.GetDepth() int
dt.GetNLeaves() int
```

### `DecisionTreeRegressor`

Same methods, but `Predict` returns `([]float64, error)` and the default criterion is `MSE`.

### Common parameters (struct fields)

| Field | Type | Description |
| :--- | :--- | :--- |
| `Criterion` | `string` | `"gini"`, `"entropy"` (classif.) or `"mse"` (regress.) |
| `Splitter` | `string` | Split search strategy |
| `MaxDepth` | `int` | Maximum depth (`0` = unlimited) |
| `MinSamplesSplit` | `int` | Minimum samples to split a node |
| `MinSamplesLeaf` | `int` | Minimum samples per leaf |
| `MaxFeatures` | `interface{}` | `int`, `float64` or `"sqrt"`/`"log2"` |
| `RandomState` | `int` | Reproducibility seed |
| `CCPAlpha` | `float64` | Cost-complexity pruning |

Fitted attributes after `Fit`: `Tree`, `NFeatures`, `NOutputs`, `NClasses`, `Classes`, `IsFitted`.

### Extensible components

`Criterion` (`Gini`, `Entropy`, `MSE`), `Splitter` (`BestSplitter`, `RandomSplitter`), `TreeBuilder` (`DepthFirstTreeBuilder`) and `BuildTree(...)` to build low-level trees. `PrintTree(t, featureNames)` prints the tree.

---

## `ensemble` — Random Forest

### `RandomForestClassifier`

```go
rf := ensemble.NewRandomForestClassifier()

rf.Fit(X [][]float64, y []float64) error
rf.Predict(X) ([]int, error)
rf.PredictProba(X) ([][]float64, error)
rf.PredictSingle(x []float64) (int, error)          // hot path
rf.PredictProbaSingle(x []float64) ([]float64, error) // hot path
rf.Score(X, y) (float64, error)                       // accuracy
rf.ValidateLoaded() error                             // called by persistence.Load
```

### `RandomForestRegressor`

```go
rf := ensemble.NewRandomForestRegressor()

rf.Fit(X, y) error
rf.Predict(X) ([]float64, error)
rf.PredictSingle(x) (float64, error)
rf.Score(X, y) (float64, error)  // R²
rf.ValidateLoaded() error
```

### Hyperparameters (embedded from `BaseForest`)

| Field | Type | Description |
| :--- | :--- | :--- |
| `NEstimators` | `int` | Number of trees (trained in parallel) |
| `MaxDepth` | `int` | Maximum depth per tree |
| `MinSamplesSplit` / `MinSamplesLeaf` | `int` | Node/Leaf size restrictions |
| `MaxFeatures` | `interface{}` | Features considered by split |
| `Bootstrap` | `bool` | Sampling with replacement |
| `RandomState` | `int64` | Seed (reproducibility) |
| `IsFitted` | `bool` | Status after training |

Attributes: `Classes []int`, `NClasses int`, `Estimators []…`.

Utilities: `ForestParams` / `DefaultForestParams()`, `GenerateBootstrapIndices(n, rng)`, `GetBootstrapSample(X, y, rng)`.

---

## `model_selection` — Validation and search

### `KFold`

```go
kf := model_selection.KFold{NSplits: 5, Shuffle: true, RandomState: 42}
splits := kf.Split(nSamples) // []SplitRecord{TrainIndices, TestIndices}
```

### `CrossValScore`

```go
func CrossValScore(model Estimator, X [][]float64, y []float64, cv KFold) ([]float64, error)
```

`Estimator` is any model with `Fit(X, y) error` and `Score(X, y) (float64, error)` — e.g. Random Forest or `LogisticRegression` (trees from `tree/`, with `Fit` of 3 arguments, do not meet this interface).

### `GridSearchCV`

Comprehensive **parallel** search on an explicit list of combinations:

```go
gs := &model_selection.GridSearchCV{
    EstimatorFactory: func(params map[string]interface{}) model_selection.Estimator {
        lr := linear_model.NewLogisticRegression()
        lr.LearningRate = params["lr"].(float64)
        lr.MaxIter = params["iter"].(int)
        return lr
    },
    ParamGrid: []map[string]interface{}{
        {"lr": 0.01, "iter": 100},
        {"lr": 0.10, "iter": 200},
    },
    CV:   model_selection.KFold{NSplits: 3, Shuffle: true, RandomState: 42},
    NJobs: -1, // -1 = all cores
}

err := gs.Fit(X, y)
// gs.BestParams, gs.BestScore, gs.CVResults ([]CVResult{Params, MeanScore, Err})
```

---

## `persistence` — Save and Load Models

| Function | Signature | Description |
| :--- | :--- | :--- |
| `Save` | `(filename string, model interface{}) error` | Serialize to `gob` |
| `Load` | `(filename string, model interface{}) error` | Deserialize + validate (see `Validator`) |
| `Dump` | `(model interface{}, w io.Writer) error` | Serialize to any `io.Writer` |
| `LoadFrom` | `(r io.Reader, model interface{}) error` | Deserialize from any `io.Reader` |

```go
type Validator interface {
    ValidateLoaded() error
}
```

Models that implement `Validator` (Random Forest, etc.) are automatically validated after `Load`: it detects corrupt files and untrained models.

---

## `crypto` — OHLCV → signal pipeline

### Data Types

```go
type OHLCV struct { Open, High, Low, Close, Volume float64 }

type Signal int
const (
    SignalHold  Signal = 0
    SignalLong  Signal = 1
    SignalShort Signal = 2
)
```

### Features (`FeatureConfig`)

Boolean fields (all `true` with `NewFeatureConfig()`): `Returns`, `LogReturns`, `PriceRange`, `VolumeRatio`, `Momentum`, `VWAP`.

| Role | Signature |
| :--- | :--- |
| `ToFeatureVector` | `(candles []OHLCV, cfg FeatureConfig) ([]float64, error)` — a window |
| `BuildFeatureMatrix` | `(candles []OHLCV, windowSize int, cfg FeatureConfig) ([][]float64, error)` — sliding window |
| `FeatureSize` | `(windowSize int, cfg FeatureConfig) int` — # of resulting columns |
| `MinWindowSize` | `() int` — minimum window allowed |

### `CryptoPredictor`

```go
type PredictorConfig struct {
    NEstimators    int           // forest trees (default 100)
    MaxDepth       int           // depth (default 10)
    WindowSize     int           // candles per window (default 20)
    MinConfidence  float64       // threshold for non-HOLD signal (default 0.55)
    Features       FeatureConfig // feature groups
}
cfg := crypto.DefaultPredictorConfig()

p := crypto.NewCryptoPredictor(cfg)
p.Fit(candles []OHLCV, labels []float64) error  // len(labels) == len(candles)-WindowSize+1
p.PredictSignal(candles []OHLCV) (Signal, float64, error) // len(candles) == WindowSize
p.SaveModel(path string) error
p.LoadModel(path string) error                  // model + scaler together
p.Config() PredictorConfig
p.IsFitted() bool
```

Complete user guide in [../USAGE.md](../USAGE.md).
