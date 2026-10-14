# 🚀 Quick Start with GopherML

Practical guide to train your first model with GopherML, from scratch to saving it to disk. All examples are complete and executable programs.

**Requirements:** Go 1.23+ · see [README.md](../README.md) for installation.

---

## 1. Concepts in 30 seconds

GopherML replicates scikit-learn's API with Go types:

| Concept | GopherML Type | Form |
| :--- | :--- | :--- |
| **X** (features) | `[][]float64` | `[nSamples][nFeatures]` |
| **y** (labels) | `[]float64` | `[nSamples]` |
| **Predictions** | `[]float64` or `[]int` | depending on model |

Model lifecycle:

```
NewModel() → Fit(X, y) → Predict(X) / PredictSingle(x) → Score(X, y)
```

> ⚠️ All `Fit` and `Predict` methods validate that **no data is `NaN` or `Inf`** and return error if dimensions don't match. Never ignore the error.

---

## 2. Your first model (classification)

```go
package main

import (
	"fmt"
	"log"

	"github.com/ateto3/gopherml/ensemble"
	"github.com/ateto3/gopherml/metrics"
	"github.com/ateto3/gopherml/utils"
)

func main() {
	X := [][]float64{
		{2.5, 1.0}, {3.0, 1.2}, {2.2, 0.9}, {2.8, 1.1},
		{7.5, 5.0}, {8.0, 4.8}, {7.2, 5.2}, {7.9, 4.9},
	}
	y := []float64{0, 0, 0, 0, 1, 1, 1, 1}

	// Train/test split 75/25 with deterministic RNG
	XTrain, XTest, yTrain, yTest, err := utils.TrainTestSplit(X, y, 0.25, 42)
	if err != nil {
		log.Fatal(err)
	}

	// Train a Random Forest
	rf := ensemble.NewRandomForestClassifier()
	rf.NEstimators = 100
	rf.MaxDepth = 4
	rf.RandomState = 42
	if err := rf.Fit(XTrain, yTrain); err != nil {
		log.Fatal(err)
	}

	// Predict and convert to []float64 for metrics
	preds, err := rf.Predict(XTest)
	if err != nil {
		log.Fatal(err)
	}
	yPred := make([]float64, len(preds))
	for i, p := range preds {
		yPred[i] = float64(p)
	}

	acc, err := metrics.AccuracyScore(yTest, yPred)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Accuracy: %.2f%%\n", acc*100)
	fmt.Println("Predictions:", preds)

	// Hot path: single sample prediction (faster, no extra allocations)
	label, err := rf.PredictSingle([]float64{7.7, 5.1})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Single sample:", label)

	// Probabilities per class (to filter by confidence)
	proba, err := rf.PredictProbaSingle([]float64{7.7, 5.1})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Class 1 confidence: %.1f%%\n", proba[1]*100)
}
```

---

## 3. Linear models: scale first

Linear models require normalized features. **Fit the scaler only on train**:

```go
scaler := preprocessing.NewStandardScaler()
XTrainS, err := scaler.FitTransform(XTrain) // Fit only on train
XTestS, err  := scaler.Transform(XTest)       // test uses same parameters

lr := linear_model.NewLinearRegression()
if err := lr.Fit(XTrainS, yTrain); err != nil {
	log.Fatal(err)
}
preds, err := lr.Predict(XTestS)
r2, err := lr.Score(XTestS, yTest) // R² for regression
```

> Calling `Transform` before `Fit` returns an explicit error: it never fails silently.

---

## 4. Evaluate with metrics

```go
import "github.com/ateto3/gopherml/metrics"

// Classification
acc, _ := metrics.AccuracyScore(yTrue, yPred)
precision, _ := metrics.PrecisionScore(yTrue, yPred, "macro") // "macro", "weighted" or ""
recall, _ := metrics.RecallScore(yTrue, yPred, "macro")
f1, _ := metrics.F1Score(yTrue, yPred, "macro")
auc, _ := metrics.ROCAUCScore(yTrue, yScores)
report, _ := metrics.ClassificationReport(yTrue, yPred)

// Regression
mse, _ := metrics.MeanSquaredError(yTrue, yPred)
rmse, _ := metrics.RootMeanSquaredError(yTrue, yPred)
mae, _ := metrics.MeanAbsoluteError(yTrue, yPred)
r2, _ := metrics.R2Score(yTrue, yPred)
```

---

## 5. Save and load the model

```go
import "github.com/ateto3/gopherml/persistence"

// Save (gob format)
if err := persistence.Save("model.gob", rf); err != nil {
	log.Fatal(err)
}

// Load into a newly created model — integrity is validated automatically
loaded := ensemble.NewRandomForestClassifier()
if err := persistence.Load("model.gob", loaded); err != nil {
	// The error detects corrupt files or untrained models
	log.Fatal(err)
}
```

---

## 6. Cross-validation before trusting the model

```go
import "github.com/ateto3/gopherml/model_selection"

kf := model_selection.KFold{NSplits: 5, Shuffle: true, RandomState: 42}
scores, err := model_selection.CrossValScore(rf, X, y, kf)
if err != nil {
	log.Fatal(err)
}
fmt.Println("Scores per fold:", scores)
```

---

## 7. Next steps

- **Crypto pipeline:** [USAGE.md § crypto package](../USAGE.md)
- **Hyperparameter search:** [examples/grid_search](../examples/grid_search/)
- **Complete reference:** [API_REFERENCE.md](API_REFERENCE.md)
- **Performance and hot path:** [PERFORMANCE.md](PERFORMANCE.md)
