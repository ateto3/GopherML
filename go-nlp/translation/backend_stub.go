package translation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type unsupportedBackend struct {
	modelPath string
	kind      string
}

func newBackend(cfg Config) (Backend, error) {
	kind := modelArtifactKind(cfg.ModelPath)

	if stat, err := os.Stat(cfg.ModelPath); err == nil && stat.IsDir() {
		kind = "directory"
	}

	// Validate that the artifact exists early so deployment failures are obvious.
	if _, err := os.Stat(cfg.ModelPath); err != nil {
		return nil, fmt.Errorf("translation: model path not accessible: %w", err)
	}

	// This repository ships the engine abstraction and production-ready
	// language handling, but the actual ONNX/GGUF runtime adapter should be
	// selected at build time in a follow-up patch when the exact exported model
	// format is decided. Keeping this explicit avoids a silent stub in prod.
	return &unsupportedBackend{
		modelPath: cfg.ModelPath,
		kind:      kind,
	}, nil
}

func (b *unsupportedBackend) Translate(text string, sourceLang string, targetLang string) (string, error) {
	return "", fmt.Errorf(
		"translation: no runtime backend wired for %s artifact at %q (requested %s -> %s)",
		b.kind,
		b.modelPath,
		strings.ToLower(sourceLang),
		strings.ToLower(targetLang),
	)
}

func (b *unsupportedBackend) Close() error { return nil }

func modelBundlePath(root string, parts ...string) string {
	elems := append([]string{root}, parts...)
	return filepath.Join(elems...)
}
