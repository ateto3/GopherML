# 🤝 Contributing to GopherML

Thank you for your interest in improving GopherML! This document collects the workflow, code standards and performance guidelines we follow. By submitting a pull request you also accept our [Code of Conduct](CODE_OF_CONDUCT.md).

---

## 1. Getting started

```bash
# 1. Fork the repository and clone it
git clone https://github.com/YOUR-USERNAME/gopherml.git
cd gopherml

# 2. Verify that everything compiles and the tests pass
go build ./...
go vet ./...
go test ./... -race
```

**Requirements:** Go 1.23 or later. The project has no external dependencies, so no extra `go mod download` is needed.

> **Note:** `go-nlp/` is a standalone module with its own `go.mod`. If you touch that code,
> also run `cd go-nlp && go build ./... && go test ./...`.

---

## 2. Code standards

1. **Idiomatic Go:** no Python or Java patterns. Small interfaces, short names, zero unnecessary abstraction.
2. **Explicit error handling:** never ignore the error returned by `Fit`, `Predict`, `Load` or `Save`. Errors are wrapped with context (`fmt.Errorf("...: %w", err)`).
3. **Constructors:** always use the `New...()` constructors; they set critical default hyperparameters (`NEstimators`, `MaxIter`, etc.).
4. **Formatting:** code must pass `gofmt -w .` and `go vet ./...` without warnings.
5. **Documentation of exports:** every exported function, type or method has a doc comment in English that starts with its name (Go convention).

---

## 3. Performance guidelines (critical)

This is **not** a toy project: we aim for near-C++ efficiency in low-latency inference.

1. **Performance priority:** prefer pre-allocated *slices* and avoid memory allocations on the *hot path*. Use `PredictSingle` instead of `Predict` inside high-frequency loops.
2. **Native concurrency:** use `sync.WaitGroup` and channels only when the parallelism pays for the cost of communication. We optimize for multi-core CPUs (ARM64 and x86_64).
3. **Zero dependencies:** the core depends only on the Go standard library. Do not add external dependencies or Cgo.
4. **Data safety:** all input data goes through `utils.CheckFinite` (or `CheckFiniteSingle` on the hot path) before reaching any model.
5. **No silent regressions:** if you optimize, include or update a benchmark in `*_test.go` that proves it.

---

## 4. Guidelines for AI assistants (Windsurf / Antigravity / Copilot…)

If you are generating code with an AI for this repository, follow these rules without exception:

- **`PredictSingle` in loops**, never `Predict` with a single sample.
- **Always validate `if err != nil`** after `Fit()`, `Predict()`, `Load()`.
- **Struct embedding:** models like `Lasso` embed `LinearModel`; coefficients are accessed as `model.Coef`.
- **`New...()` constructors**, never uninitialized struct literals.
- **Zero new dependencies** in the main module's `go.mod`.
- The generated code must compile with `go vet ./...` and pass `go test ./... -race` **before** opening the PR.

---

## 5. Tests

- Every behavior change includes tests.
- Run the full suite before submitting: `go test ./... -race`.
- Tests must be deterministic: use `utils.NewRNG(seed)` instead of the global `rand`.
- For performance changes, add a benchmark (`Benchmark*`) in the affected package.

---

## 6. Pull request workflow

1. Create a descriptive branch: `feature/my-feature` or `fix/my-fix`.
2. One goal per PR (mixing features with fixes makes review harder).
3. Update the affected documentation (`README.md`, `USAGE.md`, `docs/`, `CHANGELOG.md`).
4. Check the [PR template](.github/PULL_REQUEST_TEMPLATE.md) checklist.
5. Describe **what** changes, **why** it changes and **how** it was validated.

### Commit convention (recommended)

```
feat(linear_model): add alternative solver for Ridge
fix(crypto): validate candles with zero volume
docs(api): update model_selection reference
test(ensemble): cover PredictProbaSingle
```

---

## 7. Reporting issues

- **Bugs:** use the [bug report template](https://github.com/ateto3/gopherml/issues/new?template=bug_report.yml).
- **Features:** use the [feature request template](https://github.com/ateto3/gopherml/issues/new?template=feature_request.yml).
- **Security vulnerabilities:** do not open a public issue; follow [SECURITY.md](SECURITY.md).

---

## 8. License

By contributing, you grant your contributions under the project's [MIT license](LICENSE).
