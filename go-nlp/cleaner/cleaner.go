// Package cleaner provides text cleaning utilities for NLP preprocessing.
//
// This package includes utilities for:
//   - HTML tag removal
//   - URL and email removal
//   - Special character removal
//   - Whitespace normalization
//   - Case normalization
//   - Stop word removal (from corpus/stopwords)
//   - Emoji removal
//   - Punctuation removal
//   - Number normalization
//   - Unicode normalization
package cleaner

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/ateto3/gopherml/go-nlp/corpus/stopwords"
)

// Options holds configuration for text cleaning.
type Options struct {
	RemoveHTML             bool
	RemoveURLs             bool
	RemoveEmails           bool
	RemoveHashtags         bool
	RemoveMentions         bool
	RemovePunctuation      bool
	Lowercase              bool
	Uppercase              bool
	NormalizeSpace         bool
	RemoveNumbers          bool
	RemoveStopwords        bool
	RemoveEmoji            bool
	RemoveAccents          bool
	RemoveExtraWhitespaces bool
	MinWordLength          int
	MaxWordLength          int
	Language               string
	Stopwords              map[string]bool
}

// DefaultOptions returns default cleaning options.
func DefaultOptions() *Options {
	return &Options{
		RemoveHTML:             false,
		RemoveURLs:             false,
		RemoveEmails:           false,
		RemoveHashtags:         false,
		RemoveMentions:         false,
		RemovePunctuation:      false,
		Lowercase:              false,
		Uppercase:              false,
		NormalizeSpace:         false,
		RemoveNumbers:          false,
		RemoveStopwords:        false,
		RemoveEmoji:            false,
		RemoveAccents:          false,
		RemoveExtraWhitespaces: false,
		MinWordLength:          1,
		MaxWordLength:          100,
		Language:               "en",
		Stopwords:              nil,
	}
}

// OptionsFunc is a function that modifies Options.
type OptionsFunc func(*Options)

// WithHTML sets HTML removal.
func WithHTML(remove bool) OptionsFunc {
	return func(o *Options) { o.RemoveHTML = remove }
}

// WithURLs sets URL removal.
func WithURLs(remove bool) OptionsFunc {
	return func(o *Options) { o.RemoveURLs = remove }
}

// WithEmails sets email removal.
func WithEmails(remove bool) OptionsFunc {
	return func(o *Options) { o.RemoveEmails = remove }
}

// WithPunctuation sets punctuation removal.
func WithPunctuation(remove bool) OptionsFunc {
	return func(o *Options) { o.RemovePunctuation = remove }
}

// WithLowercase sets lowercase conversion.
func WithLowercase(lower bool) OptionsFunc {
	return func(o *Options) { o.Lowercase = lower }
}

// WithNormalizeSpace sets whitespace normalization.
func WithNormalizeSpace(normalize bool) OptionsFunc {
	return func(o *Options) { o.NormalizeSpace = normalize }
}

// WithStopwords sets stopword removal.
func WithStopwords(remove bool) OptionsFunc {
	return func(o *Options) { o.RemoveStopwords = remove }
}

// WithLanguage sets the language for stopwords.
func WithLanguage(lang string) OptionsFunc {
	return func(o *Options) { o.Language = lang }
}

// WithStopwordList sets a custom stopword list.
func WithStopwordList(words map[string]bool) OptionsFunc {
	return func(o *Options) { o.Stopwords = words }
}

// WithEmoji sets emoji removal.
func WithEmoji(remove bool) OptionsFunc {
	return func(o *Options) { o.RemoveEmoji = remove }
}

// WithAccents sets accent removal.
func WithAccents(remove bool) OptionsFunc {
	return func(o *Options) { o.RemoveAccents = remove }
}

// Cleaner provides text cleaning functionality.
type Cleaner struct {
	Options *Options
}

// New creates a new Cleaner with given options.
func New(opts ...OptionsFunc) *Cleaner {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &Cleaner{Options: options}
}

// Clean cleans the input text.
func (c *Cleaner) Clean(text string) string {
	if text == "" {
		return ""
	}

	result := text

	// Remove HTML tags
	if c.Options.RemoveHTML {
		result = removeHTML(result)
	}

	// Remove URLs
	if c.Options.RemoveURLs {
		result = removeURLs(result)
	}

	// Remove emails
	if c.Options.RemoveEmails {
		result = removeEmails(result)
	}

	// Remove mentions
	if c.Options.RemoveMentions {
		result = removeMentions(result)
	}

	// Remove hashtags
	if c.Options.RemoveHashtags {
		result = removeHashtags(result)
	}

	// Remove emoji
	if c.Options.RemoveEmoji {
		result = removeEmoji(result)
	}

	// Remove accents
	if c.Options.RemoveAccents {
		result = removeAccents(result)
	}

	// Remove punctuation
	if c.Options.RemovePunctuation {
		result = removePunctuation(result)
	}

	// Remove numbers
	if c.Options.RemoveNumbers {
		result = removeNumbers(result)
	}

	// Normalize case
	if c.Options.Lowercase {
		result = strings.ToLower(result)
	} else if c.Options.Uppercase {
		result = strings.ToUpper(result)
	}

	// Normalize whitespace
	if c.Options.NormalizeSpace || c.Options.RemoveExtraWhitespaces {
		result = normalizeWhitespace(result)
	}

	// Remove stopwords
	if c.Options.RemoveStopwords {
		if c.Options.Stopwords == nil {
			c.Options.Stopwords = stopwords.Default(c.Options.Language).Words()
		}
		result = removeStopwords(result, c.Options.Stopwords, "")
	}

	// Filter by word length
	result = filterByWordLength(result, c.Options.MinWordLength, c.Options.MaxWordLength)

	return strings.TrimSpace(result)
}

// CleanTokens cleans a slice of tokens.
func (c *Cleaner) CleanTokens(tokens []string) []string {
	if c.Options.Stopwords == nil {
		c.Options.Stopwords = stopwords.Default(c.Options.Language).Words()
	}

	var result []string
	for _, tok := range tokens {
		cleaned := strings.TrimSpace(tok)

		// Skip empty tokens
		if cleaned == "" {
			continue
		}

		// Remove punctuation
		if c.Options.RemovePunctuation {
			cleaned = removePunctuation(cleaned)
		}

		// Remove numbers
		if c.Options.RemoveNumbers {
			cleaned = removeNumbers(cleaned)
		}

		// Skip if empty after cleaning
		if cleaned == "" {
			continue
		}

		// Normalize case
		if c.Options.Lowercase {
			cleaned = strings.ToLower(cleaned)
		} else if c.Options.Uppercase {
			cleaned = strings.ToUpper(cleaned)
		}

		// Filter by length
		if len(cleaned) < c.Options.MinWordLength || len(cleaned) > c.Options.MaxWordLength {
			continue
		}

		// Remove stopwords
		if c.Options.RemoveStopwords {
			if c.Options.Stopwords[cleaned] {
				continue
			}
		}

		result = append(result, cleaned)
	}

	return result
}

// removeHTML removes HTML tags from text.
func removeHTML(text string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(text, "")
}

// removeURLs removes URLs from text.
func removeURLs(text string) string {
	re := regexp.MustCompile(`https?://[^\s]+|www\.[^\s]+|[^\s]+\.(com|org|net|edu|gov|io|co)[^\s]*`)
	return re.ReplaceAllString(text, "")
}

// removeEmails removes email addresses from text.
func removeEmails(text string) string {
	re := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	return re.ReplaceAllString(text, "")
}

// removeMentions removes @mentions from text.
func removeMentions(text string) string {
	re := regexp.MustCompile(`@\w+`)
	return re.ReplaceAllString(text, "")
}

// removeHashtags removes #hashtags from text.
func removeHashtags(text string) string {
	re := regexp.MustCompile(`#\w+`)
	return re.ReplaceAllString(text, "")
}

// removeEmoji removes emoji from text.
func removeEmoji(text string) string {
	re := regexp.MustCompile(`[\U0001F600-\U0001F64F]|` +
		`[\U0001F300-\U0001F5FF]|` +
		`[\U0001F680-\U0001F6FF]|` +
		`[\U0001F1E0-\U0001F1FF]|` +
		`[\U00002702-\U000027B0]|` +
		`[\U0001F900-\U0001F9FF]|` +
		`[\U0001FA70-\U0001FAFF]|` +
		`[\U00002600-\U000026FF]`)
	return re.ReplaceAllString(text, "")
}

// removePunctuation removes punctuation from text.
func removePunctuation(text string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) {
			return -1
		}
		return r
	}, text)
}

// removeNumbers removes digits from text.
func removeNumbers(text string) string {
	re := regexp.MustCompile(`\d+`)
	return re.ReplaceAllString(text, "")
}

// normalizeWhitespace normalizes whitespace in text.
func normalizeWhitespace(text string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(text, " ")
}

// filterByWordLength filters words by length.
func filterByWordLength(text string, min, max int) string {
	words := strings.Fields(text)
	var result []string
	for _, w := range words {
		if len(w) >= min && len(w) <= max {
			result = append(result, w)
		}
	}
	return strings.Join(result, " ")
}

// removeAccents removes accents from text.
func removeAccents(text string) string {
	// This is a simplified version; a full implementation would use unicode/utf8
	// and mapping tables for proper accent removal
	replacements := map[rune]rune{
		'á': 'a', 'é': 'e', 'í': 'i', 'ó': 'o', 'ú': 'u',
		'Á': 'A', 'É': 'E', 'Í': 'I', 'Ó': 'O', 'Ú': 'U',
		'ñ': 'n', 'Ñ': 'N',
	}
	return strings.Map(func(r rune) rune {
		if rep, ok := replacements[r]; ok {
			return rep
		}
		return r
	}, text)
}

// removeStopwords removes stopwords from text.
func removeStopwords(text string, stopwords map[string]bool, lang string) string {
	words := strings.Fields(text)
	var result []string
	for _, w := range words {
		if !stopwords[strings.ToLower(w)] {
			result = append(result, w)
		}
	}
	return strings.Join(result, " ")
}

// Stopwords returns a copy of stopwords for the configured language.
func (c *Cleaner) Stopwords() map[string]bool {
	return stopwords.Default(c.Options.Language).Words()
}

// AddStopwords adds words to the stopword list.
func (c *Cleaner) AddStopwords(words ...string) {
	if c.Options.Stopwords == nil {
		c.Options.Stopwords = stopwords.Default(c.Options.Language).Words()
	}
	for _, w := range words {
		c.Options.Stopwords[strings.ToLower(w)] = true
	}
}

// RemoveStopwords removes words from the stopword list.
func (c *Cleaner) RemoveStopwords(words ...string) {
	if c.Options.Stopwords == nil {
		return
	}
	for _, w := range words {
		delete(c.Options.Stopwords, strings.ToLower(w))
	}
}

// CleanFunc is a function that cleans text.
type CleanFunc func(string) string

// Pipe creates a chain of cleaning functions.
func Pipe(funcs ...CleanFunc) CleanFunc {
	return func(text string) string {
		for _, f := range funcs {
			text = f(text)
		}
		return text
	}
}

// Various cleaning functions as standalone functions.
func RemoveHTML(text string) string          { return removeHTML(text) }
func RemoveURLs(text string) string          { return removeURLs(text) }
func RemoveEmails(text string) string        { return removeEmails(text) }
func RemoveMentions(text string) string      { return removeMentions(text) }
func RemoveHashtags(text string) string      { return removeHashtags(text) }
func RemoveEmoji(text string) string         { return removeEmoji(text) }
func RemovePunctuation(text string) string   { return removePunctuation(text) }
func RemoveNumbers(text string) string       { return removeNumbers(text) }
func NormalizeWhitespace(text string) string { return normalizeWhitespace(text) }
func Lowercase(text string) string           { return strings.ToLower(text) }
func Uppercase(text string) string           { return strings.ToUpper(text) }

// WithNumbers sets number removal.
func WithNumbers(remove bool) OptionsFunc {
	return func(o *Options) { o.RemoveNumbers = remove }
}

// WithMinWordLength sets minimum word length.
func WithMinWordLength(min int) OptionsFunc {
	return func(o *Options) { o.MinWordLength = min }
}

// WithMaxWordLength sets maximum word length.
func WithMaxWordLength(max int) OptionsFunc {
	return func(o *Options) { o.MaxWordLength = max }
}

// WithUppercase sets uppercase conversion.
func WithUppercase(upper bool) OptionsFunc {
	return func(o *Options) { o.Uppercase = upper }
}
