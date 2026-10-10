# 🧪 Ejemplos de GopherML

Programas ejecutables (`package main`) que demuestran el uso completo de la librería. Todos funcionan con datos incluidos en el propio código: no necesitan descargas externas.

**Requisitos:** Go 1.23+ y el módulo `github.com/ateto3/gopherml` disponible (ver [README](../README.md#2-instalación)).

---

## Ejemplos disponibles

### 1. Clasificación básica

**Archivo:** [`classification/basic_classification.go`](classification/basic_classification.go)

Random Forest para clasificación binaria: escalado con `StandardScaler`, entrenamiento, predicción *batch* y de muestra única (`PredictSingle`), probabilidades (`PredictProba`), validación cruzada con `KFold` y métricas detalladas (precision, recall, F1).

```bash
go run examples/classification/basic_classification.go
```

### 2. Regresión básica

**Archivo:** [`regression/basic_regression.go`](regression/basic_regression.go)

Comparación de `LinearRegression` y `RandomForestRegressor` en un problema de predicción de precios: escalado, predicciones, R², MSE/MAE y comparación entre modelos.

```bash
go run examples/regression/basic_regression.go
```

### 3. Pipeline de trading de cripto

**Archivo:** [`crypto/crypto_trading_bot.go`](crypto/crypto_trading_bot.go)

Pipeline completo con el paquete `crypto/`: genera velas OHLCV sintéticas, construye etiquetas por ventana deslizante, entrena un `CryptoPredictor`, guarda y recarga el modelo (`SaveModel`/`LoadModel`) y emite una señal `LONG`/`SHORT`/`HOLD` con su nivel de confianza.

```bash
go run examples/crypto/crypto_trading_bot.go
```

### 4. Persistencia de modelos

**Archivo:** [`persistence/model_persistence.go`](persistence/model_persistence.go)

Entrena un Random Forest, lo guarda en `gob`, lo recarga en un modelo nuevo y verifica que las predicciones coinciden. Además demuestra la detección de archivos corruptos mediante la interfaz `Validator`.

```bash
go run examples/persistence/model_persistence.go
```

### 5. Búsqueda de hiperparámetros

**Archivo:** [`grid_search/hyperparameter_tuning.go`](grid_search/hyperparameter_tuning.go)

`GridSearchCV` paralelo (`NJobs: -1`) sobre `LogisticRegression`: fábrica de modelos, lista de combinaciones, puntuaciones por combinación y selección de los mejores parámetros.

```bash
go run examples/grid_search/hyperparameter_tuning.go
```

### 6. Pipeline de preprocesamiento

**Archivo:** [`preprocessing/preprocessing_pipeline.go`](preprocessing/preprocessing_pipeline.go)

Uso de todos los transformadores: `StandardScaler` (incluyendo `InverseTransform`), `MinMaxScaler`, `LabelEncoder` (texto → enteros) y `PolynomialFeatures`, más el manejo de errores al usar `Transform` antes de `Fit`.

```bash
go run examples/preprocessing/preprocessing_pipeline.go
```

---

## Ejecutar todos

```bash
for f in examples/*/*.go; do
  echo "▶ $f"
  go run "$f" > /dev/null || exit 1
done
echo "✔ Todos los ejemplos ejecutan correctamente"
```

## Notas

- Cada directorio es un paquete independiente: puedes ejecutar un archivo suelto con `go run <ruta>` sin afectar al resto.
- Los `.gob` que se crean durante la ejecución (persistencia/cripto) se borran automáticamente al finalizar y están excluidos del repositorio vía `.gitignore`.
- Para ejemplos de rendimiento y *benchmarking*, ver [docs/PERFORMANCE.md](../docs/PERFORMANCE.md).

