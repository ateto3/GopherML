# 🐹 GopherML 🚀

> **Machine Learning nativo en Go: inferencia de baja latencia, cero dependencias, API estilo scikit-learn.**

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/ateto3/gopherml.svg)](https://pkg.go.dev/github.com/ateto3/gopherml)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen)](#5-tests-y-ci)
[![CI](https://github.com/ateto3/gopherml/actions/workflows/ci.yml/badge.svg)](https://github.com/ateto3/gopherml/actions/workflows/ci.yml)

**GopherML** es una librería de Machine Learning escrita íntegramente en **Go**, diseñada para proporcionar herramientas de entrenamiento e inferencia ultra-rápidas aprovechando la concurrencia nativa (goroutines) y **sin Cgo ni dependencias externas**. Sigue la arquitectura y la API de *scikit-learn* de Python, adaptada al mundo de los bots de criptomonedas y la inferencia de alta frecuencia.

---

## 📑 Índice

1. [Características](#1-características)
2. [Instalación](#2-instalación)
3. [Inicio rápido](#3-inicio-rápido)
4. [Paquetes](#4-paquetes)
5. [Tests y CI](#5-tests-y-ci)
6. [Estructura del repositorio](#6-estructura-del-repositorio)
7. [Documentación](#7-documentación)
8. [Ejemplos](#8-ejemplos)
9. [Subproyecto: go-nlp](#9-subproyecto-go-nlp)
10. [Hoja de ruta](#10-hoja-de-ruta)
11. [Métricas de diseño](#11-métricas-de-diseño)
12. [Contribuir y licencia](#12-contribuir-y-licencia)

---

## 1. Características

- **⚡ Inferencia de baja latencia:** `PredictSingle` y `PredictProbaSingle` están optimizados para el *hot path* de bots (cero asignaciones extra por tick).
- **🧩 API estilo scikit-learn:** `NewModel()` → `Fit(X, y)` → `Predict(X)` → `Score(X, y)`. Si ya usas scikit-learn, te sientes en casa.
- **🌊 Concurrencia nativa:** entrenamiento paralelo de Random Forest, `GridSearchCV` y álgebra lineal con `goroutines` y `sync.WaitGroup` — sin Cgo.
- **🛡️ Validación estricta:** todos los datos de entrada se validan contra `NaN`/`Inf` (`utils.CheckFinite`) antes de tocar un modelo. Nunca falla en silencio.
- **💾 Persistencia con integridad:** `Save`/`Load` en formato `gob`, con validación automática post-carga mediante la interfaz `persistence.Validator`.
- **📈 Métricas completas:** Accuracy, Precision, Recall, F1, R², MAE, MSE, RMSE y ROC-AUC.
- **🪙 Pipeline de cripto:** el paquete `crypto/` encapsula el flujo completo **OHLCV → features → señal de trading** (`Hold`/`Long`/`Short`).
- **🚫 Zero dependencies:** solo la librería estándar de Go. Sin NumPy, sin CGO, sin bindings.

---

## 2. Instalación

**Requisitos:** Go 1.23 o superior.

### Como módulo publicado en GitHub

```bash
go get github.com/ateto3/gopherml@latest
```

### Como dependencia local (desarrollo)

Añade un `replace` en el `go.mod` de tu proyecto:

```go
module github.com/tu-usuario/mi-bot

go 1.23

require github.com/ateto3/gopherml v0.0.0

replace github.com/ateto3/gopherml => /ruta/a/GopherML
```

Y ejecuta:

```bash
go mod tidy
```

---

## 3. Inicio rápido

Clasificación binaria completa con Random Forest:

```go
package main

import (
	"fmt"
	"log"

	"github.com/ateto3/gopherml/ensemble"
	"github.com/ateto3/gopherml/preprocessing"
)

func main() {
	// Datos: [longitud, ancho] → etiqueta 0 o 1
	X := [][]float64{
		{1.0, 2.0}, {1.5, 1.8}, {1.2, 2.5}, {1.1, 1.9},
		{5.0, 6.0}, {5.5, 5.8}, {5.2, 6.5}, {4.9, 5.9},
	}
	y := []float64{0, 0, 0, 0, 1, 1, 1, 1}

	// 1. Escalar los datos (imprescindible para modelos lineales)
	scaler := preprocessing.NewStandardScaler()
	Xs, err := scaler.FitTransform(X)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Entrenar
	rf := ensemble.NewRandomForestClassifier()
	rf.NEstimators = 100
	rf.MaxDepth = 5
	rf.RandomState = 42
	if err := rf.Fit(Xs, y); err != nil {
		log.Fatal(err)
	}

	// 3. Predecir (batch)
	preds, err := rf.Predict(Xs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Predicciones:", preds) // [0 0 0 0 1 1 1 1]

	// 4. Evaluar
	acc, err := rf.Score(Xs, y)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Accuracy: %.2f%%\n", acc*100)

	// 5. Hot path: una sola muestra, sin overhead de batch
	label, err := rf.PredictSingle([]float64{5.1, 6.1})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Muestra única:", label)
}
```

**Pipeline completo para bots de cripto** (OHLCV → señal):

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

Ver **[USAGE.md](USAGE.md)** para la guía completa con ejemplos de trading.

---

## 4. Paquetes

| Paquete | Descripción |
| :--- | :--- |
| [`/crypto`](crypto/) | **Pipeline de cripto:** `OHLCV` → ingeniería de features → señal de trading (`CryptoPredictor`) |
| [`/linear_model`](linear_model/) | `LinearRegression`, `Ridge`, `Lasso`, `ElasticNet`, `LogisticRegression`, `SGDClassifier` + versiones paralelas |
| [`/tree`](tree/) | `DecisionTreeClassifier` y `DecisionTreeRegressor` (CART; criterios Gini, Entropía y MSE) |
| [`/ensemble`](ensemble/) | `RandomForestClassifier` y `RandomForestRegressor` (entrenamiento paralelo con bootstrap) |
| [`/metrics`](metrics/) | Accuracy, Precision, Recall, F1, R², MAE, MSE, RMSE, ROC-AUC, reportes y matriz de confusión |
| [`/model_selection`](model_selection/) | `KFold`, `CrossValScore` y `GridSearchCV` paralelo |
| [`/preprocessing`](preprocessing/) | `StandardScaler`, `MinMaxScaler`, `LabelEncoder`, `PolynomialFeatures` |
| [`/persistence`](persistence/) | `Save`/`Load` de modelos con validación de integridad post-carga |
| [`/utils`](utils/) | Validación de datos (`CheckXy`, `CheckFinite`), `TrainTestSplit`, `LoadCSV`, RNG determinista |
| [`/examples`](examples/) | Ejemplos ejecutables de extremo a extremo |

---

---

## 5. Tests y CI

```bash
# Compilar todo
go build ./...

# Analizar código (debe quedar limpio)
go vet ./...

# Tests con detector de condiciones de carrera
go test ./... -race

# Formateo (recomendado antes de un PR)
gofmt -w .
```

Estado actual: **todos los paquetes en verde** y `go vet` limpio. Cada push ejecuta la [workflow de CI](.github/workflows/ci.yml) (`go vet` + `go test -race`).

---

## 6. Estructura del repositorio

```
GopherML/
├── README.md              # Este archivo
├── USAGE.md               # Guía de uso con ejemplos de bots de cripto
├── EXPLAINING.md          # Arquitectura y lógica interna
├── CONTRIBUTING.md        # Cómo contribuir (incluye directrices para IAs)
├── CHANGELOG.md           # Historial de cambios
├── docs/                  # Documentación extendida (guías, API, FAQ)
├── crypto/                # Pipeline OHLCV → features → señal
├── ensemble/              # Random Forest
├── linear_model/          # Modelos lineales y SGD
├── tree/                  # Árboles de decisión CART
├── metrics/               # Métricas de evaluación
├── model_selection/       # Validación cruzada y grid search
├── preprocessing/         # Escalado y encoding
├── persistence/           # Guardar/cargar modelos (gob)
├── utils/                 # Validación y utilidades de datos
├── examples/              # Ejemplos ejecutables
├── go-nlp/                # Submódulo: NLP concurrente (port de NLTK)
└── .github/               # Plantillas de issues/PRs y CI
```

---

## 7. Documentación

| Documento | Contenido |
| :--- | :--- |
| **[docs/GETTING_STARTED.md](docs/GETTING_STARTED.md)** | Tutorial paso a paso desde cero |
| **[USAGE.md](USAGE.md)** | Guía de uso con ejemplos de bots de cripto |
| **[docs/API_REFERENCE.md](docs/API_REFERENCE.md)** | Referencia de la API de todos los paquetes |
| **[EXPLAINING.md](EXPLAINING.md)** | Arquitectura y lógica interna |
| **[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)** | Arquitectura detallada y decisiones de diseño |
| **[docs/PERFORMANCE.md](docs/PERFORMANCE.md)** | Guía de rendimiento y optimizaciones |
| **[docs/FAQ.md](docs/FAQ.md)** | Preguntas frecuentes |
| **[CONTRIBUTING.md](CONTRIBUTING.md)** | Guía de contribución y estándares de código |
| **[CHANGELOG.md](CHANGELOG.md)** | Historial de cambios |

Índice completo en **[docs/README.md](docs/README.md)**.

---

## 8. Ejemplos

Todos los ejemplos son programas ejecutables (`package main`):

```bash
# Clasificación con Random Forest + métricas + validación cruzada
go run examples/classification/basic_classification.go

# Regresión: LinearRegression vs RandomForestRegressor
go run examples/regression/basic_regression.go

# Pipeline de cripto: OHLCV sintéticos → señal → persistencia
go run examples/crypto/crypto_trading_bot.go

# Guardar y cargar un modelo entrenado
go run examples/persistence/model_persistence.go

# Búsqueda paralela de hiperparámetros con GridSearchCV
go run examples/grid_search/hyperparameter_tuning.go

# Pipeline de preprocesamiento con scalers y encoder
go run examples/preprocessing/preprocessing_pipeline.go
```

Detalles en **[examples/README.md](examples/README.md)**.

---

## 9. Subproyecto: go-nlp

El directorio [`go-nlp/`](go-nlp/) contiene un **módulo de Go independiente** (`github.com/ateto3/gopherml/go-nlp`): un puerto concurrente de NLTK con tokenización, tagging POS, parsing, clasificación de texto y análisis de sentimiento, además de optimización de tokens para APIs de LLM.

> ⚠️ `go-nlp` tiene su propio `go.mod`: se compila y testea de forma independiente con
> `cd go-nlp && go build ./... && go test ./...`

Ver **[go-nlp/README.md](go-nlp/README.md)** para la documentación completa.

---

## 10. Hoja de ruta

- [x] Modelos de ML básicos (regresión lineal, clasificación, árboles, bosques, ensamblaje)
- [x] Persistencia de modelos (Save/Load con validación de integridad)
- [x] Métricas de evaluación (accuracy, precision, recall, F1, R², MAE, RMSE)
- [x] Selección de modelos y validación cruzada (KFold, GridSearchCV)
- [x] Preprocesamiento de datos (StandardScaler, MinMaxScaler, LabelEncoder)
- [x] Validación NaN/Inf en todos los datos de entrada
- [x] `PredictProba` + `PredictSingle` para bots de alta frecuencia
- [x] Paquete `crypto/` — pipeline OHLCV → features → señal de trading
- [ ] GradientBoosting (XGBoost-like)
- [ ] Soporte para datos en streaming (actualizaciones online)

---

## 11. Métricas de diseño

| Objetivo | Especificación |
| :--- | :--- |
| **Lenguaje** | Go 1.23+ |
| **Módulo** | `github.com/ateto3/gopherml` |
| **Target** | Inferencia de baja latencia / bots de cripto |
| **Hardware** | Optimizado para ARM64 (Raspberry Pi 4) y x86_64 |
| **Dependencias** | Zero (solo stdlib de Go) |
| **Tests** | `go test ./... -race` — 100% de tests en verde |

---

## 12. Contribuir y licencia

¡Las contribuciones son bienvenidas! Lee **[CONTRIBUTING.md](CONTRIBUTING.md)** para el flujo de trabajo, los estándares de código y las directrices de rendimiento (incluidas las directrices para asistentes IA). Antes de enviar un PR, asegúrate de que pasan:

```bash
go vet ./... && go test ./... -race
```

Este proyecto está bajo licencia **MIT** — ver [LICENSE](LICENSE).

---

## 📝 Notas del desarrollador (Ateto)

Este proyecto es parte de un ecosistema que incluye bots de trading de criptomonedas. El código debe ser modular para que pueda ser importado fácilmente en otros binarios de Go o microservicios de automatización.

> "El código debe ser tan rápido que el recolector de basura de Go ni se entere de que estamos aquí."