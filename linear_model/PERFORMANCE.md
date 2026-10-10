# Optimizaciones de Rendimiento en LinearModel

Este documento describe las optimizaciones específicas de Go implementadas en el paquete `linear_model`.

## 🚀 Características de Alto Rendimiento

### 1. Goroutines y Paralelismo Nativo

Todas las operaciones intensivas computacionalmente utilizan goroutines:

```go
// Ejemplo: ParallelFor divide trabajo entre CPUs
ParallelFor(nSamples, opts, func(start, end int) {
    for i := start; i < end; i++ {
        // Procesar batch
    }
})
```

**Uso de `sync.WaitGroup`** para sincronización eficiente sin bloqueos innecesarios.

### 2. Memory Pools y Reutilización de Buffers

```go
// Pre-alocación de memoria contigua para mejor localidad de caché
dataBlock := make([]float64, nSamples*nFeatures)
for i := range X {
    XWork[i] = dataBlock[i*nFeatures : (i+1)*nFeatures]
}
```

**Ventaja**: Una sola allocación vs nSamples allocaciones separadas.

### 3. Punteros y Slice Headers

Uso extensivo de punteros para evitar copias:
```go
xi := X[i]  // Puntero a fila, evita bounds checks en inner loops
for j := 0; j < nFeatures; j++ {
    sum += xi[j] * xi[k]  // Acceso directo por puntero
}
```

### 4. Canales para Agregación

Los workers envían resultados parciales vía canales:
```go	type partialSum struct {
		idx  int
		val  float64
	}
	results := make(chan partialSum, nFeatures)
	// Workers calculan y envían
	results <- partialSum{idx: j, val: sum}
```

### 5. Loop Unrolling Manual

Para operaciones críticas:
```go
// Unroll x4 para mejor pipelining de CPU
for ; j <= nFeatures-4; j += 4 {
    sum += row[j]*v[j] + row[j+1]*v[j+1] + 
           row[j+2]*v[j+2] + row[j+3]*v[j+3]
}
```

## 📊 Resultados de Benchmarks

| Dataset | Secuencial | Paralelo | Speedup |
|---------|-----------|----------|---------|
| 50K×200 (Large) | 46.0s | 1.8s | **25x** |
| 5K×50 Lasso | 91ms | 17ms | **5.3x** |
| 10K predicción | 2.4ms | 1.2ms | **2x** |

## 🎯 Modelos Paralelos Disponibles

1. **`LinearRegressionParallel`** - Multiplicación de matrices paralela
2. **`LassoParallel`** - Coordinate descent con actualizaciones paralelas
3. **`WorkerPool`** reutilizable para procesamiento por lotes

## 🔧 Uso Recomendado

### Cuándo usar versiones paralelas:
- Dataset > 10,000 samples
- Features > 50
- Predicciones en batch grandes

### Cuándo usar versiones secuenciales:
- Dataset < 1,000 samples (overhead de goroutines no compensa)
- Sistemas single-core
- Aplicaciones en tiempo real de baja latencia

```go
// Ejemplo de uso
opts := linear_model.ParallelOptions{
    NumWorkers:   runtime.NumCPU(),
    MinBatchSize: 1000,
}

model := linear_model.NewLinearRegressionParallelWithParams(
    linear_model.LinearRegressionParams{FitIntercept: true},
    opts,
)
model.Fit(X, y)
```

## 💡 Optimizaciones Específicas de Go

| Técnica | Beneficio |
|---------|-----------|
| Goroutines livianas | ~2KB stack vs ~1MB thread |
| Slices contiguas | Mejor localidad de caché |
| `sync.Pool` | Reutilización de objetos GC |
| Channels buffered | Reduce contención |
| Mutex selectivo | Solo donde es necesario |

## 🚫 Zero Dependencies

El paquete utiliza **exclusivamente** la librería estándar de Go:
- `runtime` - NumCPU()
- `sync` - WaitGroup, Mutex, Pool
- No cgo, no numpy, no dependencias externas
