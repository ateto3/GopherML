# 🧪 GopherML Examples

Executable programs (`package main`) that demonstrate full usage of the library. All of them run with data included in the code itself: no external downloads needed.

**Requirements:** Go 1.23+ and the module `github.com/ateto3/gopherml` available (see [README](../README.md#2-installation)).

---

## Available examples

### 1. Basic classification

**File:** [`classification/basic_classification.go`](classification/basic_classification.go)

Random Forest for binary classification: scaling with `StandardScaler`, training, *batch* and single-sample prediction (`PredictSingle`), probabilities (`PredictProba`), cross-validation with `KFold` and detailed metrics (precision, recall, F1).

```bash
go run examples/classification/basic_classification.go
```

### 2. Basic regression

**File:** [`regression/basic_regression.go`](regression/basic_regression.go)

Comparison of `LinearRegression` and `RandomForestRegressor` on a price prediction problem: scaling, predictions, R², MSE/MAE and model comparison.

```bash
go run examples/regression/basic_regression.go
```

### 3. Crypto trading pipeline

**File:** [`crypto/crypto_trading_bot.go`](crypto/crypto_trading_bot.go)

Full pipeline with the `crypto/` package: generates synthetic OHLCV candles, builds labels with a sliding window, trains a `CryptoPredictor`, saves and reloads the model (`SaveModel`/`LoadModel`) and emits a `LONG`/`SHORT`/`HOLD` signal with its confidence level.

```bash
go run examples/crypto/crypto_trading_bot.go
```

### 4. Model persistence

**File:** [`persistence/model_persistence.go`](persistence/model_persistence.go)

Trains a Random Forest, saves it to `gob`, reloads it into a new model and verifies that predictions match. It also demonstrates corrupt-file detection through the `Validator` interface.

```bash
go run examples/persistence/model_persistence.go
```

### 5. Hyperparameter search

**File:** [`grid_search/hyperparameter_tuning.go`](grid_search/hyperparameter_tuning.go)

Parallel `GridSearchCV` (`NJobs: -1`) over `LogisticRegression`: model factory, list of combinations, scores per combination and selection of the best parameters.

```bash
go run examples/grid_search/hyperparameter_tuning.go
```

### 6. Preprocessing pipeline

**File:** [`preprocessing/preprocessing_pipeline.go`](preprocessing/preprocessing_pipeline.go)

Use of every transformer: `StandardScaler` (including `InverseTransform`), `MinMaxScaler`, `LabelEncoder` (text → integers) and `PolynomialFeatures`, plus error handling when calling `Transform` before `Fit`.

```bash
go run examples/preprocessing/preprocessing_pipeline.go
```

---

## Run them all

```bash
for f in examples/*/*.go; do
  echo "▶ $f"
  go run "$f" > /dev/null || exit 1
done
echo "✔ All examples run correctly"
```

## Notes

- Each directory is a standalone package: you can run a single file with `go run <path>` without affecting the rest.
- The `.gob` files created during execution (persistence/crypto) are deleted automatically at the end and are excluded from the repository via `.gitignore`.
- For performance and *benchmarking* examples, see [docs/PERFORMANCE.md](../docs/PERFORMANCE.md).
