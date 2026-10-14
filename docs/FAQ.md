Frequently Answered Questions (FAQs)

## Errors and Validation

### `Fit`/`Predict` returns` NaN `or` Inf `— why?

GopherML explicitly rejects non-finite values in the input. The error indicates the position of the problem. Typical solution: review the calculation of features (divisions by zero, `log(0)`) or clean the CSV before training. Never 'patch' with `math.IsNaN` out: validation is intentional.

### He tells me `dimension mismatch` — what do I check?

`len(X)` (samples) must match `len(y)`, and **all rows** of `X` must have the same number of columns. For prediction, the number of columns must be the same as in training.

### `Transform` returns error before `Fit`

Correct: Scalers need to know mean/deviation (or min/max) with `Fit` (or `FitTransform`) before transforming. It never fails in silence.

### `model not fitted` / `tree is not fitted yet`

You are calling `Predict` before `Fit`, or you uploaded an incomplete model. In `persistence.Load`, models with the `Validator` interface detect this automatically.

---

## API and usage

### What is the difference between `Predict` and `PredictSingle`?

`Predict` processes a batch (`[][]float64`) and reserves memory for the result; `PredictSingle` processes a sample with the minimum of allocations — use it in trading loops (see [PERFORMANCE.md](PERFORMANCE.md)).

### Why does` tree.DecisionTreeClassifier.Fit (X, y, sampleWeight) `have 3 arguments?

Low-level trees support weights per sample. That's why they * * don't * * implement the `Estimator` interface of `model_selection`; for cross-validation use `ensemble.RandomForest *`, whose signature is` Fit(X, y) `.

### How do I get probabilities to filter by trust?

```go
proba, _ := rf.PredictProbaSingle(featureVector)
if proba[1] >= 0.65 { /* confianza suficiente */ }
```

In `crypto`, this filter is already embedded: `PredictSignal` returns `SignalHold` if confidence doesn't exceed `MinConfidence`.

### How do I make results reproducible?

Set `RandomState` (and `Shuffle: true` in `KFold`). GopherML uses seeded RNG (`utils.NewRNG`) instead of `math/rand` global state.

### What does `Score` return?

Accuracy for classifiers and R² for regressors (equivalent to scikit-learn).

---

## Persistence

### Are `.gob` files compatible between versions?

`gob` serializes the concrete struct structure: if you change the GopherML version and model internal fields change, retrain or save the version as well. In production, version model files together with the binary.

### Can I save an untrained model?

Yes it would serialize, but `Load` will fail if the model implements `Validator` (which checks `IsFitted`). That's the integrity protection.

### Is the scaler saved with the model?

With `CryptoPredictor.SaveModel`, yes (model + scaler together). For standalone models, save the scaler separately or fit `StandardScaler` with the same training data.

---

## Performance

### Should I use parallel versions?

Only if the dataset is large (>10,000 samples or >50 features). For small datasets, goroutines add overhead. See [PERFORMANCE.md](PERFORMANCE.md).

### Does GopherML use Cgo or GPU?

No. Only Go standard library: static binaries, cross-platform, easy to deploy on Raspberry Pi (ARM64) or x86_64 servers.

### Why not use TensorFlow/GoML directly?

GopherML covers classical algorithms (trees, forests, linear, logistic) with scikit-learn API and zero dependencies. If you need deep learning, this is not the right project.

---

## The project

### What is `go-nlp/`?

An **independent** module (its own `go.mod`) within the repository: a concurrent NLTK port (tokenization, POS tagging, parsing, sentiment) with token optimization for LLMs. It compiles and tests separately: `cd go-nlp && go build ./...`. See [go-nlp/README.md](../go-nlp/README.md).

### Can I use GopherML in a commercial bot?

Yes: MIT license ([LICENSE](../LICENSE)). No warranties — always validate signals in a paper-trading environment before risking capital.

### How do I report a bug or propose a feature?

Use the [issue templates](https://github.com/ateto3/gopherml/issues/new/choose). Vulnerabilities: see [SECURITY.md](../SECURITY.md).

### Where do I start as a contributor?

[CONTRIBUTING.md](../CONTRIBUTING.md) — includes verification commands (`go vet`, `go test -race`) and performance guidelines.
