# GopherML — User Guide

GopherML is a high-performance Machine Learning library for Go, designed for cryptocurrency bots and low-latency inference systems.

**Module:** `github.com/ateto3/gopherml`

> 📚 ** Related documentation:** [Quick Start](docs/GETTING_STARTED.md) · [API Reference](docs/API_REFERENCE.md) · [Performance](docs/PERFORMANCE.md) · [FAQ](docs/FAQ.md) · [Contribute](CONTRIBUTING.md)

---

## 1. Installation

### As a local dependency (development)

Add a `replace` in your bot's `go.mod`:

```go
module github.com/your-user/your-bot

go 1.23

require github.com/ateto3/gopherml v0.0.0

replace github.com/ateto3/gopherml => /path/to/GopherML
```

Then run `go mod tidy`.

### As a module published on GitHub

```bash
go get github.com/ateto3/gopherml@latest
```

---

## 2. Key Concepts

**Data Types**
- **X (Features):** `[][]float64` — shape array `[nSamples, nFeatures]`
- **y (Targets):** `[]float64` — form vector `[nSamples]`
- **OHLCV:** `crypto.OHLCV{Open, High, Low, Close, Volume float64}`

### Model Standard API
```
NewModel()          → Initialize with Default Parameters
Fit(X, y)          → Train (validates NaN/Inf automatically)
Predict(X)         → Batch prediction
PredictSingle(x)   → Single sample prediction (hot path, zero alloc extra)
PredictProba(X)    → Probabilities per class (for confidence thresholds)
Score(X, y)        → Evaluate (R² for regression, Accuracy for classification)
```

> **Important:** All `Fit` and `Predict` methods automatically reject `NaN` or `Inf` values. Never ignore returned errors.

---

## 3. `crypto/` package — Pipeline for Crypto Bots

This is the highest level package. Encapsulates the entire pipeline: → OHLCV features trading → signal.

### Quick use (bot in production)

```go
import "github.com/ateto3/gopherml/crypto"

// 1. Create predictor with settings
cfg := crypto.DefaultPredictorConfig()
cfg.NEstimators   = 200    // More trees = more stable
cfg.MaxDepth      = 8      // Avoid overfitting
cfg.WindowSize    = 20     // Last 20 candles per prediction
cfg.MinConfidence = 0.65   // Only trade if >65% confident

p := crypto.NewCryptoPredictor(cfg)

// 2. Load pre-trained model at bot startup
if err := p.LoadModel("btcusdt_1h_rf.gob"); err != nil {
    log.Fatalf("Error loading model: %v", err)
}

// 3. In the trading loop — for each new candle:
signal, confidence, err := p.PredictSignal(lastCandles) // len == WindowSize
if err != nil {
    log.Printf("Prediction failed: %v", err)
    return
}

switch signal {
case crypto.SignalLong:
fmt.Printf("LONG with %.1f%% confidence\n", confidence*100)
    // executeOrder(BUY)
case crypto.SignalShort:
fmt.Printf("SHORT with %.1f%% confidence\n", confidence*100)
    // executeOrder(SELL)
case crypto.SignalHold:
    // Don't trade
}
```

### Available signals
| Constant | Value | Meaning |
|:---|:---|:---|
| `crypto.SignalHold` | 0 | No clear signal — don't trade |
| `crypto.SignalLong` | 1 | Model predicts price rise |
| `crypto.SignalShort` | 2 | Model predicts price drop |

### Training and saving a model

```go
import "github.com/ateto3/gopherml/crypto"

cfg := crypto.DefaultPredictorConfig()
cfg.WindowSize = 20

p := crypto.NewCryptoPredictor(cfg)

// candles: historical OHLCV time series
// labels: label per window (0=hold, 1=long, 2=short)
// len(labels) must be len(candles) - WindowSize + 1
if err := p.Fit(candles, labels); err != nil {
    log.Fatal(err)
}

if err := p.SaveModel("btcusdt_1h_rf.gob"); err != nil {
    log.Fatal(err)
}
```

### Extract features manually

```go
cfg := crypto.NewFeatureConfig()
// cfg.Returns = true, cfg.LogReturns = true, cfg.VWAP = true, ...

// Feature vector for a single window
fv, err := crypto.ToFeatureVector(windowCandles, cfg)

// Complete matrix by sliding window (for training)
X, err := crypto.BuildFeatureMatrix(allCandles, windowSize, cfg)
```

**Features included by default:**
| Feature | Description |
|:---|:---|
| `Returns` | Relative returns between consecutive candles |
| `LogReturns` | Log-returns — more stationary |
| `PriceRange` | `(high-low)/close` per candle — normalized volatility |
| `VolumeRatio` | `volume/avg_volume` — volume relative to the average |
| `Momentum` | Fraction of bullish candles in the window |
| `VWAP` | Closure deviation from window VWAP |

---

## 4. Individual Models

### Random Forest — Confidence Threshold Classifier

```go
import "github.com/ateto3/gopherml/ensemble"

rf := ensemble.NewRandomForestClassifier()
rf.NEstimators = 100
rf.MaxDepth    = 10
rf.RandomState = 42

if err := rf.Fit(X, y); err != nil {
    log.Fatal(err)
}

// Batch prediction
predictions, _ := rf.Predict(XTest)

// Probabilities per class (to filter by confidence)
proba, _ := rf.PredictProba(XTest)
// test[i][1] = probability of class 1 for sample i

// Single sample prediction (hot path in bots)
label, _ := rf.PredictSingle(featureVector)
probaSingle, _ := rf.PredictProbaSingle(featureVector)
```

### Random Forest — Regressor (Price Prediction)

```go
rf := ensemble.NewRandomForestRegressor()
rf.NEstimators = 100

rf.Fit(X, y) // y = future price

// Single sample prediction (no batch overhead)
predictedPrice, _ := rf.PredictSingle(featureVector)
```

### Linear Regression

```go
import "github.com/ateto3/gopherml/linear_model"

model := linear_model.NewLinearRegression()
model.Fit(X, y)
preds, _ := model.Predict(XTest)
```

---

## 5. Preprocessing

### StandardScaler (recommended for linear models)

```go
import "github.com/ateto3/gopherml/preprocessing"

scaler := preprocessing.NewStandardScaler()

// Option A: Fit + Transform separately (train/test split)
scaler.Fit(XTrain)          // Fit only on train
XTrainScaled, _ := scaler.Transform(XTrain)
XTestScaled, _  := scaler.Transform(XTest)  // Same parameters

// Option B: FitTransform (more convenient, same result)
XScaled, _ := scaler.FitTransform(XTrain)

// Reverse the transformation
XOriginal, _ := scaler.InverseTransform(XScaled)
```

> **Attention:** Calling `Transform` before `Fit` returns an explicit error — it never fails silently.

### MinMaxScaler

```go
scaler := preprocessing.NewMinMaxScaler()
XScaled, _ := scaler.FitTransform(X) // Scales to [0, 1]
```

---

## 6. Model Persistence

```go
import "github.com/ateto3/gopherml/persistence"
import "github.com/ateto3/gopherml/ensemble"

// Save
rf := ensemble.NewRandomForestClassifier()
rf.Fit(X, y)
persistence.Save("model.gob", rf)

// Load (with automatic integrity validation)
rfLoaded := ensemble.NewRandomForestClassifier()
if err := persistence.Load("model.gob", rfLoaded); err != nil {
// Error includes detection of corrupt files or untrained models
    log.Fatal(err)
}
```

> Models implementing the `persistence.Validator` interface are automatically validated after loading. `RandomForestClassifier`, `RandomForestRegressor`, etc. already implement this interface.

---

## 7. Data validation

```go
import "github.com/ateto3/gopherml/utils"

// Validate 2D matrix and target vector (includes NaN/Inf check)
if err := utils.CheckXy(X, y); err != nil {
    log.Fatal(err)
}

// Validate only a matrix
if err := utils.CheckFinite(X); err != nil {
    log.Fatal(err)
}

// Validate an individual feature vector (hot path)
if err := utils.CheckFiniteSingle(featureVector); err != nil {
    log.Fatal(err)
}
```

---

## 8. Model Selection

```go
import (
    "fmt"
    "log"

    "github.com/ateto3/gopherml/linear_model"
    "github.com/ateto3/gopherml/model_selection"
)

// K-Fold cross-validation (the API is a KFold struct, not an integer)
kf := model_selection.KFold{NSplits: 5, Shuffle: true, RandomState: 42}
scores, err := model_selection.CrossValScore(model, X, y, kf)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Scores per fold:", scores)

// parallel grid search: a model factory + list of combinations
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
    CV:    model_selection.KFold{NSplits: 3, Shuffle: true, RandomState: 42},
    NJobs: -1, // -1 = all cores
}
if err := gs.Fit(X, y); err != nil {
    log.Fatal(err)
}
fmt.Printf("Best parameters: %v (score: %.4f)\n", gs.BestParams, gs.BestScore)
```

> **Note:** `CrossValScore` and `GridSearchCV` require models with `Fit(X, y)` and `Score(X, y)` (`Estimator` interface). Individual `tree/` trees use `Fit(X, y, sampleWeight)`, so use `ensemble/` forests.

---

## 9. Best Practices for Bots in Production

| Practice | Reason |
|:---|:---|
| Use `PredictSingle` on high-frequency loops | Avoid unnecessary allocations per tick |
| Load the model at startup, not on every prediction | `Load` is expensive — do it once |
| Use `MinConfidence >= 0.60` | Reduce false signals in sideways markets |
| Validate candles before predicting | Exchanges can send candles with gaps or volume 0 |
| Save the scaler next to the model | The `CryptoPredictor` does it automatically via `SaveModel` |
| Use `go test ./... -race` before deploying | Detects race conditions in concurrent prediction |

---

*GopherML — Go Native Machine Learning for Cryptocurrency Bots.*
