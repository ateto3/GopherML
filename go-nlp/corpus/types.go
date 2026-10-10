package corpus

import (
	"io/fs"
	"sync"
)

// Reader defines the interface for corpus readers.
type Reader interface {
	FileIDs() []string
	Words(fileIDs ...string) ([]string, error)
	Raw(fileIDs ...string) (string, error)
	Size() (int, error)
	Vocab() (map[string]bool, error)
}

// FileStats holds statistics for a single corpus file.
type FileStats struct {
	ID            string
	WordCount     int
	CharCount     int
	SentenceCount int
}

// CorpusStats holds overall statistics for a corpus.
type CorpusStats struct {
	FileCount      int
	TotalWords     int
	TotalChars     int
	TotalSentences int
	VocabSize      int
	Files          []FileStats
}

// ReaderConfig holds configuration for corpus readers.
type ReaderConfig struct {
	SplitPattern      string
	Lowercase         bool
	RemovePunctuation bool
	MinWordLength     int
	MaxWordLength     int
	Language          string
}

// NewConfig creates a new ReaderConfig with default values.
func NewConfig() *ReaderConfig {
	return &ReaderConfig{
		Language:        "en",
		MinWordLength:   1,
		MaxWordLength:   100,
	}
}

// StatsReader wraps a Reader and computes statistics.
type StatsReader struct {
	Reader
	config *ReaderConfig
}

// NewStatsReader creates a StatsReader.
func NewStatsReader(r Reader, config *ReaderConfig) *StatsReader {
	return &StatsReader{Reader: r, config: config}
}

// FSReader reads corpus files from an fs.FS with optimizations.
type FSReader struct {
	FS           fs.FS
	Files        []string
	SplitFn      func(string) []string
	config       *ReaderConfig
	rawCache     map[string]string
	wordsCache   map[string][]string
	loadedFiles  map[string]bool
	loadedMutex  sync.RWMutex
}
