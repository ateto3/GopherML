# 🏗️ Arquitectura de GopherML

Este documento describe la arquitectura del proyecto, los patrones de diseño y las decisiones técnicas. Para una explicación más narrativa de los algoritmas, ver [EXPLAINING.md](../EXPLAINING.md).

---

## 1. Filosofía

1. **Espejo de scikit-learn:** paquetes, nombres de métodos y flujos idénticos a `sklearn` para que la curva de aprendizaje sea cero si vienes de Python.
2. **Cero dependencias:** exclusivamente librería estándar de Go. Sin Cgo, sin NumPy, sin bindings.
3. **Eficiencia de datos:** *flat arrays*, memoria pre-asociada y acceso directo para mejorar la localidad de caché.
4. **Concurrencia donde compensa:** goroutines solo en operaciones con trabajo suficiente para justificar el coste de sincronización.
5. **Fallar rápido y ruidoso:** toda entrada inválida (`NaN`, `Inf`, dimensiones inconsistentes, modelo sin entrenar) produce un `error` explícito — nunca resultados silenciosamente erróneos.

---

## 2. Capas del sistema

```
                ┌─────────────────────────────────────────────┐
   Capa alta    │  crypto/   (OHLCV → features → señal)       │  ← orquestación
                ├─────────────────────────────────────────────┤
                │  model_selection/   (KFold, CV, GridSearch) │
                │  persistence/       (Save/Load + Validator) │
                ├─────────────────────────────────────────────┤
   Modelos      │  ensemble/ (RandomForest) ──► tree/ (CART)   │
                │  linear_model/ (OLS, Ridge, Lasso, Logistic,│
                │                  ElasticNet, SGD)           │
                ├─────────────────────────────────────────────┤
   Soporte      │  preprocessing/ (scalers, encoders)         │
                │  metrics/ (accuracy, F1, R², ROC-AUC...)    │
                ├─────────────────────────────────────────────┤
   Base         │  utils/ (CheckFinite, TrainTestSplit, RNG)  │  ← todo pasa por aquí
                └─────────────────────────────────────────────┘
```

Regla de dependencias: **una capa solo puede usar las que están debajo**. `utils` no importa de nadie; `crypto` puede usar todo lo demás.

---

## 3. Flujo de datos de un entrenamiento

```
Datos crudos ([][]float64)
      │
      ▼
utils.CheckXy / CheckFinite      ← valida NaN/Inf y dimensiones
      │
      ▼
preprocessing (opcional)         ← Fit en train, Transform en train y test
      │
      ▼
model.Fit(X, y)                  ← cada modelo vuelve a validar internamente
      │                              (CheckFinite + CheckIsFitted)
      ▼
model.Predict(X) / PredictSingle ← inferencia batch o de una muestra
      │
      ▼
metrics.* (yTrue, yPred)         ← evaluación
      │
      ▼
persistence.Save("modelo.gob")   ← gob + validación de integridad al cargar
```

---

## 4. Patrones de diseño

### 4.1 Struct embedding (heredaría a lo Python)

```go
type LinearModel struct { Coef []float64; Intercept float64 /* ... */ }

type Lasso struct {
    LinearModel          // ← embebe el modelo base
    Alpha float64
    // ...
}
```

`Ridge`, `Lasso`, `ElasticNet` y `LinearRegression` embeben `LinearModel` y reutilizan su lógica de resolución de sistemas lineales (Eliminación Gaussiana). Los hiperparámetros de Random Forest viven en `ensemble.BaseForest`.

**Convención:** los coeficientes se acceden directamente (`model.Coef`), sin getters.

### 4.2 Interfaces pequeñas

| Interfaz | Paquete | Propósito |
| :--- | :--- | :--- |
| `Estimator` (`Fit` + `Score`) | `model_selection` | Cualquier modelo evaluável por `CrossValScore`/`GridSearchCV` |
| `Validator` (`ValidateLoaded`) | `persistence` | Auto-validación de integridad tras `Load` |
| `Criterion`, `Splitter`, `TreeBuilder` | `tree` | Estrategias intercambiables de partición (Gini/Entropía/MSE, mejor split/random) |

### 4.3 Constructores con defaults seguros

Siempre `New...()`:

```go
rf := ensemble.NewRandomForestClassifier() // NEstimators, MaxDepth, RandomState inicializados
```

Un struct literal vacío (`&ensemble.RandomForestClassifier{}`) puede entrenar con hiperparámetros cero.

### 4.4 Errores sentinela

Paquetes como `tree` y `linear_model` exportan errores comparables con `errors.Is`:

```go
var ErrNotFitted         = errors.New("tree is not fitted yet")
var ErrDimensionMismatch = errors.New("dimension mismatch")
var ErrEmptyData         = errors.New("empty input data")
```

---

## 5. Paralelismo inteligente

GopherML no paraleliza por paralelizar:

| Operación | ¿Paralela? | Mecanismo |
| :--- | :--- | :--- |
| Entrenamiento de Random Forest | ✅ | Una goroutine por árbol + `sync.WaitGroup` |
| `GridSearchCV` | ✅ | Una goroutine por combinación de hiperparámetros (`NJobs`) |
| Álgebra lineal (Gram matrix, `MatVec`) | ✅ | Bloques de filas vía `ParallelFor` / `WorkerPool` |
| Predicción de una muestra (`PredictSingle`) | ❌ | El overhead de goroutines supera al trabajo útil |
| Escalado de datasets pequeños | ❌ | Secuencial es más rápido |

Detalles y benchmarks en [PERFORMANCE.md](PERFORMANCE.md) y [../linear_model/PERFORMANCE.md](../linear_model/PERFORMANCE.md).

---

## 6. Diseño del paquete `tree` (base de `ensemble`)

- **Nodos indexados:** el árbol guarda los nodos en un `[]Node` plano con índices enteros en lugar de punteros entre structs → menos presión sobre el recolector de basura y mejor localidad de caché.
- **Constructor DFS iterativo:** `DepthFirstTreeBuilder` usa una pila explícita (`StackRecord`) en lugar de recursión → imposible sufrir stack overflow con `MaxDepth` grandes.
- **Splitter con criterio intercambiable:** `BestSplitter` (búsqueda exhaustiva) o `RandomSplitter` (usado por Random Forest para la aleatoriedad por nodo), combinado con `Gini`, `Entropy` o `MSE`.

---

## 7. Persistencia

- Formato **`gob`** de la librería estándar: rápido, nativo de Go, sin dependencias.
- `persistence.Save` serializa el modelo completo; `Load` lo deserializa en un modelo que tú proporcionas (el tipo debe coincidir).
- Si el modelo implementa `Validator`, `Load` ejecuta `ValidateLoaded()` automáticamente y falla si el modelo quedó a medias (archivo corrupto, estructura incompleta).
- `CryptoPredictor.SaveModel` guarda internamente modelo **y** scaler juntos, evitando el error clásico de desajuste entre train e inferencia.

---

## 8. Convenciones de código

- Comentarios de doc en **inglés** que empiezan por el nombre del símbolo (`// RandomForestClassifier implements...`).
- Nombres de paquete en `snake_case` estilo sklearn (`linear_model`, `model_selection`).
- Los tests viven en el mismo paquete (`*_test.go`) y son deterministas (`utils.NewRNG(seed)`).
- Nada de estado global: los modelos son structs autocontenidos; las predicciones no mutan estado, por lo que son seguras en paralelo **después** de `Fit`.

Ver [CONTRIBUTING.md](../CONTRIBUTING.md) para el flujo completo de contribución.
