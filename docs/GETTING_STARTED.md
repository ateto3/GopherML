# 🚀 Inicio rápido con GopherML

Guía práctica para entrenar tu primer modelo con GopherML, desde cero hasta guardarlo en disco. Todos los ejemplos son programas completos y ejecutables.

**Requisitos:** Go 1.23+ · ver [README.md](../README.md) para la instalación.

---

## 1. Conceptos de 30 segundos

GopherML replica la API de scikit-learn con tipos de Go:

| Concepto | Tipo en GopherML | Forma |
| :--- | :--- | :--- |
| **X** (features) | `[][]float64` | `[nMuestras][nFeatures]` |
| **y** (etiquetas) | `[]float64` | `[nMuestras]` |
| **Predicciones** | `[]float64` o `[]int` | según el modelo |

Ciclo de vida de un modelo:

```
NewModel() → Fit(X, y) → Predict(X) / PredictSingle(x) → Score(X, y)
```

> ⚠️ Todos los métodos `Fit` y `Predict` validan que **ningún dato sea `NaN` ni `Inf`** y devuelven error si las dimensiones no coinciden. Nunca ignores el error.

---

## 2. Tu primer modelo (clasificación)

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

	// División train/test 75/25 con RNG determinista
	XTrain, XTest, yTrain, yTest, err := utils.TrainTestSplit(X, y, 0.25, 42)
	if err != nil {
		log.Fatal(err)
	}

	// Entrenar un Random Forest
	rf := ensemble.NewRandomForestClassifier()
	rf.NEstimators = 100
	rf.MaxDepth = 4
	rf.RandomState = 42
	if err := rf.Fit(XTrain, yTrain); err != nil {
		log.Fatal(err)
	}

	// Predecir y convertir a []float64 para las métricas
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
	fmt.Println("Predicciones:", preds)

	// Hot path: predicción de una sola muestra (más rápida, sin allocations extra)
	label, err := rf.PredictSingle([]float64{7.7, 5.1})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Predicción única:", label)

	// Probabilidades por clase (para filtrar por confianza)
	proba, err := rf.PredictProbaSingle([]float64{7.7, 5.1})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Confianza clase 1: %.1f%%\n", proba[1]*100)
}
```

---

## 3. Modelos lineales: escala primero

Los modelos lineales exigen características normalizadas. **Ajusta el scaler solo con train**:

```go
scaler := preprocessing.NewStandardScaler()
XTrainS, err := scaler.FitTransform(XTrain) // Fit solo en train
XTestS, err := scaler.Transform(XTest)       // test usa los mismos parámetros

lr := linear_model.NewLinearRegression()
if err := lr.Fit(XTrainS, yTrain); err != nil {
	log.Fatal(err)
}
preds, err := lr.Predict(XTestS)
r2, err := lr.Score(XTestS, yTest) // R² para regresión
```

> Llamar a `Transform` antes de `Fit` devuelve un error explícito: nunca falla en silencio.

---

## 4. Evaluar con métricas

```go
import "github.com/ateto3/gopherml/metrics"

// Clasificación
acc, _ := metrics.AccuracyScore(yTrue, yPred)
precision, _ := metrics.PrecisionScore(yTrue, yPred, "macro") // "macro", "weighted" o ""
recall, _ := metrics.RecallScore(yTrue, yPred, "macro")
f1, _ := metrics.F1Score(yTrue, yPred, "macro")
auc, _ := metrics.ROCAUCScore(yTrue, yScores)
report, _ := metrics.ClassificationReport(yTrue, yPred)

// Regresión
mse, _ := metrics.MeanSquaredError(yTrue, yPred)
rmse, _ := metrics.RootMeanSquaredError(yTrue, yPred)
mae, _ := metrics.MeanAbsoluteError(yTrue, yPred)
r2, _ := metrics.R2Score(yTrue, yPred)
```

---

## 5. Guardar y cargar el modelo

```go
import "github.com/ateto3/gopherml/persistence"

// Guardar (formato gob)
if err := persistence.Save("modelo.gob", rf); err != nil {
	log.Fatal(err)
}

// Cargar en un modelo recién creado — la integridad se valida automáticamente
cargado := ensemble.NewRandomForestClassifier()
if err := persistence.Load("modelo.gob", cargado); err != nil {
	// El error detecta archivos corruptos o modelos no entrenados
	log.Fatal(err)
}
```

---

## 6. Validación cruzada antes de confiar en el modelo

```go
import "github.com/ateto3/gopherml/model_selection"

kf := model_selection.KFold{NSplits: 5, Shuffle: true, RandomState: 42}
scores, err := model_selection.CrossValScore(rf, X, y, kf)
if err != nil {
	log.Fatal(err)
}
fmt.Println("Scores por fold:", scores)
```

---

## 7. Siguientes pasos

- **Pipeline de cripto:** [USAGE.md § Paquete crypto](../USAGE.md)
- **Búsqueda de hiperparámetros:** [examples/grid_search](../examples/grid_search/)
- **Referencia completa:** [API_REFERENCE.md](API_REFERENCE.md)
- **Rendimiento y hot path:** [PERFORMANCE.md](PERFORMANCE.md)
