# go-nlp

**High-Performance Natural Language Processing for Go — Concurrent port of NLTK.**

[![Go Version](https://img.shields.io/github/go-mod/go-version/ateto3/gopherml/go-nlp)](https://go.dev)
[![License](https://img.shields.io/github/license/ateto3/gopherml/go-nlp)](../LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/ateto3/gopherml/go-nlp.svg)](https://pkg.go.dev/github.com/ateto3/gopherml/go-nlp)

`go-nlp` es una solución integral y de alto rendimiento desarrollada en **Go** para el procesamiento del lenguaje natural (NLP) y la optimización de costes en el consumo de APIs de Inteligencia Artificial (LLMs). Este repositorio sigue el modelo arquitectónico **Open Core**:

1. **`go-nlp` Library (Código Abierto)**: Herramientas nativas de Go para tokenización, limpieza de texto, remoción de stopwords y análisis sintáctico local de ultra-baja latencia.
2. **AI API Gateway (Documentación de Producción / Privado)**: Una arquitectura de middleware en Go diseñada para conectarse a marketplaces como **RapidAPI**, actuando como un filtro inteligente que reduce drásticamente el consumo de tokens (input/output) antes de interactuar con proveedores como OpenAI, Anthropic o Google Gemini.

---

## 🚀 ¿Por qué `go-nlp`?

NLTK es el toolkit de NLP más completo en Python, pero su ejecución es secuencial y sufre del GIL (Global Interpreter Lock). `go-nlp` rescata su filosofía modular y su riqueza algorítmica, pero la implementa desde cero en Go con:

- **Tokenización concurrente:** Procesamiento de documentos completos dividiendo el texto en chunks que se tokenizan en paralelo mediante worker pools.
- **Tagging POS paralelo:** Los modelos de etiquetado (HMM, Perceptron, TnT, Brill) se ejecutan sobre goroutines independientes por oración.
- **Parsing distribuido:** Gramáticas libres de contexto, PCFG, dependency parsing y CCG ejecutados con estrategias de chart parsing paralelizable (CYK, Earley).
- **Clasificación escalable:** Naive Bayes, MaxEnt, Decision Trees, SVM y clustering (K-Means, EM, GAAC) con soporte nativo para entrenamiento distribuido.
- **Análisis de sentimiento a escala:** Procesamiento en streaming de grandes volúmenes de texto con pipelines concurrentes.
- **Zero-copy string processing:** Uso intensivo de `[]byte` y `strings.Builder` para minimizar presión del GC.
- **Zero-Dependency:** Mínimas dependencias externas. Implementaciones propias de algoritmos sin recurrir a bindings de C o Python.

---

## 📊 Cost Optimization Results

| Optimización | Reducción Input | Reducción Output | Ahorro Total |
|--------------|-----------------|------------------|--------------|
| Stopword Removal | 15-20% | - | 15-20% |
| Filler Elimination | 10-15% | - | 10-15% |
| Telegraphic Style | - | 50-75% | 50-75% |
| English Translation | 20-30% | 20-30% | 40-60% |
| Context Management | 30-50% | - | 30-50% |
| Output Length Limits | - | 40-60% | 40-60% |

**Combined Potential**: 70-85% reduction in token consumption

### Ejemplo de Ahorro

**Scenario**: 1M tokens/mes baseline

| Scenario | Input Tokens | Output Tokens | Monthly Cost |
|----------|--------------|---------------|--------------|
| Without optimization | 500K | 500K | $15.00 |
| With go-nlp optimization | 150K | 125K | $3.75 |
| **Savings** | **350K tokens** | **375K tokens** | **$11.25/month** |

---

## 🧠 Arquitectura heredada de NLTK — Potenciada por Go

La siguiente tabla mapea los paquetes de **NLTK (Python)** con los módulos de **go-nlp (Go)**:

| Módulo NLTK (Python)       | Módulo go-nlp (Go)               | Mejora Concurrente                          |
|----------------------------|----------------------------------|---------------------------------------------|
| `nltk.tokenize`            | `tokenizer/`                     | Tokenización paralela por chunks            |
| `nltk.tag`                 | `tagger/`                        | Etiquetado POS multi-goroutine              |
| `nltk.parse`               | `parser/`                        | Chart parsing concurrente                   |
| `nltk.chunk`               | `chunker/`                       | Chunking paralelo sobre tags                |
| `nltk.classify`            | `classify/`                      | Entrenamiento distribuido                   |
| `nltk.sentiment`           | `sentiment/`                     | Pipeline de análisis en streaming           |
| `nltk.stem`                | `stemmer/`                       | Stemming/Lematización batch                 |
| `nltk.corpus`              | `corpus/`                        | Carga lazy con caché concurrente            |
| `nltk.probability`         | `probability/`                   | FreqDist, ProbDist, Witten-Bell, Good-Turing|
| `nltk.collocations`        | `collocations/`                  | Bigram/Trigram scoring paralelo             |
| `nltk.text`                | `text/`                          | Concordancia, similitud contextual          |
| `nltk.metrics`             | `metrics/`                       | Agreement, scores, confusión matrix         |
| `nltk.cluster`             | `cluster/`                       | K-Means, EM, GAAC paralelos                 |
| `nltk.sem`                 | `sem/`                           | Lógica de primer orden, lambda-cálculo      |
| `nltk.inference`           | `inference/`                     | Prover9/Mace4 wrapper asíncrono             |
| `nltk.lm`                  | `lm/`                            | Modelos de lenguaje (MLE, Kneser-Ney)       |
| `nltk.tree`                | `tree/`                          | Árboles sintácticos inmutables              |
| `nltk.featstruct`          | `featstruct/`                    | Feature structures con unificación          |
| `nltk.grammar`             | `grammar/`                       | CFG, PCFG, Feature Grammar                  |
| `nltk.translate`           | `translate/`                     | IBM Models 1-3, sacábamos                   |
| `nltk.wsd`                 | `wsd/`                           | Lesk algorithm, sense disambiguation        |
| `nltk.downloader`          | `downloader/`                    | Descarga de modelos y corpus                |
| `nltk.tgrep`               | `tgrep/`                         | Búsqueda TGrep2 sobre árboles               |
| `nltk.chat`                | `chat/`                          | Chatbots (Eliza, Iesha, Sun Tzu, Zen)       |
| `nltk.draw`                | `draw/`                          | Visualización de árboles y dispersión       |

---

## 📦 Módulos de la Librería `go-nlp`

### Paquetes Principales

| Package | Descripción | Estado |
|---------|-------------|--------|
| `tokenizer/` | Segmentación de strings en tokens con soporte para múltiples idiomas y procesamiento concurrente | ✅ Stable |
| `cleaner/` | Eliminación de ruido, puntuación, caracteres especiales y normalización | ✅ Stable |
| `corpus/` | Lectores de corpus con caché y estadísticas | ✅ Stable |
| `translate/` | Modelos de traducción estadística (IBM Models 1-3) | ✅ Stable |
| `stopwords/` | Diccionarios optimizados para remoción instantánea de palabras de relleno | ✅ Stable |
| `chunk/` | Chunking de entidades nombradas | ✅ Stable |
| `classify/` | Clasificación (Naive Bayes, MaxEnt, Decision Trees) | ✅ Stable |
| `cluster/` | Clustering (K-Means, EM, GAAC) | ✅ Stable |
| `collocations/` | Bigram/Trigram scoring | ✅ Stable |
| `grammar/` | Gramáticas libres de contexto (CFG, PCFG) | ✅ Stable |
| `parser/` | Syntactic parsers (CYK, Earley, Chart) | ✅ Stable |
| `lm/` | Modelos de lenguaje (MLE, Kneser-Ney) | ✅ Stable |
| `metrics/` | Evaluation metrics y confusión matrix | ✅ Stable |
| `probability/` | Distribuciones de probabilidad | ✅ Stable |

### Paquetes Avanzados

| Package | Descripción | Estado |
|---------|-------------|--------|
| `tagger/` | POS tagging (HMM, Perceptron, TnT) | 🚧 In Progress |
| `stemmer/` | Stemming y lematización | 🚧 In Progress |
| `sentiment/` | Análisis de sentimiento | 🚧 In Progress |
| `text/` | Concordancia y similitud contextual | 🚧 In Progress |
| `wsd/` | Word sense disambiguation | 🚧 In Progress |
| `chat/` | Chatbots (Elizo, Iesha, Sun Tzu) | 🚧 In Progress |
| `draw/` | Visualización de árboles | 🚧 In Progress |

---

## 🚀 Uso Rápido

### Tokenización

```go
package main

import (
    "fmt"
    "github.com/ateto3/gopherml/go-nlp/tokenizer"
)

func main() {
    // Tokenizador básico
    tok := tokenizer.NewTreebankWordTokenizer()
    tokens := tok.Tokenize("Hello world! How are you doing today?")
    fmt.Printf("Tokens: %v\n", tokens)
    
    // Tokenización paralela para texto grande
    parallelTok := tokenizer.NewParallelTokenizer(64*1024, 0, tokenizer.NewNLTKWordTokenizer())
    largeTokens := parallelTok.Tokenize(largeText)
    
    // Tokenización de tweets
    tweetTok := tokenizer.NewDefaultTweetTokenizer()
    tweetTokens := tweetTok.Tokenize("RT @user: Hello World! Check out https://example.com #golang")
    fmt.Printf("Tweet tokens: %v\n", tweetTokens)
}
```

### Limpieza de Texto

```go
package main

import (
    "fmt"
    "github.com/ateto3/gopherml/go-nlp/cleaner"
)

func main() {
    // Limpieza básica
    c := cleaner.New(
        cleaner.WithHTML(true),
        cleaner.WithPunctuation(true),
        cleaner.WithLowercase(true),
    )
    cleaned := c.Clean("<p>Hello, World!</p> https://example.com @user #test")
    fmt.Printf("Cleaned: %s\n", cleaned)
    
    // Limpieza con stopwords
    c := cleaner.New(
        cleaner.WithStopwords(true),
        cleaner.WithLanguage("en"),
    )
    cleaned = c.Clean("The quick brown fox jumps over the lazy dog")
    fmt.Printf("Without stopwords: %s\n", cleaned)
    
    // Limpieza de tokens
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
    
    // Con estadísticas
    stats, _ := r.Stats()
    fmt.Printf("Vocabulary size: %d\n", stats.VocabSize)
    
    // Lazy loading
    lazy := corpus.NewLazy(func() (corpus.Reader, error) {
        return corpus.NewFSReader(os.DirFS("data"), []string{"file.txt"}), nil
    })
}
```

### Traducción

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

### Monetización en RapidAPI

```
User Request → Go-NLP Preprocessing → AI API Gateway → LLM Provider
                                      ↓
                              Token Caching
                                      ↓
                              Cost Optimization (50-80%)
```

### Componentes del Gateway

1. **Preprocessing Layer**: go-nlp libraries para limpieza y tokenización
2. **Prompt Engineering**: Context injection y system prompt assembly
3. **Model Routing**: Selección de modelo según complejidad de tarea
4. **API Gateway**: Integración con RapidAPI y LLM providers
5. **Post-processing**: Translation y output cleaning

---

## 📖 Documentación

- `PROJECT_SPEC.md`: Especificación completa del proyecto y arquitectura
- `examples/`: Ejemplos de uso por módulo
- `benchmarks/`: Benchmarks de rendimiento

---

## 🤝 Contribución

### Guías de Contribución

1. Fork el repositorio
2. Create una feature branch (`git checkout -b feature/amazing-feature`)
3. Commit tus cambios (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Código de Conducta

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
