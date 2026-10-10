# GopherML — Guía de Uso

GopherML es una librería de Machine Learning de alto rendimiento para Go, diseñada para bots de criptomonedas y sistemas de inferencia de baja latencia.

**Módulo:** `github.com/ateto3/gopherml`

> 📚 **Documentación relacionada:** [Inicio rápido](docs/GETTING_STARTED.md) · [Referencia de API](docs/API_REFERENCE.md) · [Rendimiento](docs/PERFORMANCE.md) · [FAQ](docs/FAQ.md) · [Contribuir](CONTRIBUTING.md)

---

## 1. Instalación

### Como dependencia local (desarrollo)

Añade un `replace` en el `go.mod` de tu bot:

```go
module github.com/tu-usuario/mi-bot

go 1.23

require github.com/ateto3/gopherml v0.0.0

replace github.com/ateto3/gopherml => /ruta/a/GopherML
```

Luego ejecuta `go mod tidy`.

### Como módulo publicado en GitHub

```bash
go get github.com/ateto3/gopherml@latest
```

---

## 2. Conceptos Clave

### Tipos de datos
- **X (Features):** `[][]float64` — matriz de forma `[nSamples, nFeatures]`
- **y (Targets):** `[]float64` — vector de forma `[nSamples]`
- **OHLCV:** `crypto.OHLCV{Open, High, Low, Close, Volume float64}`

### API estándar de modelos
```
NewModel()          → Inicializar con parámetros por defecto
Fit(X, y)          → Entrenar (valida NaN/Inf automáticamente)
Predict(X)         → Predicción batch
PredictSingle(x)   → Predicción de muestra única (hot path, zero alloc extra)
PredictProba(X)    → Probabilidades por clase (para umbrales de confianza)
Score(X, y)        → Evaluar (R² para regresión, Accuracy para clasificación)
```

> **Importante:** Todos los métodos `Fit` y `Predict` rechazan automáticamente valores `NaN` o `Inf`. Nunca ignorar los errores devueltos.

---

## 3. Paquete `crypto/` — Pipeline para Bots de Cripto

Este es el paquete de más alto nivel. Encapsula todo el pipeline: OHLCV → features → señal de trading.

### Uso rápido (bot en producción)

```go
import "github.com/ateto3/gopherml/crypto"

// 1. Crear predictor con configuración
cfg := crypto.DefaultPredictorConfig()
cfg.NEstimators   = 200    // Más árboles = más estable
cfg.MaxDepth      = 8      // Evitar overfitting
cfg.WindowSize    = 20     // Últimas 20 velas por predicción
cfg.MinConfidence = 0.65   // Solo operar si >65% de confianza

p := crypto.NewCryptoPredictor(cfg)

// 2. Cargar modelo pre-entrenado al arrancar el bot
if err := p.LoadModel("btcusdt_1h_rf.gob"); err != nil {
    log.Fatalf("Error cargando modelo: %v", err)
}

// 3. En el loop de trading — por cada nueva vela:
signal, confidence, err := p.PredictSignal(lastCandles) // len == WindowSize
if err != nil {
    log.Printf("Predicción fallida: %v", err)
    return
}

switch signal {
case crypto.SignalLong:
    fmt.Printf("LONG con %.1f%% de confianza\n", confidence*100)
    // executeOrder(BUY)
case crypto.SignalShort:
    fmt.Printf("SHORT con %.1f%% de confianza\n", confidence*100)
    // executeOrder(SELL)
case crypto.SignalHold:
    // No operar
}
```

### Señales disponibles
| Constante | Valor | Significado |
|:---|:---|:---|
| `crypto.SignalHold` | 0 | Sin señal clara — no operar |
| `crypto.SignalLong` | 1 | El modelo predice subida de precio |
| `crypto.SignalShort` | 2 | El modelo predice bajada de precio |

### Entrenar y guardar un modelo

```go
import "github.com/ateto3/gopherml/crypto"

cfg := crypto.DefaultPredictorConfig()
cfg.WindowSize = 20

p := crypto.NewCryptoPredictor(cfg)

// candles: serie temporal histórica de OHLCV
// labels: etiqueta por cada ventana (0=hold, 1=long, 2=short)
// len(labels) debe ser len(candles) - WindowSize + 1
if err := p.Fit(candles, labels); err != nil {
    log.Fatal(err)
}

if err := p.SaveModel("btcusdt_1h_rf.gob"); err != nil {
    log.Fatal(err)
}
```

### Extraer features manualmente

```go
cfg := crypto.NewFeatureConfig()
// cfg.Returns = true, cfg.LogReturns = true, cfg.VWAP = true, ...

// Vector de features para una sola ventana
fv, err := crypto.ToFeatureVector(windowCandles, cfg)

// Matriz completa por sliding window (para entrenamiento)
X, err := crypto.BuildFeatureMatrix(allCandles, windowSize, cfg)
```

**Features incluidas por defecto:**
| Feature | Descripción |
|:---|:---|
| `Returns` | Retornos relativos entre velas consecutivas |
| `LogReturns` | Log-retornos — más estacionarios |
| `PriceRange` | `(high-low)/close` por vela — volatilidad normalizada |
| `VolumeRatio` | `volume/avg_volume` — volumen relativo a la media |
| `Momentum` | Fracción de velas alcistas en la ventana |
| `VWAP` | Desviación del cierre respecto al VWAP de la ventana |

---

## 4. Modelos Individuales

### Random Forest — Clasificador con umbral de confianza

```go
import "github.com/ateto3/gopherml/ensemble"

rf := ensemble.NewRandomForestClassifier()
rf.NEstimators = 100
rf.MaxDepth    = 10
rf.RandomState = 42

if err := rf.Fit(X, y); err != nil {
    log.Fatal(err)
}

// Predicción batch
predictions, _ := rf.Predict(XTest)

// Probabilidades por clase (para filtrar por confianza)
proba, _ := rf.PredictProba(XTest)
// proba[i][1] = probabilidad de clase 1 para la muestra i

// Predicción de muestra única (hot path en bots)
label, _ := rf.PredictSingle(featureVector)
probaSingle, _ := rf.PredictProbaSingle(featureVector)
```

### Random Forest — Regresor (predicción de precio)

```go
rf := ensemble.NewRandomForestRegressor()
rf.NEstimators = 100

rf.Fit(X, y) // y = precio futuro

// Predicción de muestra única (sin overhead de batch)
predictedPrice, _ := rf.PredictSingle(featureVector)
```

### Regresión Lineal

```go
import "github.com/ateto3/gopherml/linear_model"

model := linear_model.NewLinearRegression()
model.Fit(X, y)
preds, _ := model.Predict(XTest)
```

---

## 5. Preprocesamiento

### StandardScaler (recomendado para modelos lineales)

```go
import "github.com/ateto3/gopherml/preprocessing"

scaler := preprocessing.NewStandardScaler()

// Opción A: Fit + Transform por separado (train/test split)
scaler.Fit(XTrain)          // Solo fitear en train
XTrainScaled, _ := scaler.Transform(XTrain)
XTestScaled, _  := scaler.Transform(XTest)  // Mismos parámetros

// Opción B: FitTransform (más conveniente, mismo resultado)
XScaled, _ := scaler.FitTransform(XTrain)

// Invertir la transformación
XOriginal, _ := scaler.InverseTransform(XScaled)
```

> **Atención:** Llamar `Transform` antes de `Fit` devuelve un error explícito — nunca falla silenciosamente.

### MinMaxScaler

```go
scaler := preprocessing.NewMinMaxScaler()
XScaled, _ := scaler.FitTransform(X) // Escala a [0, 1]
```

---

## 6. Persistencia de Modelos

```go
import "github.com/ateto3/gopherml/persistence"
import "github.com/ateto3/gopherml/ensemble"

// Guardar
rf := ensemble.NewRandomForestClassifier()
rf.Fit(X, y)
persistence.Save("modelo.gob", rf)

// Cargar (con validación automática de integridad)
rfLoaded := ensemble.NewRandomForestClassifier()
if err := persistence.Load("modelo.gob", rfLoaded); err != nil {
    // El error incluye detección de archivos corruptos o modelos no entrenados
    log.Fatal(err)
}
```

> Los modelos que implementan la interfaz `persistence.Validator` son validados automáticamente después de cargar. `RandomForestClassifier`, `RandomForestRegressor`, etc. ya implementan esta interfaz.

---

## 7. Validación de Datos

```go
import "github.com/ateto3/gopherml/utils"

// Validar matriz 2D y vector objetivo (incluye NaN/Inf check)
if err := utils.CheckXy(X, y); err != nil {
    log.Fatal(err)
}

// Validar solo una matriz
if err := utils.CheckFinite(X); err != nil {
    log.Fatal(err)
}

// Validar un vector de features individual (hot path)
if err := utils.CheckFiniteSingle(featureVector); err != nil {
    log.Fatal(err)
}
```

---

## 8. Selección de Modelos

```go
import (
    "fmt"
    "log"

    "github.com/ateto3/gopherml/linear_model"
    "github.com/ateto3/gopherml/model_selection"
)

// Validación cruzada K-Fold (la API es un struct KFold, no un entero)
kf := model_selection.KFold{NSplits: 5, Shuffle: true, RandomState: 42}
scores, err := model_selection.CrossValScore(model, X, y, kf)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Scores por fold:", scores)

// Grid Search paralelo: una fábrica de modelos + lista de combinaciones
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
    NJobs: -1, // -1 = todos los núcleos
}
if err := gs.Fit(X, y); err != nil {
    log.Fatal(err)
}
fmt.Printf("Mejores parámetros: %v (score: %.4f)\n", gs.BestParams, gs.BestScore)
```

> **Nota:** `CrossValScore` y `GridSearchCV` requieren modelos con `Fit(X, y)` y `Score(X, y)` (interfaz `Estimator`). Los árboles individuales de `tree/` usan `Fit(X, y, sampleWeight)`, así que usa los bosques de `ensemble/`.

---

## 9. Buenas Prácticas para Bots en Producción

| Práctica | Motivo |
|:---|:---|
| Usar `PredictSingle` en loops de alta frecuencia | Evita allocations innecesarias por tick |
| Cargar el modelo al arrancar, no en cada predicción | `Load` es costoso — hacerlo una vez |
| Usar `MinConfidence >= 0.60` | Reduce señales falsas en mercados laterales |
| Validar candles antes de predecir | Exchanges pueden enviar velas con gaps o volumen 0 |
| Guardar el scaler junto al modelo | El `CryptoPredictor` lo hace automáticamente via `SaveModel` |
| Usar `go test ./... -race` antes de desplegar | Detecta race conditions en predicción concurrente |

---

*GopherML — Machine Learning nativo en Go para bots de criptomonedas.*
