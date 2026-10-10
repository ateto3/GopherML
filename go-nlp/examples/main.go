// Example: Real-world API Gateway Pipeline using go-nlp
// This example demonstrates how to use go-nlp for preprocessing
// text before sending it to an LLM API, showing cost optimization
// through token reduction.

package main

import (
	"fmt"
	"strings"

	"github.com/ateto3/gopherml/go-nlp/cleaner"
	"github.com/ateto3/gopherml/go-nlp/tokenizer"
)

// AIPipeline represents a complete NLP pipeline for LLM preprocessing
type AIPipeline struct {
	cleaner   *cleaner.Cleaner
	tokenizer tokenizer.TokenizerI
}

// NewAIPipeline creates a new AI pipeline with optimized settings
func NewAIPipeline() *AIPipeline {
	return &AIPipeline{
		cleaner: cleaner.New(
			cleaner.WithPunctuation(true),
			cleaner.WithLowercase(true),
			cleaner.WithNormalizeSpace(true),
		),
		tokenizer: tokenizer.NewNLTKWordTokenizer(),
	}
}

// PreprocessPrompt cleans and optimizes a prompt before sending to LLM
func (p *AIPipeline) PreprocessPrompt(input string) (string, error) {
	// Step 1: Clean the input text
	cleaned := p.cleaner.Clean(input)

	// Step 2: Tokenize and check token count
	tokens := p.tokenizer.Tokenize(cleaned)
	tokenCount := len(tokens)

	if tokenCount > 4096 {
		return "", fmt.Errorf("input exceeds maximum token limit: %d > 4096", tokenCount)
	}

	fmt.Printf("Input tokens: %d\n", tokenCount)
	return cleaned, nil
}

// ProcessResponse processes the LLM response for the user
func (p *AIPipeline) ProcessResponse(output string) string {
	// Clean and format the output
	return p.cleaner.Clean(output)
}

func main() {
	fmt.Println("=== go-nlp AI Gateway Pipeline Example ===\n")

	// Create pipeline
	pipeline := NewAIPipeline()

	// Example 1: Preprocess a user prompt with fillers
	fmt.Println("Example 1: Prompt with fillers")
	rawPrompt := "Hello! Could you please help me with this? I have a question about the code. Thank you for your time."
	cleaned, _ := pipeline.PreprocessPrompt(rawPrompt)
	fmt.Printf("  Original: %s\n", rawPrompt)
	fmt.Printf("  Cleaned:  %s\n\n", cleaned)

	// Example 2: Process a Spanish prompt (pre-translation)
	fmt.Println("Example 2: Spanish prompt (pre-translation)")
	spanishPrompt := "Por favor, procesa este texto para mí. Es muy importante que me ayudes."
	cleaned, _ = pipeline.PreprocessPrompt(spanishPrompt)
	fmt.Printf("  Original: %s\n", spanishPrompt)
	fmt.Printf("  Cleaned:  %s\n\n", cleaned)

	// Example 3: Process a prompt with HTML and URLs
	fmt.Println("Example 3: Prompt with HTML and URLs")
	htmlPrompt := "<p>Check out this website: https://example.com and tell me what you think.</p>"
	cleaned, _ = pipeline.PreprocessPrompt(htmlPrompt)
	fmt.Printf("  Original: %s\n", htmlPrompt)
	fmt.Printf("  Cleaned:  %s\n\n", cleaned)

	// Example 4: Corpus statistics
	fmt.Println("\nExample 4: Corpus statistics")
	fmt.Println("  (requires corpus package initialization)")

	// Example 5: Tokenization comparison
	fmt.Println("\nExample 5: Tokenization methods")
	text := "Hello world! How are you doing today?"

	// Basic tokenizer
	tok := tokenizer.NewTreebankWordTokenizer()
	basicTokens := tok.Tokenize(text)
	fmt.Printf("  Basic tokens: %d\n", len(basicTokens))

	// Tokenize with span support
	spans := tok.SpanTokenize(text)
	fmt.Printf("  Token spans: %d\n", len(spans))
	for _, span := range spans {
		token := text[span.Start:span.End]
		fmt.Printf("    [%d:%d] = '%s'\n", span.Start, span.End, token)
	}

	// Example 6: Pipeline for multiple documents
	fmt.Println("\nExample 6: Batch processing")
	documents := []string{
		"The quick brown fox jumps over the lazy dog.",
		"Machine learning is transforming artificial intelligence.",
		"Natural language processing enables computers to understand text.",
	}

	for i, doc := range documents {
		cleaned, _ := pipeline.PreprocessPrompt(doc)
		tokens := strings.Fields(cleaned)
		fmt.Printf("  Doc %d: %d tokens - %s\n", i+1, len(tokens), strings.Join(tokens[:min(5, len(tokens))], " "))
	}

	fmt.Println("\n=== Example Complete ===")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
