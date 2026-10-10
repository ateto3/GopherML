# Translation Local Engine

This document describes the translation module added to `go-nlp` so another
project can integrate it safely.

## What Was Added

We introduced a new package:

- `translation`

It exposes the public API requested for a low-latency gateway:

```go
package translation

type Config struct {
    ModelPath string
    Threads   int
}

type Engine interface {
    ToEnglish(text string, sourceLang string) (string, error)
    FromEnglish(text string, targetLang string) (string, error)
    Close() error
}

func NewEngine(cfg Config) (Engine, error)
```

## Current Design

The new package is structured in two layers:

- Public API and concurrency control in `translation/translation.go`
- Internal runtime backend abstraction in `translation/backend_stub.go`

The engine:

- validates `ModelPath`
- normalizes language codes and aliases
- is safe for concurrent use
- ensures the model is loaded once and reused
- supports graceful `Close()`

## Language Support

The engine normalizes these language aliases:

- English: `en`, `eng`, `english`, `en-us`, `en-gb`
- Spanish: `es`, `spa`, `español`, `espanol`, `spanish`
- Hindi: `hi`, `hin`, `hindi`
- Arabic: `ar`, `ara`, `arb`, `arabic`
- Japanese: `ja`, `jpn`, `japanese`
- Korean: `ko`, `kor`, `korean`

This is enough for the gateway use case where prompts are translated between
English and the main token-expensive languages.

## Why a New Package

The repository already had a package named `translate`, but it is for classic
alignment and IBM Model 1 experimentation. That code is unrelated to runtime
machine translation inference, so the new module was added as `translation` to
avoid breaking existing behavior.

## Files Added

- `translation/translation.go`
- `translation/backend_stub.go`
- `docs/translation-local-engine.md`

## Implementation Notes

The current commit does not yet embed a concrete ONNX or GGUF runtime adapter.
Instead, it provides a clean contract so the production backend can be added
without changing the public API.

This was intentional because:

- the gateway needs a stable interface now
- the exact exported model bundle still needs to be selected
- ONNX encoder-decoder translation models require a tokenizer plus runtime
  wiring that is best added as a focused follow-up patch

## Recommended Model Direction

For a single multilingual model with good coverage across the requested
languages, the recommended starting point is:

- `facebook/nllb-200-distilled-600M`

Why:

- broad language coverage
- strong quality for the target set
- suitable for ONNX export

If the deployment favors absolute minimum latency per pair over breadth, then
pair-specific MarianMT / OPUS-MT tiny models can be considered later.

## Recommended Deployment Shape

Suggested filesystem layout:

```text
models/
  nllb-200-distilled-600M-onnx/
    config.json
    tokenizer.json
    sentencepiece.bpe.model
    encoder_model.onnx
    decoder_model.onnx
    decoder_with_past_model.onnx
```

Suggested loading pattern:

```go
eng, err := translation.NewEngine(translation.Config{
    ModelPath: "models/nllb-200-distilled-600M-onnx",
    Threads:   4,
})
if err != nil {
    return err
}
defer eng.Close()
```

## Status

The repository now has:

- the public API required by the gateway
- normalization for the target languages
- thread-safe lifecycle management
- a placeholder backend contract ready for ONNX integration

The next step is to wire the actual ONNX runtime implementation into the
backend layer.
