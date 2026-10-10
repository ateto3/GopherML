# 📖 Referencia de API de GopherML

Referencia de los tipos y funciones exportados de cada paquete, generada a partir del código fuente (Go 1.23, módulo `github.com/ateto3/gopherml`).

> **Convención de la API estándar de modelos:**
>
> ```
> NewModel()          → inicializar con parámetros por defecto
> Fit(X, y)           → entrenar (valida NaN/Inf y dimensiones)
> Predict(X)          → predicción batch
> PredictSingle(x)    → predicción de una muestra (hot path)
> PredictProba(X)     → probabilidades por clase
> Score(X, y)         → evaluar (R² en regresión, accuracy en clasificación)
> ```
>
> `X` es `[][]float64` (`[nMuestras][nFeatures]`), `y` es `[]float64`. Nunca ignores los errores.

---

## `utils` — Validación y utilidades de datos

Toda entrada de datos debe pasar por aquí (los modelos lo hacen internamente).

| Función | Firma | Descripción |
| :--- | :--- | :--- |
| `CheckArray` | `(X [][]float64) error` | Valida que la matriz no esté vacía y tenga filas homogéneas |
| `CheckFinite` | `(X [][]float64) error` | Rechaza `NaN`/`Inf` en toda la matriz |
| `CheckFiniteSingle` | `(x []float64) error` | Ídem para un vector (hot path, más rápido) |
| `CheckFiniteVec` | `(y []float64) error` | Ídem para el vector objetivo |
| `CheckXy` | `(X [][]float64, y []float64) error` | Valida matriz + vector: dimensiones y finite |
| `CheckIsFitted` | `(fitted bool, modelName string) error` | Devuelve error si el modelo no está entrenado |
| `TrainTestSplit` | `(X, y, testSize float64, randomState int64) (XTrain, XTest, yTrain, yTest [][]float64/[]float64, error)` | División train/test; `testSize ∈ (0,1)` |
| `LoadCSV` | `(path string, hasHeader bool) (X [][]float64, y []float64, error)` | Carga CSV: última columna = target |
| `NewRNG` | `(seed int64) *rand.Rand` | RNG determinista para tests reproducibles |
| `Shuffle` | `(n int, rng *rand.Rand) []int` | Permutación determinista |

---

## `preprocessing` — Escalado y codificación

Todos siguen el patrón `Fit` → `Transform` (o `FitTransform`). `Transform` antes de `Fit` devuelve error.

### `StandardScaler`

Estandariza restando la media y dividiendo por la desviación típica (equivalente a `sklearn.preprocessing.StandardScaler`).

```go
s := preprocessing.NewStandardScaler()   // *StandardScaler{Mean, Std, IsFitted}
s.Fit(X) error
s.Transform(X) ([][]float64, error)
s.FitTransform(X) ([][]float64, error)
s.InverseTransform(X) ([][]float64, error)
```

### `MinMaxScaler`

Escala cada feature al rango `[0, 1]` (`Min`, `Max`, `IsFitted`). Mismos métodos que `StandardScaler`.

### `LabelEncoder`

Codifica etiquetas de texto a enteros `0..n-1`:

```go
le := preprocessing.NewLabelEncoder()    // *LabelEncoder{Classes []string, Mapping map[string]int}
le.Fit(y []string)
le.Transform(y []string) ([]float64, error)
le.InverseTransform(y []float64) ([]string, error)
```

### `PolynomialFeatures`

Genera features polinómicas y de interacción **hasta grado 2** (mayor grado devuelve error):

```go
pf := preprocessing.NewPolynomialFeatures(2)  // *PolynomialFeatures{Degree, IncludeBias}
pf.Fit(X) error
pf.FitTransform(X) ([][]float64, error)
pf.Transform(X) ([][]float64, error)
```

---

## `metrics` — Evaluación

Todas las funciones comparan `yTrue` vs `yPred` (ambos `[]float64`) y devuelven `(float64, error)` salvo indicación.

### Clasificación

| Función | Firma | Notas |
| :--- | :--- | :--- |
| `AccuracyScore` | `(yTrue, yPred []float64) (float64, error)` | Fracción de aciertos |
| `PrecisionScore` | `(yTrue, yPred []float64, average string) (float64, error)` | `average`: `"macro"`, `"weighted"` o `""` (binario) |
| `RecallScore` | `(yTrue, yPred []float64, average string) (float64, error)` | Ídem |
| `F1Score` | `(yTrue, yPred []float64, average string) (float64, error)` | Ídem |
| `ROCAUCScore` | `(yTrue, yScores []float64) (float64, error)` | `yScores` = probabilidades de la clase positiva |
| `ConfusionMatrix` | `(yTrue, yPred []float64) ([][]int, []float64, error)` | Matriz de confusión + clases |
| `ClassificationReport` | `(yTrue, yPred []float64) (string, error)` | Informe por clase estilo sklearn |

### Regresión

| Función | Firma |
| :--- | :--- |
| `MeanSquaredError` | `(yTrue, yPred []float64) (float64, error)` |
| `RootMeanSquaredError` | `(yTrue, yPred []float64) (float64, error)` |
| `MeanAbsoluteError` | `(yTrue, yPred []float64) (float64, error)` |
| `R2Score` | `(yTrue, yPred []float64) (float64, error)` |

> Los modelos también incluyen su propio `Score(X, y)`: **accuracy** para clasificadores, **R²** para regresores.

---

## `linear_model` — Regresión y clasificación lineal

### Errores sentinela

```go
linear_model.ErrEmptyData         // entrada vacía
linear_model.ErrDimensionMismatch // dimensiones incompatibles
linear_model.ErrNotFitted         // Predict antes de Fit
```

### Modelos

| Modelo | Constructor | Métodos | Hiperparámetros destacados |
| :--- | :--- | :--- | :--- |
| `LinearRegression` (OLS) | `NewLinearRegression()`, `NewLinearRegressionWithParams(LinearRegressionParams)` | `Fit`, `Predict`, `Score` | `FitIntercept`, `CopyX` |
| `Ridge` (L2) | `NewRidge()`, `NewRidgeWithParams(RidgeParams)` | `Fit`, `Predict`, `Score` | `Alpha`, `Solver`, `MaxIter`, `Tol` |
| `Lasso` (L1) | `NewLasso()`, `NewLassoWithParams(LassoParams)` | `Fit`, `Predict`, `Score` | `Alpha`, `MaxIter`, `Tol`, `Positive` |
| `ElasticNet` (L1+L2) | `NewElasticNet()`, `NewElasticNetWithParams(ElasticNetParams)` | `Fit`, `Predict`, `Score` | `Alpha`, `L1Ratio`, `MaxIter`, `Tol` |
| `LogisticRegression` | `NewLogisticRegression()` | `Fit`, `Predict`, `PredictProba`, `Score` | `LearningRate`, `MaxIter`, `Tol` (binaria) |
| `SGDClassifier` | `NewSGDClassifier()` | `Fit`, `Predict` | `LearningRate`, `MaxIter`, `Tol`, `RandomState` |

`LinearRegression`, `Ridge`, `Lasso` y `ElasticNet` embeben `LinearModel`:

```go
type LinearModel struct {
    Coef         []float64 // coeficientes
    Intercept    float64   // sesgo
    FitIntercept bool
    NFeaturesIn  int
    IsFitted     bool
}
// métodos: Predict(X) ([]float64, error), Score(X, y) (float64, error)
```

### Versiones paralelas y utilidades

| Símbolo | Firma / descripción |
| :--- | :--- |
| `ParallelOptions` | `{NumWorkers, MinBatchSize int}` — `DefaultParallelOptions()` usa todos los núcleos |
| `ParallelFor` | `(n int, opts ParallelOptions, fn func(start, end int))` — paraleliza un rango |
| `ParallelMatrixMultiply` | `(X [][]float64, opts) [][]float64` |
| `ParallelMatVecMul` | `(X [][]float64, v []float64, opts) []float64` |
| `ParallelCenterData` | centra datos en paralelo, devuelve offsets |
| `LinearRegressionParallel` | `NewLinearRegressionParallel()` — OLS paralelo |
| `LassoParallel` | `NewLassoParallel()`, `NewLassoParallelWithParams(LassoParams, ParallelOptions)` |
| `WorkerPool` | `NewWorkerPool(workers int)` reutilizable para lotes |
| `DotProduct`, `Mean`, `Transpose`, `MeanSquaredError` | utilidades matemáticas |

Guía de cuándo paralelizar en [PERFORMANCE.md](PERFORMANCE.md) y [../linear_model/PERFORMANCE.md](../linear_model/PERFORMANCE.md).

---

## `tree` — Árboles de decisión CART

> ⚠️ **Atención:** `Fit` de los árboles recibe un tercer argumento `sampleWeight` (usa `nil` si no hay pesos). Por eso **no** implementan directamente la interfaz `Estimator` de `model_selection`; usa los bosques de `ensemble` para validación cruzada.

### `DecisionTreeClassifier`

```go
dt := tree.NewDecisionTreeClassifier()
// o con parámetros explícitos:
dt := tree.NewDecisionTreeClassifierWithParams(tree.DecisionTreeClassifierParams{ /* ... */ })

dt.Fit(X [][]float64, y []float64, sampleWeight []float64) error
dt.Predict(X) ([]int, error)
dt.PredictProba(X) ([][]float64, error)
dt.Score(X, y) (float64, error)
dt.GetDepth() int
dt.GetNLeaves() int
```

### `DecisionTreeRegressor`

Mismos métodos, pero `Predict` devuelve `([]float64, error)` y el criterio por defecto es `MSE`.

### Parámetros comunes (campos del struct)

| Campo | Tipo | Descripción |
| :--- | :--- | :--- |
| `Criterion` | `string` | `"gini"`, `"entropy"` (clasif.) o `"mse"` (regres.) |
| `Splitter` | `string` | Estrategia de búsqueda de split |
| `MaxDepth` | `int` | Profundidad máxima (`0` = sin límite) |
| `MinSamplesSplit` | `int` | Mínimo de muestras para dividir un nodo |
| `MinSamplesLeaf` | `int` | Mínimo de muestras por hoja |
| `MaxFeatures` | `interface{}` | `int`, `float64` o `"sqrt"`/`"log2"` |
| `RandomState` | `int` | Semilla de reproducibilidad |
| `CCPAlpha` | `float64` | Poda de coste-complejidad |

Atributos ajustados tras `Fit`: `Tree`, `NFeatures`, `NOutputs`, `NClasses`, `Classes`, `IsFitted`.

### Componentes extensibles

`Criterion` (`Gini`, `Entropy`, `MSE`), `Splitter` (`BestSplitter`, `RandomSplitter`), `TreeBuilder` (`DepthFirstTreeBuilder`) y `BuildTree(...)` para construir árboles a bajo nivel. `PrintTree(t, featureNames)` imprime el árbol.

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
rf.ValidateLoaded() error                             // llamado por persistence.Load
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

### Hiperparámetros (embebidos de `BaseForest`)

| Campo | Tipo | Descripción |
| :--- | :--- | :--- |
| `NEstimators` | `int` | Número de árboles (entrenados en paralelo) |
| `MaxDepth` | `int` | Profundidad máxima por árbol |
| `MinSamplesSplit` / `MinSamplesLeaf` | `int` | Restricciones de tamaño de nodo/hoja |
| `MaxFeatures` | `interface{}` | Features consideradas por split |
| `Bootstrap` | `bool` | Muestreo con reemplazo por árbol |
| `RandomState` | `int64` | Semilla (reproducibilidad) |
| `IsFitted` | `bool` | Estado tras entrenar |

Atributos: `Classes []int`, `NClasses int`, `Estimators []…`.

Utilidades: `ForestParams` / `DefaultForestParams()`, `GenerateBootstrapIndices(n, rng)`, `GetBootstrapSample(X, y, rng)`.

---

## `model_selection` — Validación y búsqueda

### `KFold`

```go
kf := model_selection.KFold{NSplits: 5, Shuffle: true, RandomState: 42}
splits := kf.Split(nSamples) // []SplitRecord{TrainIndices, TestIndices}
```

### `CrossValScore`

```go
func CrossValScore(model Estimator, X [][]float64, y []float64, cv KFold) ([]float64, error)
```

`Estimator` es cualquier modelo con `Fit(X, y) error` y `Score(X, y) (float64, error)` — p. ej. Random Forest o `LogisticRegression` (los árboles de `tree/`, con `Fit` de 3 argumentos, no cumplen esta interfaz).

### `GridSearchCV`

Búsqueda exhaustiva **paralela** sobre una lista explícita de combinaciones:

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
    NJobs: -1, // -1 = todos los núcleos
}

err := gs.Fit(X, y)
// gs.BestParams, gs.BestScore, gs.CVResults ([]CVResult{Params, MeanScore, Err})
```

---

## `persistence` — Guardar y cargar modelos

| Función | Firma | Descripción |
| :--- | :--- | :--- |
| `Save` | `(filename string, model interface{}) error` | Serializa a `gob` |
| `Load` | `(filename string, model interface{}) error` | Deserializa + valida (ver `Validator`) |
| `Dump` | `(model interface{}, w io.Writer) error` | Serializa a cualquier `io.Writer` |
| `LoadFrom` | `(r io.Reader, model interface{}) error` | Deserializa desde cualquier `io.Reader` |

```go
type Validator interface {
    ValidateLoaded() error
}
```

Los modelos que implementan `Validator` (Random Forest, etc.) se validan automáticamente tras `Load`: detecta archivos corruptos y modelos sin entrenar.

---

## `crypto` — Pipeline OHLCV → señal

### Tipos de datos

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

Campos booleanos (todos `true` con `NewFeatureConfig()`): `Returns`, `LogReturns`, `PriceRange`, `VolumeRatio`, `Momentum`, `VWAP`.

| Función | Firma |
| :--- | :--- |
| `ToFeatureVector` | `(candles []OHLCV, cfg FeatureConfig) ([]float64, error)` — una ventana |
| `BuildFeatureMatrix` | `(candles []OHLCV, windowSize int, cfg FeatureConfig) ([][]float64, error)` — sliding window |
| `FeatureSize` | `(windowSize int, cfg FeatureConfig) int` — nº de columnas resultantes |
| `MinWindowSize` | `() int` — ventana mínima admitida |

### `CryptoPredictor`

```go
type PredictorConfig struct {
    NEstimators    int           // árboles del bosque (default 100)
    MaxDepth       int           // profundidad (default 10)
    WindowSize     int           // velas por ventana (default 20)
    MinConfidence  float64       // umbral para señal no-HOLD (default 0.55)
    Features       FeatureConfig // grupos de features
}
cfg := crypto.DefaultPredictorConfig()

p := crypto.NewCryptoPredictor(cfg)
p.Fit(candles []OHLCV, labels []float64) error  // len(labels) == len(candles)-WindowSize+1
p.PredictSignal(candles []OHLCV) (Signal, float64, error) // len(candles) == WindowSize
p.SaveModel(path string) error
p.LoadModel(path string) error                  // modelo + scaler juntos
p.Config() PredictorConfig
p.IsFitted() bool
```

Guía completa de uso en [../USAGE.md](../USAGE.md).
