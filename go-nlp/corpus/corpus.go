package corpus

import (
	"io"
	"io/fs"
	"strings"
	"sync"
)

// NewFSReader creates an optimized corpus reader from fs.FS.
func NewFSReader(fsys fs.FS, files []string, opts ...func(*ReaderConfig)) *FSReader {
	config := NewConfig()
	for _, opt := range opts {
		opt(config)
	}

	reader := &FSReader{
		FS:          fsys,
		Files:       files,
		SplitFn:     strings.Fields,
		config:      config,
		rawCache:    make(map[string]string),
		wordsCache:  make(map[string][]string),
		loadedFiles: make(map[string]bool),
	}

	if config.SplitPattern != "" {
		reader.SplitFn = func(text string) []string {
			return strings.FieldsFunc(text, func(r rune) bool {
				return strings.ContainsRune(config.SplitPattern, r)
			})
		}
	}

	return reader
}

// NewLowercase creates an option to convert text to lowercase.
func NewLowercase() func(*ReaderConfig) {
	return func(c *ReaderConfig) { c.Lowercase = true }
}

// NewRemovePunctuation creates an option to remove punctuation.
func NewRemovePunctuation() func(*ReaderConfig) {
	return func(c *ReaderConfig) { c.RemovePunctuation = true }
}

// NewMinWordLength creates an option to set minimum word length.
func NewMinWordLength(min int) func(*ReaderConfig) {
	return func(c *ReaderConfig) { c.MinWordLength = min }
}

// NewMaxWordLength creates an option to set maximum word length.
func NewMaxWordLength(max int) func(*ReaderConfig) {
	return func(c *ReaderConfig) { c.MaxWordLength = max }
}

// FileIDs returns all available file IDs.
func (r *FSReader) FileIDs() []string {
	return append([]string(nil), r.Files...)
}

// Raw returns raw text content with caching.
func (r *FSReader) Raw(fileIDs ...string) (string, error) {
	if len(fileIDs) == 0 {
		fileIDs = r.Files
	}

	var builder strings.Builder
	for _, id := range fileIDs {
		if cached, ok := r.getRawCached(id); ok {
			builder.WriteString(cached)
		} else {
			f, err := r.FS.Open(id)
			if err != nil {
				return "", err
			}
			data, err := io.ReadAll(f)
			f.Close()
			if err != nil {
				return "", err
			}
			raw := string(data)
			r.cacheRaw(id, raw)
			builder.WriteString(raw)
		}
	}
	return builder.String(), nil
}

// RawFile returns raw text for a single file.
func (r *FSReader) RawFile(fileID string) (string, error) {
	if cached, ok := r.getRawCached(fileID); ok {
		return cached, nil
	}

	f, err := r.FS.Open(fileID)
	if err != nil {
		return "", err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	raw := string(data)
	r.cacheRaw(fileID, raw)
	return raw, nil
}

// Words returns words from files with caching and filtering.
func (r *FSReader) Words(fileIDs ...string) ([]string, error) {
	if len(fileIDs) == 0 {
		fileIDs = r.Files
	}

	var allWords []string
	for _, id := range fileIDs {
		words, err := r.WordsFile(id)
		if err != nil {
			return nil, err
		}
		allWords = append(allWords, words...)
	}
	return allWords, nil
}

// WordsFile returns words from a single file with caching.
func (r *FSReader) WordsFile(fileID string) ([]string, error) {
	if cached, ok := r.getWordsCached(fileID); ok {
		return cached, nil
	}

	raw, err := r.RawFile(fileID)
	if err != nil {
		return nil, err
	}

	words := r.processWords(raw)
	r.cacheWords(fileID, words)
	return words, nil
}

// processWords applies all configured processing to text.
func (r *FSReader) processWords(text string) []string {
	if r.config.Lowercase {
		text = strings.ToLower(text)
	}

	tokens := r.SplitFn(text)

	var result []string
	for _, tok := range tokens {
		if r.config.RemovePunctuation {
			tok = strings.Map(func(r rune) rune {
				if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' {
					return r
				}
				return -1
			}, tok)
		}

		if len(tok) >= r.config.MinWordLength && len(tok) <= r.config.MaxWordLength {
			if len(tok) > 0 {
				result = append(result, tok)
			}
		}
	}

	return result
}

// Size returns total word count.
func (r *FSReader) Size() (int, error) {
	words, err := r.Words()
	if err != nil {
		return 0, err
	}
	return len(words), nil
}

// Vocab returns unique words in the corpus.
func (r *FSReader) Vocab() (map[string]bool, error) {
	words, err := r.Words()
	if err != nil {
		return nil, err
	}

	vocab := make(map[string]bool)
	for _, w := range words {
		vocab[w] = true
	}
	return vocab, nil
}

// Caching helpers
func (r *FSReader) getRawCached(fileID string) (string, bool) {
	r.loadedMutex.RLock()
	defer r.loadedMutex.RUnlock()
	val, ok := r.rawCache[fileID]
	return val, ok
}

func (r *FSReader) cacheRaw(fileID, raw string) {
	r.loadedMutex.Lock()
	defer r.loadedMutex.Unlock()
	r.rawCache[fileID] = raw
	r.loadedFiles[fileID] = true
}

func (r *FSReader) getWordsCached(fileID string) ([]string, bool) {
	r.loadedMutex.RLock()
	defer r.loadedMutex.RUnlock()
	val, ok := r.wordsCache[fileID]
	return val, ok
}

func (r *FSReader) cacheWords(fileID string, words []string) {
	r.loadedMutex.Lock()
	defer r.loadedMutex.Unlock()
	r.wordsCache[fileID] = words
}

// Lazy wraps a reader and loads it once.
type Lazy struct {
	once   sync.Once
	reader Reader
	load   func() (Reader, error)
	err    error
}

// NewLazy creates a lazy corpus handle.
func NewLazy(load func() (Reader, error)) *Lazy {
	return &Lazy{load: load}
}

func (l *Lazy) ensure() error {
	l.once.Do(func() { l.reader, l.err = l.load() })
	return l.err
}

func (l *Lazy) FileIDs() []string {
	if l.ensure() != nil {
		return nil
	}
	return l.reader.FileIDs()
}

func (l *Lazy) Words(fileIDs ...string) ([]string, error) {
	if err := l.ensure(); err != nil {
		return nil, err
	}
	return l.reader.Words(fileIDs...)
}

func (l *Lazy) Raw(fileIDs ...string) (string, error) {
	if err := l.ensure(); err != nil {
		return "", err
	}
	return l.reader.Raw(fileIDs...)
}
