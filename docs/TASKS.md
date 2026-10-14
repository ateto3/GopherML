# 🗂️ Tasks and Development Phases

GopherML technical roadmap. The goal is to offer an experience similar to Python's Machine Learning libraries — especially scikit-learn — while maintaining an idiomatic Go implementation, without Cgo and without external dependencies.

This document separates planned work from available work. A task is only marked as completed when it meets its acceptance criteria, includes tests, and the affected documentation is updated.

## How to read this document

- `[x]` completed and verified.
- `[ ]` pending.
- `[~]` in progress; should include a linked issue or PR when it exists.
- **P0** blocks quality or basic compatibility; **P1** improves the main experience; **P2** is an extension.
- Each phase produces a usable result and can be divided into small PRs.

## Current status summary

| Area | Status | Base available |
| :--- | :---: | :--- |
| Model API | 🟡 | `Fit`, `Predict`, `Score` and single sample APIs in several models |
| Supervised models | 🟡 | linear, trees, Random Forest and SGD |
| Preprocessing | 🟡 | scalers, label encoder and polynomial features |
| Evaluation and selection | 🟡 | metrics, K-Fold, cross-validation and grid search |
| Persistence | 🟢 | `gob` with post-load validation |
| Production and performance | 🟡 | native concurrency, benchmarks and crypto pipeline |
| Documentation and public API | 🟡 | guides, API reference and examples; need to consolidate parity and versions |
| `go-nlp` | 🟡 | separate subproject with own modules and tests |

To understand existing behavior, consult [API_REFERENCE.md](API_REFERENCE.md),
[ARCHITECTURE.md](ARCHITECTURE.md), [PERFORMANCE.md](PERFORMANCE.md) and
[CONTRIBUTING.md](../CONTRIBUTING.md).

---

## Phase 0 — Public contract and baseline quality

**Objective:** establish the rules that all estimators must share before expanding the library surface.

- [x] **P0 — Constructor and error conventions.** Maintain `New...` constructors, explicit errors, and rejection of non-finite inputs.
- [ ] **P0 — API inventory.** Document for each estimator its methods, input/output types, fitted state, and differences from `sklearn`.
- [ ] **P0 — Common contracts.** Define small interfaces for estimators, classifiers, regressors, and transformers without forcing an artificial hierarchy.
- [ ] **P0 — Compatibility policy.** Publish versioning rules, incompatible changes, error stability, and persisted model compatibility.
- [ ] **P1 — Platform matrix.** Verify Go 1.23+, Linux/macOS/Windows, `amd64` and `arm64` in CI.

**Exit criteria**

- All exported APIs have GoDoc comments and example or test when applicable.
- `go build ./...`, `go vet ./...` and `go test ./... -race` pass in CI.
- The API reference doesn't promise methods that the code doesn't expose.

---

## Phase 1 — scikit-learn-style API

**Objective:** make the `constructor → fit → predict → score` flow consistent between packages and easy to discover for those coming from Python.

- [ ] **P0 — Uniform lifecycle.** Align names, "not fitted" errors, dimension validation, and retraining behavior.
- [ ] **P0 — Composable transformers.** Normalize `Fit`, `Transform`, `FitTransform`, and `InverseTransform`; document what state is learned only with training data.
- [ ] **P1 — Pipelines.** Add a typed pipeline to chain transformers and estimator, avoiding information leakage in cross-validation.
- [ ] **P1 — Columns and composition.** Design an idiomatic alternative to `ColumnTransformer` to apply transformations by feature subset.
- [ ] **P1 — Cloning and parameters.** Allow cloning estimators and querying or updating hyperparameters safely for `GridSearchCV`.
- [ ] **P2 — Fit results.** Expose structures equivalent to `cv_results_`, with scores, times, parameters, and best configuration.

**Exit criteria**

- A complete pipeline appears in an executable example and in the usage guide.
- Tests verify that the transformer fits only with the training split.
- Batch and single sample APIs have compatible errors and semantics.

---

## Phase 2 — Preprocessing and data

**Objective:** cover operations that normally precede training.

- [x] **P0 — Basic scaling.** `StandardScaler` and `MinMaxScaler`.
- [x] **P0 — Basic encoding.** `LabelEncoder` and class validation.
- [ ] **P1 — Imputation.** Add missing value imputation with documented strategies (`mean`, `median`, constant) without silently converting `NaN`.
- [ ] **P1 — Categorical encoding.** Incorporate one-hot encoding with explicit control of unknown categories.
- [ ] **P1 — Robust splits.** Add `StratifiedKFold`, `GroupKFold`, and reproducible shuffle options.
- [ ] **P2 — Auxiliary datasets.** Provide small, reproducible loaders for examples, without making them a runtime dependency.

**Exit criteria:** each transformer covers empty data, incorrect dimensions, unknown categories, `NaN`/`Inf`, retraining, and concurrent use after `Fit`.

---

## Phase 3 — Supervised models

**Objective:** expand functional coverage while maintaining deterministic results comparable to reference implementations.

- [x] **P0 — Base regression and classification.** OLS, Ridge, Lasso, Elastic Net, Logistic Regression, and SGD.
- [x] **P0 — Trees and forests.** CART and Random Forest for classification and regression.
- [ ] **P1 — Boosting.** Evaluate Gradient Boosting with a small first version, explainable, and without native dependencies.
- [ ] **P1 — Neighbors and centroids.** Add KNN and K-Means with clear contracts for distances, ties, and seeds.
- [ ] **P1 — Linear SVM.** Prioritize a linear version before kernels to preserve memory footprint and portability.
- [ ] **P2 — Probabilistic models.** Naive Bayes and probability calibration.

**Exit criteria:** each model has safe default parameters, edge case tests, representative benchmark, minimal example, and numerical comparison with a Python reference on small datasets.

---

## Phase 4 — Evaluation, selection, and explainability

**Objective:** make reliable model comparison and result reading.

- [x] **P0 — Main metrics.** Accuracy, precision, recall, F1, ROC-AUC, confusion matrix, MAE, MSE, RMSE, and R².
- [x] **P0 — Validation and search.** K-Fold, cross-validation, and grid search.
- [ ] **P1 — Additional metrics.** Balanced accuracy, log loss, average precision, median absolute error, and well-specified multiclass metrics.
- [ ] **P1 — Curves and thresholds.** ROC/PR as serializable data, threshold selection, and positive class documentation.
- [ ] **P1 — Feature importance.** Tree importance and permutation, with warnings about correlation and information leakage.
- [ ] **P2 — Local explainability.** Design a small API for per-feature contributions without promising full SHAP compatibility.

**Exit criteria:** metrics specify averaging, labels, missing classes, and division by zero; selection results are reproducible with a fixed seed.

---

## Phase 5 — Production, performance, and persistence

**Objective:** make a trained model safe to deploy and measure.

- [x] **P0 — Basic persistence.** Save/load with `gob` and validate reconstructed state.
- [ ] **P0 — Artifact versioning.** Include format version, model type, metadata, and actionable error on incompatibility.
- [ ] **P0 — Load security.** Limit sizes and validate structures before accepting untrusted artifacts; document that `gob` is not a format for unknown origin data.
- [x] **P1 — Hot path.** Maintain `PredictSingle`, fast finite validation, and allocation benchmarks.
- [ ] **P1 — Performance budgets.** Set latency p50/p99, memory, and throughput objectives per model and platform.
- [ ] **P1 — Observability.** Add hooks or optional results for inference time, errors, and model version, without global logging.
- [ ] **P2 — Portable format.** Study JSON/CBOR or ONNX only if there's a real interoperability use case and it can be maintained without Cgo.

**Exit criteria:** each optimization includes before/after benchmark, `go test -race` detects no races, and an incompatible model fails explicitly before predicting.

---

## Phase 6 — Documentation, examples, and distribution

**Objective:** enable a new user to install, train, evaluate, and deploy a model without reading the source code.

- [ ] **P0 — Single documentation.** Maintain correct relative links between README, `docs/`, examples, and API reference.
- [ ] **P0 — Verifiable examples.** Compile all programs in `examples/` in CI and maintain a minimal example per model family.
- [ ] **P1 — Migration guides.** Document frequent equivalences with `sklearn`, especially arrays, classes, errors, and `random_state`.
- [ ] **P1 — Releases.** Add changelog, release notes, support policy, and reproducible publication procedure.
- [ ] **P1 — Public benchmarks.** Publish datasets, hardware, Go version, command, and criteria for interpreting each benchmark.
- [ ] **P2 — Generated reference.** Evaluate automatic generation from GoDoc to reduce divergence with `API_REFERENCE.md`.

**Exit criteria:** a clean checkout can run build, tests, examples, and linked documentation following only [GETTING_STARTED.md](GETTING_STARTED.md).

---

## Phase 7 — `go-nlp` and advanced scope

**Objective:** consolidate the NLP subproject without mixing its contracts with the core Machine Learning.

- [ ] **P1 — Module boundary.** Clearly document its `go.mod`, API, and independent release lifecycle.
- [ ] **P1 — Coverage and benchmarks.** Add regression coverage for parser, tagger, stemmer, translation, and classification.
- [ ] **P2 — Optional integration.** Define explicit adapters to use NLP features from GopherML without introducing circular dependencies.
- [ ] **P2 — Corpus compatibility.** Document supported formats, encoding, and memory limits.

**Exit criteria:** the subproject can be compiled and tested in isolation, and documentation makes clear which APIs are stable.

---

## Definition of Done for any task

Before closing a task or PR:

- [ ] The scope and technical decision are described.
- [ ] There are deterministic tests for normal behavior and relevant errors.
- [ ] Ran `gofmt`, `go vet ./...`, `go test ./... -race`, and if applicable, `go test ./... -bench=. -benchmem`.
- [ ] No dependencies or Cgo were added without an explicit decision.
- [ ] The exported API has English GoDoc and idiomatic Go names.
- [ ] Updated the reference, examples, changelog, or this roadmap when applicable.
- [ ] The PR explains what changes, why, how it was validated, and any difference from the reference Python library.

## Recommended execution order

1. Close Phase 0 to stabilize contracts and CI.
2. Complete Phase 1 before adding many new model families.
3. Continue with Phases 2 and 4 to avoid difficult-to-validate pipelines.
4. Address Phase 5 before declaring models production-ready.
5. Maintain Phase 6 in parallel with each public change.
6. Treat Phase 7 as an independent work line.

Priority may change based on actual user needs, but any change must maintain the project's quality criteria and zero-dependency scope.
