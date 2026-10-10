# ❓ Preguntas frecuentes (FAQ)

## Errores y validación

### `Fit`/`Predict` me devuelve `NaN` o `Inf` — ¿por qué?

GopherML rechaza explícitamente los valores no finitos en la entrada. El error indica la posición del problema. Solución típica: revisar el cálculo de features (divisiones por cero, `log(0)`) o limpiar el CSV antes de entrenar. Nunca "parchees" con `math.IsNaN` fuera: la validación es intencionada.

### Me dice `dimension mismatch` — ¿qué compruebo?

`len(X)` (muestras) debe coincidir con `len(y)`, y **todas las filas** de `X` deben tener el mismo número de columnas. Para predicción, el número de columnas debe ser el mismo que en entrenamiento.

### `Transform` devuelve un error antes de `Fit`

Correcto: los scalers necesitan conocer media/desviación (o min/max) con `Fit` (o `FitTransform`) antes de transformar. Nunca falla en silencio.

### `model not fitted` / `tree is not fitted yet`

Estás llamando a `Predict` antes de `Fit`, o cargaste un modelo incompleto. En `persistence.Load`, los modelos con la interfaz `Validator` detectan esto automáticamente.

---

## API y uso

### ¿Cuál es la diferencia entre `Predict` y `PredictSingle`?

`Predict` procesa un lote (`[][]float64`) y reserva memoria para el resultado; `PredictSingle` procesa una muestra con el mínimo de allocations — úsala en bucles de trading (ver [PERFORMANCE.md](PERFORMANCE.md)).

### ¿Por qué `tree.DecisionTreeClassifier.Fit(X, y, sampleWeight)` tiene 3 argumentos?

Los árboles de bajo nivel soportan pesos por muestra. Por eso **no** implementan la interfaz `Estimator` de `model_selection`; para validación cruzada usa `ensemble.RandomForest*`, cuya firma es `Fit(X, y)`.

### ¿Cómo obtengo probabilidades para filtrar por confianza?

```go
proba, _ := rf.PredictProbaSingle(featureVector)
if proba[1] >= 0.65 { /* confianza suficiente */ }
```

En `crypto`, este filtro ya está embebido: `PredictSignal` devuelve `SignalHold` si la confianza no supera `MinConfidence`.

### ¿Cómo hago que los resultados sean reproducibles?

Fija `RandomState` (y `Shuffle: true` en `KFold`). GopherML usa RNG con semilla (`utils.NewRNG`) en lugar del estado global de `math/rand`.

### ¿Qué devuelve `Score`?

Accuracy para clasificadores y R² para regresores (equivalente a scikit-learn).

---

## Persistencia

### ¿Los `.gob` son compatibles entre versiones?

`gob` serializa la estructura concreta del struct: si cambias la versión de GopherML y cambian campos internos de un modelo, reentrena o guarda también la versión. En producción, versiona los archivos de modelo junto al binario.

### ¿Se puede guardar un modelo sin entrenar?

Sí se serializaría, pero `Load` fallará si el modelo implementa `Validator` (que comprueba `IsFitted`). Esa es la protección de integridad.

### ¿El scaler se guarda con el modelo?

Con `CryptoPredictor.SaveModel`, sí (modelo + scaler juntos). Para modelos sueltos, guarda tú el scaler aparte o fitea `StandardScaler` con los mismos datos de entrenamiento.

---

## Rendimiento

### ¿Debo usar las versiones paralelas?

Solo si el dataset es grande (>10 000 muestras o >50 features). Para datasets pequeños, las goroutines añaden overhead. Ver [PERFORMANCE.md](PERFORMANCE.md).

### ¿GopherML usa Cgo o GPU?

No. Solo librería estándar de Go: binarios estáticos, multiplataforma, fáciles de desplegar en Raspberry Pi (ARM64) o servidores x86_64.

### ¿Por qué no usar TensorFlow/GoML directamente?

GopherML cubre los algoritmos clásicos (árboles, bosques, lineales, logística) con API de scikit-learn y cero dependencias. Si necesitas deep learning, este no es el proyecto adecuado.

---

## El proyecto

### ¿Qué es `go-nlp/`?

Un módulo **independiente** (su propio `go.mod`) dentro del repositorio: un puerto concurrente de NLTK (tokenización, POS tagging, parsing, sentimiento) con optimización de tokens para LLMs. Se compila y testea por separado: `cd go-nlp && go build ./...`. Ver [go-nlp/README.md](../go-nlp/README.md).

### ¿Puedo usar GopherML en un bot comercial?

Sí: licencia MIT ([LICENSE](../LICENSE)). Sin garantías — valida siempre las señales en un entorno de paper-trading antes de arriesgar capital.

### ¿Cómo reporto un bug o propongo una función?

Usa las [plantillas de issues](https://github.com/ateto3/gopherml/issues/new/choose). Vulnerabilidades: ver [SECURITY.md](../SECURITY.md).

### ¿Dónde empiezo como contribuidor?

[CONTRIBUTING.md](../CONTRIBUTING.md) — incluye los comandos de verificación (`go vet`, `go test -race`) y las directrices de rendimiento.
