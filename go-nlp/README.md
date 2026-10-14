# go-nlp

**High-Performance Natural Language Processing for Go — Concurrent port of NLTK.**

[![Go Version](https://img.shields.io/github/go-mod/go-version/ateto3/gopherml/go-nlp)](https://go.dev)
[![License](https://img.shields.io/github/license/ateto3/gopherml/go-nlp)](../LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/ateto3/gopherml/go-nlp.svg)](https://pkg.go.dev/github.com/ateto3/gopherml/go-nlp)

`go-nlp` is a comprehensive, high-performance solution developed in **Go** for natural language processing (NLP) and cost optimization in the consumption of Artificial Intelligence APIs (LLMs). This repository follows the * * Open Core * * architectural model:

1. **`go-nlp` Library (Open Source)**: Native Go tools for tokenization, text cleaning, stopword removal, and ultra-low latency local syntactic analysis.
2. **AI API Gateway (Production/Private Documentation)**: A Go middleware architecture designed to connect to marketplaces like **RapidAPI**, acting as an intelligent filter that drastically reduces token consumption (input/output) before interacting with providers like OpenAI, Anthropic, or Google Gemini.

---

## 🚀 Why `go-nlp`?

NLTK is the most complete NLP toolkit in Python, but its execution is sequential and suffers from GIL (Global Interpreter Lock). `go-nlp` rescues its modular philosophy and algorithmic richness, but implements it from scratch in Go with:

- **Concurrent tokenization:** Processing complete documents by dividing text into chunks that are tokenized in parallel via worker pools.
- **Parallel POS tagging:** Tagging models (HMM, Perceptron, TnT, Brill) run on independent goroutines per sentence.
- **Distributed parsing:** Context-free grammars, PCFG, dependency parsing, and CCG executed with parallelizable chart parsing strategies (CYK, Earley).
- **Scalable classification:** Naive Bayes, MaxEnt, Decision Trees, SVM, and clustering (K-Means, EM, GAAC) with native support for distributed training.
- **Sentiment analysis at scale:** Streaming processing of large text volumes with concurrent pipelines.
- **Zero-copy string processing:** Intensive use of `[]byte` and `strings.Builder` to minimize GC pressure.
- **Zero-Dependency:** Minimal external dependencies. Own algorithm implementations without resorting to C or Python bindings.

---

## 📊 Cost Optimization Results

| Optimization | Input Reduction | Output Reduction | Total Savings |
|--------------|-----------------|------------------|--------------|
| Stopword Removal | 15-20% | - | 15-20% |
| Filler Elimination | 10-15% | - | 10-15% |
| Telegraphic Style | - | 50-75% | 50-75% |
| English Translation | 20-30% | 20-30% | 40-60% |
| Context Management | 30-50% | - | 30-50% |
| Output Length Limits | - | 40-60% | 40-60% |

**Combined Potential**: 70-85% reduction in token consumption

Savings Example

**Scenario**: 1M tokens/mes baseline

| Scenario | Input Tokens | Output Tokens | Monthly Cost |
|----------|--------------|---------------|--------------|
| Without optimization | 500K | 500K | $15.00 |
| With go-nlp optimization | 150K | 125K | $3.75 |
| **Savings** | **350K tokens** | **375K tokens** | **$11.25/month** |

---

## NLTK Legacy 🧠 Architecture — Powered by Go

The following table maps the **NLTK (Python)** packages with the ** go-nlp (Go)** modules:

| NLTK Module (Python)       | go-nlp Module (Go)               | Concurrent Improvement                          |
|----------------------------|----------------------------------|---------------------------------------------|
| `nltk.tokenize`            | `tokenizer/`                     | Parallel chunk tokenization            |
| `nltk.tag`                 | `tagger/`                        | Multi-goroutine POS tagging              |
| `nltk.parse`               | `parser/`                        | Concurrent chart parsing                   |
| `nltk.chunk`               | `chunker/`                       | Parallel chunking over tags                |
| `nltk.classify`            | `classify/`                      | Distributed training                   |
| `nltk.sentiment`           | `sentiment/`                     | Streaming analysis pipeline           |
| `nltk.stem`                | `stemmer/`                       | Batch stemming/lemmatization                 |
| `nltk.corpus`              | `corpus/`                        | Lazy loading with concurrent cache            |
| `nltk.probability`         | `probability/`                   | FreqDist, ProbDist, Witten-Bell, Good-Turing|
| `nltk.collocations`        | `collocations/`                  | Parallel Bigram/Trigram scoring             |
| `nltk.text`                | `text/`                          | Concordance, contextual similarity          |
| `nltk.metrics`             | `metrics/`                       | Agreement, scores, confusion matrix         |
| `nltk.cluster`             | `cluster/`                       | Parallel K-Means, EM, GAAC                 |
| `nltk.sem`                 | `sem/`                           | First-order logic, lambda calculus      |
| `nltk.inference`           | `inference/`                     | Async Prover9/Mace4 wrapper             |
| `nltk.lm`                  | `lm/`                            | Language models (MLE, Kneser-Ney)       |
| `nltk.tree`                | `tree/`                          | Immutable syntactic trees              |
| `nltk.featstruct`          | `featstruct/`                    | Feature structures with unification          |
| `nltk.grammar`             | `grammar/`                       | CFG, PCFG, Feature Grammar                  |
| `nltk.translate`           | `translate/`                     | IBM Models 1-3, statistical models                     |
| `nltk.wsd`                 | `wsd/`                           | Lesk algorithm, sense disambiguation        |
| `nltk.downloader`          | `downloader/`                    | Model and corpus download                |
| `nltk.tgrep`               | `tgrep/`                         | TGrep2 search on trees               |
| `nltk.chat`                | `chat/`                          | Chatbots (Eliza, Iesha, Sun Tzu, Zen)       |
| `nltk.draw`                | `draw/`                          | Tree and dispersion visualization       |

---

## 'go-nlp` Library 📦 Modules

### Main Packages

| Package | Description | Status |
|---------|-------------|--------|
| `tokenizer/` | String segmentation into tokens with multi-language support and concurrent processing | ✅ Stable |
| `cleaner/` | Noise removal, punctuation, special characters, and normalization | ✅ Stable |
| `corpus/` | Corpus readers with cache and statistics | ✅ Stable |
| `translate/` | Statistical translation models (IBM Models 1-3) | ✅ Stable |
| `stopwords/` | Optimized dictionaries for instant filler word removal | ✅ Stable |
| `chunk/` | Named entity chunking | ✅ Stable |
| `classify/` | Classification (Naive Bayes, MaxEnt, Decision Trees) | ✅ Stable |
| `cluster/` | Clustering (K-Means, EM, GAAC) | ✅ Stable |
| `collocations/` | Bigram/Trigram scoring | ✅ Stable |
| `grammar/` | Context-free grammars (CFG, PCFG) | ✅ Stable |
| `parser/` | Syntactic parsers (CYK, Earley, Chart) | ✅ Stable |
| `lm/` | Language models (MLE, Kneser-Ney) | ✅ Stable |
| `metrics/` | Evaluation metrics and confusion matrix | ✅ Stable |
| `probability/` | Probability distributions | ✅ Stable |

### Advanced Packages

| Package | Description | Status |
|---------|-------------|--------|
| `tagger/` | POS tagging (HMM, Perceptron, TnT) | 🚧 In Progress |
| `stemmer/` | Stemming and lemmatization | 🚧 In Progress |
| `sentiment/` | Sentiment Analysis | 🚧 In Progress |
| `text/` | Contextual concordance and similarity | 🚧 In Progress |
| `wsd/` | Word without disambiguation | 🚧 In Progress |
| `chat/` | Chatbots (Elizo, Iesha, Sun Tzu) | 🚧 In Progress |
| `draw/` | Tree display | 🚧 In Progress |

---

## Quick 🚀 Use

Tokenization

```go
package main

import (
    "fmt"
    "github.com/ateto3/gopherml/go-nlp/tokenizer"
)

func main() {
    // Basic tokenizer
    tok := tokenizer.NewTreebankWordTokenizer()
    tokens := tok.Tokenize("Hello world! How are you doing today?")
    fmt.Printf("Tokens: %v\n", tokens)
    
    // Parallel tokenization for large text
    parallelTok := tokenizer.NewParallelTokenizer(64*1024, 0, tokenizer.NewNLTKWordTokenizer())
    largeTokens := parallelTok.Tokenize(largeText)
    
    // Tweet tokenization
    tweetTok := tokenizer.NewDefaultTweetTokenizer()
    tweetTokens := tweetTok.Tokenize("RT @user: Hello World! Check out https://example.com #golang")
    fmt.Printf("Tweet tokens: %v\n", tweetTokens)
}
```

### Text Cleaning

```go
package main

import (
    "fmt"
    "github.com/ateto3/gopherml/go-nlp/cleaner"
)

func main() {
    // Basic cleaning
    c := cleaner.New(
        cleaner.WithHTML(true),
        cleaner.WithPunctuation(true),
        cleaner.WithLowercase(true),
    )
    cleaned := c.Clean("<p>Hello, World!</p> https://example.com @user #test")
    fmt.Printf("Cleaned: %s\n", cleaned)
    
    // Cleaning with stopwords
    c := cleaner.New(
        cleaner.WithStopwords(true),
        cleaner.WithLanguage("en"),
    )
    cleaned = c.Clean("The quick brown fox jumps over the lazy dog")
    fmt.Printf("Without stopwords: %s\n", cleaned)
    
    // Token cleaning
    tokens := []string{"the", "quick", "brown", "fox", "jumps"}
    cleanedTokens := c.CleanTokens(tokens)
    fmt.Printf("Cleaned tokens: %v\n", cleanedTokens)
}
```

### Corpus Reader

```go
package main

import (
    "fmt"
    "os"
    "github.com/ateto3/gopherml/go-nlp/corpus"
)

func main() {
    // Leer corpus de archivos
    files := []string{"file1.txt", "file2.txt"}
    r := corpus.NewFSReader(os.DirFS("data"), files)
    
    words, _ := r.Words()
    fmt.Printf("Total words: %d\n", len(words))
    
    // With statistics
    stats, _ := r.Stats()
    fmt.Printf("Vocabulary size: %d\n", stats.VocabSize)
    
    // Lazy loading
    lazy := corpus.NewLazy(func() (corpus.Reader, error) {
        return corpus.NewFSReader(os.DirFS("data"), []string{"file.txt"}), nil
    })
}
```

Translation

```go
package main

import (
    "fmt"
    "github.com/ateto3/gopherml/go-nlp/translate"
)

func main() {
    // Entrenar modelo IBM 1
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
    fmt.Printf("Alignment: %s\n", sent.Alignment.String())
}
```

### Pipeline Completo (API Gateway Pattern)

```go
package main

import (
    "fmt"
    "github.com/ateto3/gopherml/go-nlp/cleaner"
    "github.com/ateto3/gopherml/go-nlp/tokenizer"
    "github.com/ateto3/gopherml/go-nlp/translate"
)

type AIPipeline struct {
    cleaner   *cleaner.Cleaner
    translator *translate.Translator
    tokenizer  *tokenizer.Tokenizer
}

func NewAIPipeline() *AIPipeline {
    return &AIPipeline{
        cleaner: cleaner.New(
            cleaner.WithStopwords(true),
            cleaner.WithLanguage("en"),
        ),
        translator: translate.New(),
        tokenizer:  tokenizer.NewNLTKWordTokenizer(),
    }
}

func (p *AIPipeline) ProcessPrompt(input string) (string, error) {
    // 1. Clean input
    cleaned := p.cleaner.Clean(input)
    
    // 2. Translate to English if needed
    english := p.translator.TranslateToEnglish(cleaned)
    
    // 3. Tokenize and check length
    tokens := p.tokenizer.Tokenize(english)
    if len(tokens) > 4096 {
        return "", fmt.Errorf("input too long")
    }
    
    return english, nil
}

func (p *AIPipeline) ProcessResponse(output string) (string, error) {
    // 4. Translate back to Spanish
    spanish := p.translator.TranslateToSpanish(output)
    
    // 5. Clean output
    cleaned := p.cleaner.Clean(spanish)
    
    return cleaned, nil
}
```

---

## 🏗️ Arquitectura de Despliegue

### Monetization in RapidAPI

```
User Request → Go-NLP Preprocessing → AI API Gateway → LLM Provider
                                      ↓
                              Token Caching
                                      ↓
                              Cost Optimization (50-80%)
```

### Gateway Components

1. **Preprocessing Layer**: go-nlp libraries for cleaning and tokenization
2. **Prompt Engineering**: Context injection and system prompt assembly
3. **Model Routing**: Model selection based on task complexity
4. **API Gateway**: Integration with RapidAPI and LLM providers
5. **Post-processing**: Translation and output cleaning

---

Documentation

- `PROJECT_SPEC.md`: Full project specification and architecture
- `examples/`: Examples of use by module
- `benchmarks/`: Performance benchmarks

---

Contribution

### Contribution Guides

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Code of conduct

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

## 📞 Contact

For open source questions: [GitHub Issues](https://github.com/ateto3/gopherml/issues)  
For commercial inquiries: [burst bubble AI](https://burstable.ai)

---

**Built with ❤️ using Go for high-performance NLP**
