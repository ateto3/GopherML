# 🐹 GopherML 🚀

> **Native Go machine learning: low-latency inference, zero dependencies, scikit-learn-style API.**

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/ateto3/gopherml.svg)](https://pkg.go.dev/github.com/ateto3/gopherml)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen)](#5-tests-y-ci)
[![CI](https://github.com/ateto3/gopherml/actions/workflows/ci.yml/badge.svg)](https://github.com/ateto3/gopherml/actions/workflows/ci.yml)

**GopherML** is a machine learning library written entirely in **Go**, designed to provide ultra-fast training and inference tools using native concurrency (goroutines) and **without Cgo or external dependencies**. It follows Python’s *scikit-learn* architecture and API, adapted for cryptocurrency bots and high-frequency inference.

---

## 📑 Table of contents

1. [Features](#1-features)
2. [Installation](#2-installation)
3. [Getting started](#3-getting-started)
4. [Packages](#4-packages)
5. [Tests and CI](#5-tests-and-ci)
6. [Repository structure](#6-repository-structure)
7. [Documentation](#7-documentation)
8. [Examples](#8-examples)
9. [Subproject: go-nlp](#9-subproject-go-nlp)
10. [Roadmap](#10-roadmap)
11. [Design metrics](#11-design-metrics)
12. [Contributing and license](#12-contributing-and-license)

---

## 1. Features

- **⚡ Low-latency inference:** `PredictSingle` and `PredictProbaSingle` are optimized for bot hot paths (zero extra allocations per tick).
- **🧩 scikit-learn-style API:** `NewModel()` → `Fit(X, y)` → `Predict(X)` → `Score(X, y)`. If you already use scikit-learn, you'll feel at home.
- **🌊 Native concurrency:** parallel training of Random Forest, `GridSearchCV` and linear algebra with `goroutines` and `sync.WaitGroup` — without Cgo.
- **🛡️ Strict validation:** All input data is validated against `NaN`/`Inf` (`utils.CheckFinite`) before touching a model. It never fails in silence.
- **💾 Persistence with integrity:** `Save`/`Load` in `gob` format, with automatic post-load validation via the `persistence.Validator` interface.
- **📈 Complete metrics:** Accuracy, Precision, Recall, F1, R², MAE, MSE, RMSE, and ROC-AUC.
- **🪙 Crypto pipeline:** the `crypto/` package encapsulates the complete flow **OHLCV → features → trading signal** (`Hold`/`Long`/`Short`).
- **🚫 Zero dependencies:** only the Go standard library. No NumPy, no CGO, no bindings.

---

## 2. Installation

**Requirements:** Go 1.23 or higher.

### As a module published on GitHub

```bash
go get github.com/ateto3/gopherml@latest
```

### As a local dependency (development)

Add a `replace` in your project's `go.mod`:

```go
module github.com/your-user/your-bot

go 1.23

require github.com/ateto3/gopherml v0.0.0

replace github.com/ateto3/gopherml => /path/to/GopherML
```

Then run:

```bash
go mod tidy
```

---

## 3. Getting started

Complete binary classification with Random Forest:

```go
package main

import (
	"fmt"
	"log"

	"github.com/ateto3/gopherml/ensemble"
	"github.com/ateto3/gopherml/preprocessing"
)

func main() {
	// Data: [length, width] → label 0 or 1
	X := [][]float64{
		{1.0, 2.0}, {1.5, 1.8}, {1.2, 2.5}, {1.1, 1.9},
		{5.0, 6.0}, {5.5, 5.8}, {5.2, 6.5}, {4.9, 5.9},
	}
	y := []float64{0, 0, 0, 0, 1, 1, 1, 1}

	// 1. Scale the data (essential for linear models)
	scaler := preprocessing.NewStandardScaler()
	Xs, err := scaler.FitTransform(X)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Train
	rf := ensemble.NewRandomForestClassifier()
	rf.NEstimators = 100
	rf.MaxDepth = 5
	rf.RandomState = 42
	if err := rf.Fit(Xs, y); err != nil {
		log.Fatal(err)
	}

	// 3. Predict (batch)
	preds, err := rf.Predict(Xs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Predictions:", preds) // [0 0 0 0 1 1 1 1]

	// 4. Evaluate
	acc, err := rf.Score(Xs, y)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Accuracy: %.2f%%\n", acc*100)

	// 5. Hot path: single sample, no batch overhead
	label, err := rf.PredictSingle([]float64{5.1, 6.1})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Single sample:", label)
}
```

**Full pipeline for crypto bots** (OHLCV → signal):

```go
cfg := crypto.DefaultPredictorConfig()
cfg.WindowSize = 20
cfg.MinConfidence = 0.65

p := crypto.NewCryptoPredictor(cfg)
if err := p.LoadModel("btcusdt_1h_rf.gob"); err != nil {
	log.Fatal(err)
}

signal, confidence, err := p.PredictSignal(lastCandles) // len(lastCandles) == WindowSize
if err == nil && confidence >= cfg.MinConfidence {
	// signal: crypto.SignalHold | crypto.SignalLong | crypto.SignalShort
	_ = signal
}
```

See **[USAGE.md](USAGE.md)** for the full guide with trading examples.

---

## 4. Packages

| Package | Description |
| :--- | :--- |
| [`/crypto`](crypto/) | **Crypto pipeline:** `OHLCV` → feature engineering → trading signal (`CryptoPredictor`) |
| [`/linear_model`](linear_model/) | `LinearRegression`, `Ridge`, `Lasso`, `ElasticNet`, `LogisticRegression`, `SGDClassifier` + parallel versions |
| [`/tree`](tree/) | `DecisionTreeClassifier` and `DecisionTreeRegressor` (CART; Gini, Entropy, and MSE criteria) |
| [`/ensemble`](ensemble/) | `RandomForestClassifier` and `RandomForestRegressor` (parallel training with bootstrap) |
| [`/metrics`](metrics/) | Accuracy, Precision, Recall, F1, R², MAE, MSE, RMSE, ROC-AUC, reports, and confusion matrix |
| [`/model_selection`](model_selection/) | `KFold`, `CrossValScore`, and parallel `GridSearchCV` |
| [`/preprocessing`](preprocessing/) | `StandardScaler`, `MinMaxScaler`, `LabelEncoder`, `PolynomialFeatures` |
| [`/persistence`](persistence/) | Model `Save`/`Load` with post-load integrity validation |
| [`/utils`](utils/) | Data validation (`CheckXy`, `CheckFinite`), `TrainTestSplit`, `LoadCSV`, deterministic RNG |
| [`/examples`](examples/) | Executable end-to-end examples |

---

---

## 5. Tests and CI

```bash
# Build everything
go build ./...

# Analyze code (must be clean)
go vet ./...

# Tests with race detector
go test ./... -race

# Format (recommended before a PR)
gofmt -w .
```

Current status: **all packages in green** and clean `go vet`. Each push runs the [CI workflow](.github/workflows/ci.yml) (`go vet` + `go test -race`).

---

## 6. Repository structure

```
GopherML/
├── README.md              # This file
├── USAGE.md               # Usage guide with crypto bot examples
├── EXPLAINING.md          # Architecture and internal logic
├── CONTRIBUTING.md        # How to contribute (includes AI guidelines)
├── CHANGELOG.md           # Change history
├── docs/                  # Extended documentation (guides, API, FAQ)
├── crypto/                # OHLCV → features → signal pipeline
├── ensemble/              # Random Forest
├── linear_model/          # Linear models and SGD
├── tree/                  # CART decision trees
├── metrics/               # Evaluation metrics
├── model_selection/       # Cross-validation and grid search
├── preprocessing/         # Scaling and encoding
├── persistence/           # Save/load models (gob)
├── utils/                 # Data validation and utilities
├── examples/              # Executable examples
├── go-nlp/                # Submodule: concurrent NLP (NLTK port)
└── .github/               # Issue/PR templates and CI
```

---

## 7. Documentation

| Document | Content |
| :--- | :--- |
| **[docs/GETTING_STARTED.md](docs/GETTING_STARTED.md)** | Step-by-step tutorial from scratch |
| **[USAGE.md](USAGE.md)** | Usage guide with crypto bot examples |
| **[docs/API_REFERENCE.md](docs/API_REFERENCE.md)** | API reference for all packages |
| **[EXPLAINING.md](EXPLAINING.md)** | Architecture and internal logic |
| **[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)** | Detailed architecture and design decisions |
| **[docs/PERFORMANCE.md](docs/PERFORMANCE.md)** | Performance guide and optimizations |
| **[docs/FAQ.md](docs/FAQ.md)** | Frequently asked questions |
| **[CONTRIBUTING.md](CONTRIBUTING.md)** | Contribution guide and code standards |
| **[CHANGELOG.md](CHANGELOG.md)** | Change history |

Full index in **[docs/README.md](docs/README.md)**.

---

## 8. Examples

All examples are executable programs (`package main`):

```bash
# Classification with Random Forest + metrics + cross-validation
go run examples/classification/basic_classification.go

# Regression: LinearRegression vs RandomForestRegressor
go run examples/regression/basic_regression.go

# Crypto pipeline: synthetic OHLCV → signal → persistence
go run examples/crypto/crypto_trading_bot.go

# Save and load a trained model
go run examples/persistence/model_persistence.go

# Parallel hyperparameter search with GridSearchCV
go run examples/grid_search/hyperparameter_tuning.go

# Preprocessing pipeline with scalers and encoder
go run examples/preprocessing/preprocessing_pipeline.go
```

Details in **[examples/README.md](examples/README.md)**.

---

## 9. Subproject: go-nlp

The [`go-nlp/`](go-nlp/) directory contains a ** standalone Go module ** (`github.com/ateto3/gopherml/go-nlp`): a concurrent NLTK port with tokenization, POS tagging, parsing, text classification and sentiment analysis, plus token optimization for LLM APIs.

> ⚠️ `go-nlp` has its own `go.mod`: it is independently compiled and tested with
> `cd go-nlp && go build ./... && go test ./...`

See **[go-nlp/README.md](go-nlp/README.md)** for full documentation.

---

## 10. Roadmap

- [x] Basic ML models (linear regression, classification, trees, forests, ensembles)
- [x] Model persistence (Save/Load with integrity validation)
- [x] Evaluation metrics (accuracy, precision, recall, F1, R², MAE, RMSE)
- [x] Model selection and cross-validation (KFold, GridSearchCV)
- [x] Data preprocessing (StandardScaler, MinMaxScaler, LabelEncoder)
- [x] NaN/Inf validation on all input data
- [x] `PredictProba` + `PredictSingle` for high-frequency bots
- [x] `crypto/` package — OHLCV → features → trading signal pipeline
- [ ] GradientBoosting (XGBoost-like)
- [ ] Streaming data support (online updates)

---

## 11. Design metrics

| Objective | Specification |
| :--- | :--- |
| **Language** | Go 1.23+ |
| **Module** | `github.com/ateto3/gopherml` |
| **Target** | Low latency inference/crypto bots |
| **Hardware** | Optimized for ARM64 (Raspberry Pi 4) and x86_64 |
| **Dependencies** | Zero (Go stdlib only) |
| **Tests** | `go test ./... -race` — 100% green tests |

---

## 12. Contributing and license

Contributions are welcome! Read **[CONTRIBUTING.md](CONTRIBUTING.md)** for workflow, code standards, and performance guidelines (including guidelines for AI assistants). Before sending a PR, make sure they pass:

```bash
go vet ./... && go test ./... -race
```

This project is licensed under **mit** — see [LICENSE](LICENSE).

---

## Developer 📝 Notes (Atheist)

This project is part of an ecosystem that includes cryptocurrency trading bots. The code must be modular so that it can be easily imported into other Go binaries or automation microservices.

> "The code must be so fast that Go's garbage collector doesn't even know we're here."
