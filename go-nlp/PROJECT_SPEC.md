# Project Specification: burst bubble AI API Gateway 🚀

## 📋 Overview

**burst bubble AI API Gateway** is an Open Core architecture that combines:

1. **`go-nlp` Library (Open Source)**: High-performance NLP tools in Go for local tokenization, text cleaning, stopword removal, and syntactic analysis with ultra-low latency.
2. **AI API Gateway (Private/Production Documentation)**: Go middleware that connects to marketplaces like **RapidAPI**, acting as an intelligent filter that drastically reduces token consumption (input/output) before interacting with LLM providers like OpenAI, Anthropic, or Google Gemini.

---

## 🎯 Core Objectives

### 1. Cost Optimization
- Reduce LLM API costs by 50-80% through intelligent preprocessing
- Extend token budgets by filtering before sending to LLMs
- Maximize throughput within rate limits

### 2. Performance
- Sub-millisecond local text preprocessing
- Zero-dependency Go implementation
- Concurrent processing for high-throughput scenarios

### 3. Flexibility
- Support multiple LLM providers (OpenAI, Anthropic, Google Gemini)
- Pluggable preprocessing pipeline
- Language-agnostic design

---

## 🛠️ Token Optimization Engineering

### Input Optimization (Pre-LLM)

#### 1. Elimination of Filler Words
**Rule**: Remove greetings, pleasantries, and unnecessary introductions from prompts.

**Implementation**:
```go
// Use cleaner package with custom stopword lists
cleaner := cleaner.New(
    cleaner.WithStopwords(true),
    cleaner.WithLanguage("en"),
)
cleaned := cleaner.Clean(prompt)
```

#### 2. Telegraphic Style Enforcement
**Rule**: Force models to use direct, compact language without rhetorical flourishes.

**Impact**: Up to 75% reduction in output tokens.

#### 3. English-First Translation Strategy
**Rule**: Translate to English before sending to LLM, translate results back.

**Why**: LLM tokenizers are optimized for English. Spanish characters, accents, and morphological structures fragment words into multiple tokens.

**Implementation**:
```go
// Use translate package for bilingual pipeline
translator := translate.New()
english := translator.TranslateToEnglish(spanishText)
result := sendToLLM(english)
spanish := translator.TranslateToSpanish(result)
```

**Tokenization Impact**:
- Spanish: "Café" → 2-3 tokens (depending on model)
- English: "Coffee" → 1 token

#### 4. Context Management
**Rule**: Use persistent memories and system variables in Go backend instead of repeatedly sending context.

**Implementation**:
```go
type Conversation struct {
    SystemPrompt   string
    ContextHistory []Message
    Memory         *context.Memory
}

func (c *Conversation) BuildPrompt(userInput string) string {
    // Build prompt with context from memory, not repeated in each request
    return fmt.Sprintf("%s\n\nUser: %s", c.SystemPrompt, userInput)
}
```

#### 5. Selective Editing
**Rule**: Request specific edits rather than full rewrites.

**Pattern**:
```
Original: "Edit only the third paragraph to fix grammar"
Instead of: "Rewrite the entire document with grammar fixes"
```

---

### Output Optimization (Post-LLM)

#### 1. Parameter Controls
**Rule**: Enforce hard limits on `max_tokens` in Go code.

**Implementation**:
```go
type Request struct {
    Model   string `json:"model"`
    MaxTokens int    `json:"max_tokens"` // Hard limit enforced by Go
    // ...
}

func (r *Request) EnforceLimits() {
    if r.MaxTokens > 2048 {
        r.MaxTokens = 2048 // Hard cap
    }
}
```

#### 2. Post-Processing Filters
**Rule**: Clean LLM output before returning to user.

**Implementation**:
```go
func PostProcessOutput(raw string) string {
    c := cleaner.New(
        cleaner.WithPunctuation(true),
        cleaner.WithNumbers(false),
        cleaner.WithLowercase(true),
    )
    return c.Clean(raw)
}
```

---

## 🏗️ Architecture

### Go-NLP Package Structure

```
go-nlp/
├── tokenizer/          # Tokenization algorithms (concurrent, multi-language)
├── cleaner/            # Text cleaning and sanitization
├── corpus/             # Corpus readers with caching and statistics
│   └── stopwords/      # Stopword dictionaries (optimized in-memory, O(1) lookup)
├── translate/          # Translation models (IBM Models 1-3)
├── chunk/              # Named entity chunking
├── classify/           # Classification models
├── cluster/            # Clustering algorithms
├── collocations/       # Bigram/trigram scoring
├── grammar/            # CFG, PCFG, feature grammars
├── inference/          # Logical inference
├── lm/                 # Language models
├── metrics/            # Evaluation metrics
├── parser/             # Syntactic parsers
├── probability/        # Probability distributions
├── sem/                # Semantic analysis
├── sentiment/          # Sentiment analysis
├── stemmer/            # Stemming and lemmatization
├── tagger/             # POS tagging
├── text/               # Text analysis utilities
├── tgrep/              # Tree search
├── tree/               # Syntax trees
├── featstruct/         # Feature structures
├── downloader/         # Model and corpus downloader
├── chat/               # Chatbot implementations
└── draw/               # Tree visualization
```

### AI API Gateway Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    User Request                              │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│              Go-NLP Preprocessing Layer                      │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ 1. Text Cleaning (HTML, URLs, Punctuation)          │   │
│  │ 2. Stopword Removal                                 │   │
│  │ 3. Length Filtering                                 │   │
│  │ 4. Language Detection & Translation (ES→EN)         │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│              Prompt Engineering Layer                       │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ 1. Context Injection (from persistent memory)        │   │
│  │ 2. System Prompt Assembly                            │   │
│  │ 3. Telegraphic Style Enforcement                     │   │
│  │ 4. Token Count Pre-check                             │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│              Model Routing Layer                            │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ Task Complexity Analysis                             │   │
│  │ ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │   │
│  │ │ Simple Task  │  │ Medium Task  │  │ Complex Task │ │   │
│  │ │ (Speed)      │  │ (Balance)    │  │ (Accuracy)   │ │   │
│  │ │ - Gemini     │  │ - Claude 3   │  │ - GPT-4      │ │   │
│  │ │ - Claude 3.5 │  │ - Llama 3    │  │ - Claude Opus│ │   │
│  │ │ - Flash      │  │ - Mistral    │  │ - GPT-4 Turbo│ │   │
│  │ └──────────────┘  └──────────────┘  └──────────────┘ │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│              API Gateway Layer                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ RapidAPI Integration                                 │   │
│  │ - OpenAI API                                         │   │
│  │ - Anthropic API                                      │   │
│  │ - Google Gemini API                                  │   │
│  │ - Custom Endpoints                                   │   │
│  │ - Prompt Caching Support                             │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│              Post-Processing Layer                          │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ 1. Output Cleaning                                   │   │
│  │ 2. Translation (EN→ES)                               │   │
│  │ 3. Token Count Verification                          │   │
│  │ 4. Rate Limit Checking                               │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

---

## 📦 Go-NLP Module Specifications

### 1. `tokenizer/` - Tokenization

#### Features
- **Concurrent Processing**: Worker pools for large documents
- **Multi-language Support**: Language-specific tokenizers
- **Span Tokenization**: Track original text positions
- **Custom Regex**: User-defined tokenization patterns

#### Tokenizers
| Tokenizer | Use Case | Language |
|-----------|----------|----------|
| `TreebankWordTokenizer` | Standard English | en |
| `NLTKWordTokenizer` | Improved Treebank | en |
| `TweetTokenizer` | Social media | multi |
| `RegexpTokenizer` | Custom patterns | any |
| `WhitespaceTokenizer` | Whitespace-based | any |
| `WordPunctTokenizer` | Words + punctuation | any |
| `CharTokenizer` | Character-level | any |
| `LineTokenizer` | Line-based | any |
| `MWETokenizer` | Multi-word expressions | any |

#### Usage
```go
import "github.com/ateto3/gopherml/go-nlp/tokenizer"

// Simple tokenization
tok := tokenizer.NewTreebankWordTokenizer()
tokens := tok.Tokenize("Hello world! How are you?")

// Concurrent tokenization for large texts
parallelTok := tokenizer.NewParallelTokenizer(64*1024, 0, tokenizer.NewNLTKWordTokenizer())
tokens := parallelTok.Tokenize(largeText)

// Span tokenization for text extraction
spans := tok.SpanTokenize("Hello world")
for _, span := range spans {
    token := text[span.Start:span.End]
}
```
Is it not pushed in github yet. So you can access to it by local using the path ~/Escritorio/GopherML/go-nlp

### 2. `cleaner/` - Text Cleaning

#### Features
- **HTML Removal**: Strip HTML tags and entities
- **URL/Email Removal**: Remove links and addresses
- **Punctuation Removal**: Clean punctuation marks
- **Stopword Removal**: Language-aware filtering
- **Case Normalization**: Lowercase/Uppercase conversion
- **Whitespace Normalization**: Normalize spacing
- **Emoji Removal**: Filter emoji characters
- **Accent Removal**: Remove diacritics
- **Length Filtering**: Min/max word length constraints

#### Usage
```go
import "github.com/ateto3/gopherml/go-nlp/cleaner"

// Basic cleaning
c := cleaner.New(
    cleaner.WithHTML(true),
    cleaner.WithPunctuation(true),
    cleaner.WithLowercase(true),
)
cleaned := c.Clean("<p>Hello, World!</p> https://example.com")

// Advanced cleaning with stopwords
c := cleaner.New(
    cleaner.WithStopwords(true),
    cleaner.WithLanguage("en"),
    cleaner.WithMinWordLength(2),
    cleaner.WithMaxWordLength(15),
)

// Custom stopword list
stopwords := map[string]bool{"quick": true, "brown": true}
c := cleaner.New(cleaner.WithStopwordList(stopwords))
c.AddStopwords("lazy", "dog")
```

### 3. `corpus/` - Corpus Readers

#### Features
- **Lazy Loading**: On-demand file loading
- **Caching**: Concurrent cache for raw/word data
- **Statistics**: Compute corpus statistics
- **Configurable**: Customizable splitting and filtering

#### Usage
```go
import "github.com/ateto3/gopherml/go-nlp/corpus"

// Basic corpus reader
r := corpus.NewFSReader(fsys, []string{"file1.txt", "file2.txt"})

// With configuration
r := corpus.NewFSReader(fsys, []string{"file1.txt"}, 
    corpus.NewLowercase(),
    corpus.NewRemovePunctuation(),
    corpus.NewMinWordLength(2),
)

// Statistics
stats, err := r.Stats()
fmt.Printf("Total words: %d\n", stats.TotalWords)
fmt.Printf("Vocabulary size: %d\n", stats.VocabSize)

// Lazy corpus
lazy := corpus.NewLazy(func() (corpus.Reader, error) {
    return corpus.NewFSReader(fsys, []string{"file.txt"}), nil
})
```

### 4. `translate/` - Translation

#### Features
- **IBM Model 1**: Statistical machine translation
- **Bilingual Corpora**: Support for aligned sentence pairs
- **Alignment**: Word-level alignment between languages

#### Usage
```go
import "github.com/ateto3/gopherml/go-nlp/translate"

// Train IBM Model 1
corpus := []translate.AlignedSent{
    {Words: []string{"the", "house"}, Mots: []string{"das", "haus"}},
    {Words: []string{"the", "book"}, Mots: []string{"das", "buch"}},
}
model := translate.NewIBMModel1(corpus, 10)

// Align new sentence
sent := &translate.AlignedSent{
    Words: []string{"the", "car"},
}
model.Align(sent)
// sent.Alignment contains word alignments
```

### 5. `corpus/stopwords/` - Stopword Dictionaries (NEW)

#### Features
- **Multi-language**: English, Spanish, Portuguese, French, German, Italian, Dutch, Russian, Chinese, Japanese
- **Optimized**: In-memory hash maps for O(1) lookup
- **Customizable**: Add/remove stopwords at runtime
- **Localized in corpus package**: Integrated with corpus readers for seamless preprocessing

#### Usage
```go
import "github.com/ateto3/gopherml/go-nlp/corpus/stopwords"

// Get stopwords for language
stopwordsEn := stopwords.Default("en")
stopwordsEs := stopwords.Default("es")

// Check if word is stopword
isStopword := stopwords.Contains("en", "the") // true

// Add/remove stopwords
stopwords.Add("en", "customword")
stopwords.Remove("en", "unwanted")

// List supported languages
langs := stopwords.Languages() // ["en", "es", "pt", ...]
```

#### Integration with Cleaner
```go
import (
    "github.com/ateto3/gopherml/go-nlp/cleaner"
    "github.com/ateto3/gopherml/go-nlp/corpus/stopwords"
)

// Use custom stopwords from corpus/stopwords
c := cleaner.New(
    cleaner.WithStopwords(true),
    cleaner.WithLanguage("es"),
)
// Automatically uses stopwords.Default("es")

// Or use custom stopword list
customStopwords := map[string]bool{
    "custom1": true,
    "custom2": true,
}
c = cleaner.New(cleaner.WithStopwordList(customStopwords))
```

---

## 🚀 Deployment Architecture (RapidAPI Monetization)

### RapidAPI Integration

```
┌─────────────────────────────────────────────────────────────┐
│                    RapidAPI Gateway                          │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ API Endpoint: /api/v1/llm                            │   │
│  │ Authentication: API Key                              │   │
│  │ Rate Limiting: Per-tenant quotas                     │   │
│  │ Usage Tracking: Token consumption analytics          │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│              burst bubble AI Server (Go)                    │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ HTTP Handler                                         │   │
│  │ - Parse request                                      │   │
│  │ - Validate API key                                   │   │
│  │ - Enforce rate limits                                │   │
│  │ - Process through go-nlp pipeline                    │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│              LLM Provider (OpenAI/Anthropic/Gemini)         │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ Cost Optimization: 50-80% reduction                  │   │
│  │ Token Budget: Extended usage                         │   │
│  │ Performance: Lower latency (English optimization)    │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### API Response Format

```json
{
  "success": true,
  "data": {
    "output": "Generated text",
    "input_tokens": 150,
    "output_tokens": 85,
    "total_tokens": 235,
    "model": "gpt-4-turbo",
    "language": "es"
  },
  "metadata": {
    "preprocessing_time_ms": 12,
    "llm_time_ms": 245,
    "cost_usd": 0.0035
  }
}
```

---

## 📊 Cost Savings Analysis

### Token Reduction Breakdown

| Optimization | Input Reduction | Output Reduction | Total Savings |
|--------------|-----------------|------------------|---------------|
| Stopword Removal | 15-20% | - | 15-20% |
| Filler Elimination | 10-15% | - | 10-15% |
| Telegraphic Style | - | 50-75% | 50-75% |
| English Translation | 20-30% | 20-30% | 40-60% |
| Context Management | 30-50% | - | 30-50% |
| Output Length Limits | - | 40-60% | 40-60% |

**Combined Potential**: 70-85% reduction in token consumption

### Example Cost Calculation

**Scenario**: 1M tokens/month baseline

| Scenario | Input Tokens | Output Tokens | Monthly Cost |
|----------|--------------|---------------|--------------|
| Without optimization | 500K | 500K | $15.00 |
| With go-nlp optimization | 150K | 125K | $3.75 |
| **Savings** | **350K tokens** | **375K tokens** | **$11.25/month** |

At scale (10M+ tokens/month), savings exceed $100/month per API.

---

## 🔐 Security & Privacy

### Data Handling
- All preprocessing happens server-side
- User data never leaves the customer's infrastructure
- No training on user prompts
- Compliance with GDPR, HIPAA, and SOC2

### Rate Limiting
- Per-API-key limits
- Burst protection
- Graceful degradation

---

## 📈 Monetization Strategy

### RapidAPI Publishing
1. **Public API**: Free tier with limited usage
2. **Paid Tiers**: Based on token consumption
3. **Enterprise**: Custom SLAs and dedicated infrastructure

### Pricing Model
| Tier | Input Tokens/Mo | Output Tokens/Mo | Price |
|------|-----------------|------------------|-------|
| Free | 10K | 10K | $0 |
| Basic | 100K | 100K | $5/mo |
| Pro | 1M | 1M | $25/mo |
| Enterprise | Custom | Custom | Contact |

---

## 🧪 Testing Strategy

### Unit Tests
- All go-nlp modules have comprehensive unit tests
- Coverage: >90% for core functionality

### Integration Tests
- End-to-end API gateway flow
- LLM provider integration
- Token counting accuracy

### Performance Tests
- Benchmark tokenization throughput
- Memory usage analysis
- Concurrent request handling

---

## 📚 Documentation

### Public (Open Source)
- `README.md`: Project overview
- `PROJECT_SPEC.md`: This document
- `examples/`: Usage examples
- `benchmarks/`: Performance benchmarks

### Private (Production)
- API documentation (Postman/ Swagger)
- Architecture diagrams
- Deployment guides
- Security policies

---

## 🔄 Future Enhancements

### Short-term
- [ ] Support for more LLM providers
- [ ] Dynamic prompt caching
- [ ] Real-time token counting

### Medium-term
- [ ] Custom model training interface
- [ ] Multi-tenant isolation
- [ ] Analytics dashboard

### Long-term
- [ ] Federated learning support
- [ ] Edge computing capabilities
- [ ] Custom tokenizer generation

---

## 🤝 Contribution Guidelines

### Open Source Contributions
1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

### Code Style
- Follow Go best practices
- 80-character line limit
- Comprehensive documentation
- Benchmark tests for performance-critical code

---

## 📄 License

This project follows the **Open Core** model:

- **`go-nlp` Library**: MIT License (open source)
- **AI API Gateway**: Proprietary (production use only)

See individual `LICENSE` files for details.

---

## 🙏 Acknowledgments

- NLTK team for the original Python implementation
- Go community for the excellent tooling
- Open source contributors and users

---

**Built with ❤️ using Go for high-performance NLP**
