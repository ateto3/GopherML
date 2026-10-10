package tokenizer

import (
	"strings"
	"sync"
)

// ParallelTokenizer performs tokenization concurrently by splitting
// the input text into chunks and processing each chunk in a goroutine.
type ParallelTokenizer struct {
	chunkSize     int
	numWorkers    int
	baseTokenizer TokenizerI
}

// NewParallelTokenizer creates a new ParallelTokenizer.
// chunkSize: number of bytes per chunk for parallel processing.
// numWorkers: number of goroutines to use (0 = GOMAXPROCS).
// baseTokenizer: the underlying tokenizer to use for each chunk.
func NewParallelTokenizer(chunkSize, numWorkers int, baseTokenizer TokenizerI) *ParallelTokenizer {
	if chunkSize <= 0 {
		chunkSize = 1024 * 64 // 64KB default
	}
	return &ParallelTokenizer{
		chunkSize:     chunkSize,
		numWorkers:    numWorkers,
		baseTokenizer: baseTokenizer,
	}
}

// Tokenize tokenizes the text in parallel.
// The text is split into chunks, each tokenized concurrently,
// and results are merged in order.
func (pt *ParallelTokenizer) Tokenize(text string) []string {
	if len(text) <= pt.chunkSize {
		return pt.baseTokenizer.Tokenize(text)
	}

	chunks := splitIntoChunks(text, pt.chunkSize)
	if len(chunks) == 0 {
		return nil
	}

	// Determine number of workers
	numWorkers := pt.numWorkers
	if numWorkers <= 0 {
		numWorkers = defaultNumWorkers()
	}
	if numWorkers > len(chunks) {
		numWorkers = len(chunks)
	}

	// Create work channels
	type chunkResult struct {
		index  int
		tokens []string
	}

	chunkChan := make(chan struct {
		index int
		text  string
	}, len(chunks))
	resultChan := make(chan chunkResult, len(chunks))
	var wg sync.WaitGroup

	// Send chunks to channel
	for i, chunk := range chunks {
		chunkChan <- struct {
			index int
			text  string
		}{i, chunk}
	}
	close(chunkChan)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for chunk := range chunkChan {
				tokens := pt.baseTokenizer.Tokenize(chunk.text)
				resultChan <- chunkResult{chunk.index, tokens}
			}
		}()
	}

	// Wait for all workers and close results
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results in order
	results := make([][]string, len(chunks))
	for result := range resultChan {
		results[result.index] = result.tokens
	}

	// Merge results
	var totalTokens int
	for _, r := range results {
		totalTokens += len(r)
	}

	merged := make([]string, 0, totalTokens)
	for _, r := range results {
		merged = append(merged, r...)
	}
	return merged
}

// SpanTokenize returns token spans. For parallel processing, this uses
// single-threaded processing since spans must map to the original text.
func (pt *ParallelTokenizer) SpanTokenize(text string) []Span {
	// Span tokenization must be serial for accurate offsets
	return pt.baseTokenizer.SpanTokenize(text)
}

// ParallelTokenize is a convenience function that tokenizes text in parallel
// using the NLTKWordTokenizer (the default NLTK word tokenizer).
func ParallelTokenize(text string) []string {
	tok := NewParallelTokenizer(0, 0, NewNLTKWordTokenizer())
	return tok.Tokenize(text)
}

// ParallelSentTokenize tokenizes text into sentences using parallel processing.
func ParallelSentTokenize(text string, language string) []string {
	// For parallel sentence tokenization, we use Punkt-like rules
	// This is a simplified version; a full Punkt implementation would go here
	tok := NewParallelTokenizer(0, 0, NewSentTokenizer(language))
	return tok.Tokenize(text)
}

// SentTokenizer is a simple sentence tokenizer based on punctuation rules.
// A more complete implementation would use the Punkt algorithm.
type SentTokenizer struct {
	language string
}

// NewSentTokenizer creates a sentence tokenizer for the given language.
func NewSentTokenizer(language string) *SentTokenizer {
	return &SentTokenizer{language: language}
}

// Tokenize splits text into sentences.
func (st *SentTokenizer) Tokenize(text string) []string {
	// Simple sentence splitting based on common sentence delimiters
	// A full implementation would use PunktSentenceTokenizer
	return splitSentences(text)
}

// SpanTokenize returns sentence spans.
func (st *SentTokenizer) SpanTokenize(text string) []Span {
	// For now, return empty (would need proper implementation)
	sentences := st.Tokenize(text)
	return alignTokens(sentences, text)
}

// splitIntoChunks splits text into chunks at word boundaries.
func splitIntoChunks(text string, chunkSize int) []string {
	if len(text) <= chunkSize {
		return []string{text}
	}

	var chunks []string
	start := 0
	for start < len(text) {
		end := start + chunkSize
		if end >= len(text) {
			chunks = append(chunks, text[start:])
			break
		}

		// Try to break at a space or newline
		breakAt := strings.LastIndexAny(text[start:end], " \n\t")
		if breakAt <= 0 {
			breakAt = chunkSize
		}
		end = start + breakAt

		chunks = append(chunks, text[start:end])
		start = end
	}
	return chunks
}

// splitSentences splits text into sentences using simple rules.
func splitSentences(text string) []string {
	// Simple sentence splitter for parallel demo
	// A proper implementation would use the Punkt algorithm
	var sentences []string
	var buf strings.Builder

	for _, r := range text {
		buf.WriteRune(r)
		if r == '.' || r == '!' || r == '?' {
			sentences = append(sentences, strings.TrimSpace(buf.String()))
			buf.Reset()
		}
	}
	if buf.Len() > 0 {
		sentences = append(sentences, strings.TrimSpace(buf.String()))
	}
	return sentences
}

// defaultNumWorkers returns the default number of workers (GOMAXPROCS).
func defaultNumWorkers() int {
	// In a real implementation, this would use runtime.GOMAXPROCS(0)
	// For now, we'll just use a reasonable default
	return 4
}
