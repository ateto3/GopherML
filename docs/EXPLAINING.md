# GopherML - Architecture and Internal Logic Guide

This document provides a comprehensive explanation of the **GopherML** library, detailing the purpose of each file, the algorithmic logic implemented, and how Go features have been leveraged to optimize performance in Machine Learning tasks.

---

## 1. Project Philosophy

GopherML is designed as a high-performance alternative to scikit-learn for the Go ecosystem:
- **Native Concurrency:** Use of `goroutines` to parallelize forest training and large matrix computation.
- **Data Efficiency:** Preference for "flat arrays" and direct memory access to improve cache locality.
- **Zero External Dependencies:** Implemented exclusively with the Go standard library.
- **Mirror Architecture:** Package structure and method names similar to scikit-learn to facilitate the learning curve.

---

## 2. Directory and File Structure

### 📂 `ensemble/` - Ensemble Models (Random Forest)
This package implements algorithms that combine multiple decision trees.
- **[base.go](ensemble/base.go)**: Defines the `BaseForest` structure, which contains common hyperparameters (estimators, depth, bootstrap).
- **[forest.go](ensemble/forest.go)**: The core of the package. Implements `RandomForestClassifier` and `RandomForestRegressor`.
    - **Logic:** Uses a `sync.WaitGroup` to launch the training of each tree in an independent goroutine.
    - **Bootstrap:** Generates random subsamples with replacement for each tree, increasing model diversity.

### 📂 `linear_model/` - Linear Regression and Classification
Optimized for massive linear algebra operations.
- **[base.go](linear_model/base.go)**: Defines the `LinearModel` struct that all regressions inherit from. Includes the **Gaussian Elimination** logic for solving linear equation systems.
- **[parallel.go](linear_model/parallel.go)**: Contains the parallelism engine. Implements `ParallelMatrixMultiply` and `LinearRegressionParallel`.
- **[logistic.go](linear_model/logistic.go)**: Implements Logistic Regression using gradient descent for binary classification.
- **[coordinate_descent.go](linear_model/coordinate_descent.go)**: Unified implementation of **Lasso** and **ElasticNet** via cyclic coordinate descent.
- **[stochastic_gradient.go](linear_model/stochastic_gradient.go)**: Implements `SGDClassifier`, ideal for datasets that don't fit in memory or online training.

### 📂 `tree/` - Decision Trees (CART)
Highly optimized implementation of classification and regression trees.
- **[tree.go](tree/tree.go)**: Index-based node structure to minimize Garbage Collector pressure.
- **[builder.go](tree/builder.go)**: Iterative DFS algorithm to build the tree without Stack Overflow risk.
- **[splitter.go](tree/splitter.go)**: The engine that searches for the best split point based on information gain or variance reduction.

### 📂 `model_selection/` - Validation and Search
- **[search.go](model_selection/search.go)**: Implements **parallel** `GridSearchCV`. Trains multiple parameter combinations simultaneously, leveraging all CPU cores.
- **[split.go](model_selection/split.go)**: Utilities like `KFold` for robust data splitting.

### 📂 `metrics/` - Evaluation
- **[classification.go](metrics/classification.go)**: Accuracy, Precision, Recall, F1-Score, and Confusion Matrix.
- **[regression.go](metrics/regression.go)**: MSE, MAE, and the important **R² Score**.

---

## 3. Key Logic: Intelligent Parallelism

GopherML doesn't just use goroutines, it uses them where it matters:
1. **GridSearch**: Each parameter combination runs in its own thread.
2. **Random Forest**: Each tree is independent.
3. **Linear Algebra**: The covariance matrix (Gram matrix) calculation is divided into row blocks processed in parallel.

---

## 4. Guide for AIs (Antigravity)

If you're using this library to generate code:
- **Struct Embedding**: Models like `Lasso` embed `LinearModel`. If you need to access coefficients, use `model.Coef`.
- **Validation**: Always verify `if err != nil` after a `Fit()` or `Predict()`. GopherML is strict with data dimensions.
- **Initialization**: Prefer the `New...()` constructors as they configure critical default values like `NEstimators` or `MaxIter`.

---
*Documentation updated for the consolidated version of GopherML.*
