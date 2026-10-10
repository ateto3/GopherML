package sentiment

import "github.com/ateto3/gopherml/go-nlp/classify"

// SentimentAnalyzer wraps a classifier for sentiment tasks.
type SentimentAnalyzer struct {
	classifier classify.Classifier
}

// NewSentimentAnalyzer creates a classifier-backed analyzer.
func NewSentimentAnalyzer(c classify.Classifier) *SentimentAnalyzer {
	return &SentimentAnalyzer{classifier: c}
}

// Analyze classifies the sentiment of text.
func (s *SentimentAnalyzer) Analyze(text string) Score {
	feat := classify.Features{"text": text}
	label := s.classifier.Classify(feat)
	switch label {
	case "pos":
		return Score{Pos: 1, Comp: 1}
	case "neg":
		return Score{Neg: 1, Comp: -1}
	default:
		return Score{Neu: 1}
	}
}
