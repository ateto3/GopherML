# ⚡ Guía de rendimiento

GopherML está diseñado para inferencia de baja latencia en producción (bots de trading, microservicios). Este documento recoge las optimizaciones disponibles y **cómo usarlas correctamente**.

---

## 1. Regla de oro: el hot path

En un bucle de alta frecuencia (un tick por vela, miles de predicciones/segundo):

```go
// ❌ Evita: reserva memoria para un batch entero cada llamada
preds, _ := rf.Predict([][]float64{muestra})

// ✅ Correcto: API de muestra única, sin allocations extra
label, _ := rf.PredictSingle(muestra)
proba, _ := rf.PredictProbaSingle(muestra)

// ✅ En el pipeline de cripto
signal, conf, err := p.PredictSignal(lastCandles)
```

| API | Uso recomendado |
| :--- | :--- |
| `PredictSingle` / `PredictProbaSingle` | Bucles de trading, una muestra por tick |
| `Predict` / `PredictProba` | Evaluación por lotes, backtesting, scoring de datasets |
| `utils.CheckFiniteSingle` | Validación de una muestra en hot path |

---

## 2. Paralelismo: cuándo sí y cuándo no

| Situación | Recomendación |
| :--- | :--- |
| Dataset > 10 000 muestras o > 50 features | Versiones paralelas (`LinearRegressionParallel`, `LassoParallel`) |
| `GridSearchCV` con muchas combinaciones | `NJobs: -1` (todos los núcleos) |
| Entrenamiento de Random Forest | Ya es paralelo de fábrica (una goroutine por árbol) |
| Dataset < 1 000 muestras | Secuencial: el overhead de goroutines no compensa |
| Single-core (p. ej. contenedor con 1 CPU) | Secuencial |
| Predicción de una muestra | **Nunca** paralelizar: `PredictSingle` |

```go
opts := linear_model.ParallelOptions{
    NumWorkers:   runtime.NumCPU(),
    MinBatchSize: 1000, // no partir lotes pequeños
}
model := linear_model.NewLinearRegressionParallel()
model.Fit(X, y)
```

---

## 3. Optimizaciones internas de GopherML

| Técnica | Dónde | Beneficio |
| :--- | :--- | :--- |
| Goroutines + `sync.WaitGroup` | `ensemble`, `model_selection`, `linear_model` | Aprovecha todos los núcleos |
| Memoria contigua pre-asociada | `linear_model` (bloques planos) | Una sola allocación vs N; mejor localidad de caché |
| Punteros a fila en inner loops | `linear_model` | Evita *bounds checks* en bucles calientes |
| `sync.Pool` / buffers reutilizables | `linear_model`, `crypto` | Menos presión sobre el GC |
| *Loop unrolling* (x4) | multiplicaciones críticas | Mejor *pipelining* de CPU |
| Nodos indexados (`[]Node` + índices) | `tree` | Menos allocations que grafos de punteros |
| DFS iterativo con pila explícita | `tree.DepthFirstTreeBuilder` | Cero riesgo de stack overflow |
| `[]byte` / `strings.Builder` | `go-nlp` | *Zero-copy* y menos GC |

Benchmarks detallados de `linear_model`: **[../linear_model/PERFORMANCE.md](../linear_model/PERFORMANCE.md)** (p. ej. 50K×200: 46.0s → 1.8s secuencial vs paralelo).

---

## 4. Cómo medir tú mismo

```bash
# Benchmarks de un paquete
go test ./linear_model -bench=. -benchmem

# Tests + benchmark con profiling de CPU
go test ./ensemble -bench=. -cpuprofile=cpu.out
go tool pprof cpu.out

# Detector de condiciones de carrera (imprescindible antes de producción)
go test ./... -race
```

Indicadores a vigilar:

- **`allocs/op`** en `PredictSingle`: debe ser 0 o casi 0.
- **`B/op`**: allocations por operación; si sube tras un cambio, hay regresión.
- Tiempo **p99** de `PredictSignal` en tu hardware objetivo (Raspberry Pi ARM64, x86_64).

---

## 5. Checklist de producción para bots

| Práctica | Motivo |
| :--- | :--- |
| Cargar el modelo **al arrancar**, no en cada predicción | `Load` es costoso: hazlo una vez |
| Usar `MinConfidence >= 0.60` | Reduce señales falsas en mercados laterales |
| Validar velas antes de predecir | Exchanges envían velas con *gaps* o volumen 0 |
| Guardar scaler junto al modelo | `CryptoPredictor.SaveModel` lo hace automáticamente |
| `RandomState` fijo en entrenamiento | Resultados reproducibles |
| Entrenar en despliegues ARM64 si despliegas en ARM | Evita diferencias de FP entre arquitecturas |

---

## 6. Compilación recomendada

```bash
# Binario optimizado para producción (sin información de VCS si aún no hay git)
go build -ldflags="-s -w" -trimpath -o bot ./cmd
```

GopherML no usa Cgo, así que el binario es estático y multiplataforma por defecto (`CGO_ENABLED=0`).
