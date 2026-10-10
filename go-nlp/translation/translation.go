package translation

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
)

// Config configures a local translation engine.
//
// ModelPath should point to either:
// - a directory containing a translated model bundle, or
// - a direct .onnx/.gguf artifact, depending on the backend you wire in.
//
// Threads is the preferred number of intra-op / decode worker threads.
type Config struct {
	ModelPath string
	Threads   int
}

// Engine is the public API expected by the gateway.
type Engine interface {
	ToEnglish(text string, sourceLang string) (string, error)
	FromEnglish(text string, targetLang string) (string, error)
	Close() error
}

// Backend is the internal translation primitive.
// A production ONNX or GGUF implementation can satisfy this interface.
type Backend interface {
	Translate(text string, sourceLang string, targetLang string) (string, error)
	Close() error
}

type engine struct {
	mu      sync.RWMutex
	backend Backend
	closed  bool
	cfg     Config
}

// NewEngine creates a thread-safe engine and loads the model once.
//
// The current implementation validates the model path and wires the engine
// around an internal backend abstraction so the repo can host an ONNX or GGUF
// implementation without changing the public API.
func NewEngine(cfg Config) (Engine, error) {
	if strings.TrimSpace(cfg.ModelPath) == "" {
		return nil, errors.New("translation: ModelPath is required")
	}
	if cfg.Threads < 0 {
		return nil, errors.New("translation: Threads must be >= 0")
	}

	backend, err := newBackend(cfg)
	if err != nil {
		return nil, err
	}

	return &engine{
		backend: backend,
		cfg:     cfg,
	}, nil
}

func (e *engine) ToEnglish(text string, sourceLang string) (string, error) {
	return e.translate(text, sourceLang, "en")
}

func (e *engine) FromEnglish(text string, targetLang string) (string, error) {
	return e.translate(text, "en", targetLang)
}

func (e *engine) translate(text string, sourceLang string, targetLang string) (string, error) {
	sourceLang = normalizeLang(sourceLang)
	targetLang = normalizeLang(targetLang)

	if text == "" {
		return "", nil
	}
	if sourceLang == "" || targetLang == "" {
		return "", errors.New("translation: source and target languages are required")
	}
	if sourceLang == targetLang {
		return text, nil
	}

	e.mu.RLock()
	closed := e.closed
	backend := e.backend
	e.mu.RUnlock()
	if closed {
		return "", errors.New("translation: engine is closed")
	}
	if backend == nil {
		return "", errors.New("translation: backend not initialized")
	}

	return backend.Translate(text, sourceLang, targetLang)
}

func (e *engine) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed {
		return nil
	}
	e.closed = true
	if e.backend == nil {
		return nil
	}
	return e.backend.Close()
}

// normalizeLang maps common aliases to compact ISO-style codes.
func normalizeLang(lang string) string {
	lang = strings.TrimSpace(strings.ToLower(lang))
	switch lang {
	case "eng", "english", "en-us", "en-gb":
		return "en"
	case "spa", "espanol", "español", "spanish":
		return "es"
	case "hin", "hindi":
		return "hi"
	case "arb", "ara", "arabic":
		return "ar"
	case "jpn", "japanese":
		return "ja"
	case "kor", "korean":
		return "ko"
	default:
		return lang
	}
}

// modelArtifactKind returns a coarse hint based on the configured path.
func modelArtifactKind(modelPath string) string {
	switch strings.ToLower(filepath.Ext(modelPath)) {
	case ".onnx":
		return "onnx"
	case ".gguf":
		return "gguf"
	default:
		return "directory"
	}
}

func languagePairKey(sourceLang, targetLang string) string {
	return fmt.Sprintf("%s:%s", normalizeLang(sourceLang), normalizeLang(targetLang))
}
