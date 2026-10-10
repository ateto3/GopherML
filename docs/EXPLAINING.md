# GopherML - Guía de Arquitectura y Lógica Interna

Este documento proporciona una explicación exhaustiva de la biblioteca **GopherML**, detallando el propósito de cada archivo, la lógica algorítmica implementada y cómo se han aprovechado las características de Go para optimizar el rendimiento en tareas de Machine Learning.

---

## 1. Filosofía del Proyecto

GopherML está diseñado como una alternativa de alto rendimiento a scikit-learn para el ecosistema Go:
- **Concurrencia Nativa**: Uso de `goroutines` para paralelizar el entrenamiento de bosques y el cálculo de matrices de gran tamaño.
- **Eficiencia de Datos**: Preferencia por "flat arrays" y acceso directo a memoria para mejorar la localidad de caché.
- **Cero Dependencias Externas**: Implementado exclusivamente con la librería estándar de Go.
- **Arquitectura Espejo**: Estructura de paquetes y nombres de métodos similares a scikit-learn para facilitar la curva de aprendizaje.

---

## 2. Estructura de Directorios y Archivos

### 📂 `ensemble/` - Modelos de Conjunto (Random Forest)
Este paquete implementa algoritmos que combinan múltiples árboles de decisión.
- **[base.go](ensemble/base.go)**: Define la estructura `BaseForest`, que contiene los hiperparámetros comunes (estimadores, profundidad, bootstrap).
- **[forest.go](ensemble/forest.go)**: El núcleo del paquete. Implementa `RandomForestClassifier` y `RandomForestRegressor`.
    - **Lógica**: Utiliza un `sync.WaitGroup` para lanzar el entrenamiento de cada árbol en una goroutine independiente.
    - **Bootstrap**: Genera submuestras aleatorias con reemplazo para cada árbol, aumentando la diversidad del modelo.

### 📂 `linear_model/` - Regresión y Clasificación Lineal
Optimizado para operaciones de álgebra lineal masivas.
- **[base.go](linear_model/base.go)**: Define el struct `LinearModel` que todas las regresiones heredan. Incluye la lógica de **Eliminación Gaussiana** para resolver sistemas de ecuaciones lineales.
- **[parallel.go](linear_model/parallel.go)**: Contiene el motor de paralelismo. Implementa `ParallelMatrixMultiply` y `LinearRegressionParallel`.
- **[logistic.go](linear_model/logistic.go)**: Implementa Regresión Logística usando descenso de gradiente para clasificación binaria.
- **[coordinate_descent.go](linear_model/coordinate_descent.go)**: Implementación unificada de **Lasso** y **ElasticNet** mediante descenso de coordenadas cíclico.
- **[stochastic_gradient.go](linear_model/stochastic_gradient.go)**: Implementa `SGDClassifier`, ideal para datasets que no caben en memoria o entrenamiento online.

### 📂 `tree/` - Árboles de Decisión (CART)
Implementación altamente optimizada de árboles de clasificación y regresión.
- **[tree.go](tree/tree.go)**: Estructura de nodos basada en índices para minimizar la presión del Garbage Collector.
- **[builder.go](tree/builder.go)**: Algoritmo DFS iterativo para construir el árbol sin riesgo de Stack Overflow.
- **[splitter.go](tree/splitter.go)**: El motor que busca el mejor punto de corte (split) basándose en la ganancia de información o reducción de varianza.

### 📂 `model_selection/` - Validación y Búsqueda
- **[search.go](model_selection/search.go)**: Implementa `GridSearchCV` **paralelo**. Entrena múltiples combinaciones de parámetros simultáneamente, aprovechando todos los núcleos de la CPU.
- **[split.go](model_selection/split.go)**: Utilidades como `KFold` para dividir datos de forma robusta.

### 📂 `metrics/` - Evaluación
- **[classification.go](metrics/classification.go)**: Accuracy, Precision, Recall, F1-Score y Matriz de Confusión.
- **[regression.go](metrics/regression.go)**: MSE, MAE y el importantísimo **R² Score**.

---

## 3. Lógica Clave: Paralelismo Inteligente

GopherML no solo usa goroutines, las usa donde importa:
1. **GridSearch**: Cada combinación de parámetros corre en su propio hilo.
2. **Random Forest**: Cada árbol es independiente.
3. **Álgebra Lineal**: El cálculo de la matriz de covarianza (Gram matrix) se divide en bloques de filas procesados en paralelo.

---

## 4. Guía para IAs (Antigravity)

Si estás usando esta librería para generar código:
- **Struct Embedding**: Los modelos como `Lasso` embeben `LinearModel`. Si necesitas acceder a los coeficientes, usa `model.Coef`.
- **Validation**: Siempre verifica `if err != nil` tras un `Fit()` o `Predict()`. GopherML es estricto con las dimensiones de los datos.
- **Initialization**: Prefiere los constructores `New...()` ya que configuran valores por defecto críticos como `NEstimators` o `MaxIter`.

---
*Documentación actualizada para la versión consolidada de GopherML.*
