package sentiment

import (
	"testing"

	"github.com/ateto3/gopherml/go-nlp/classify"
)

func TestVaderAnalyzer(t *testing.T) {
	s := NewVaderAnalyzer().Analyze("this is very good")
	if s.Comp <= 0 {
		t.Fatalf("expected positive score, got %+v", s)
	}
}

func TestSentimentAnalyzer(t *testing.T) {
	c := classify.TrainMajority([]classify.LabeledFeatures{
		{Features: classify.Features{"text": "good"}, Label: "pos"},
		{Features: classify.Features{"text": "bad"}, Label: "neg"},
		{Features: classify.Features{"text": "good"}, Label: "pos"},
	})
	sa := NewSentimentAnalyzer(c)
	if got := sa.Analyze("anything"); got.Pos != 1 {
		t.Fatalf("expected positive default, got %+v", got)
	}
}

func TestVaderPositive(t *testing.T) {
	s := NewVaderAnalyzer().Analyze("I love this amazing product")
	if s.Comp <= 0 {
		t.Errorf("expected positive score for positive text, got %+v", s)
	}
	if s.Pos <= 0 {
		t.Errorf("expected Pos > 0, got %f", s.Pos)
	}
}

func TestVaderNegative(t *testing.T) {
	s := NewVaderAnalyzer().Analyze("this is terrible and awful")
	if s.Comp >= 0 {
		t.Errorf("expected negative score for negative text, got %+v", s)
	}
	if s.Neg <= 0 {
		t.Errorf("expected Neg > 0, got %f", s.Neg)
	}
}

func TestVaderNeutral(t *testing.T) {
	s := NewVaderAnalyzer().Analyze("the chair is in the room")
	if s.Neu < 0.5 {
		t.Errorf("expected high Neu for neutral text, got %+v", s)
	}
}

func TestVaderBooster(t *testing.T) {
	weak := NewVaderAnalyzer().Analyze("good")
	strong := NewVaderAnalyzer().Analyze("very good")
	if strong.Comp <= weak.Comp {
		t.Errorf("expected 'very good' to have higher score than 'good', got %f vs %f", strong.Comp, weak.Comp)
	}
}

func TestVaderNegation(t *testing.T) {
	positive := NewVaderAnalyzer().Analyze("good")
	negated := NewVaderAnalyzer().Analyze("not good")
	if negated.Comp >= positive.Comp {
		t.Errorf("expected 'not good' to have lower score than 'good', got %f vs %f", negated.Comp, positive.Comp)
	}
}

func TestVaderEmpty(t *testing.T) {
	s := NewVaderAnalyzer().Analyze("")
	if s.Comp != 0 {
		t.Errorf("expected Comp=0 for empty, got %f", s.Comp)
	}
}

func TestSentimentAnalyzerPositive(t *testing.T) {
	c := classify.TrainMajority([]classify.LabeledFeatures{
		{Features: classify.Features{"text": "happy"}, Label: "pos"},
	})
	sa := NewSentimentAnalyzer(c)
	got := sa.Analyze("any text")
	if got.Pos != 1 {
		t.Errorf("expected Pos=1, got %+v", got)
	}
}

func TestSentimentAnalyzerNegative(t *testing.T) {
	c := classify.TrainMajority([]classify.LabeledFeatures{
		{Features: classify.Features{"text": "sad"}, Label: "neg"},
	})
	sa := NewSentimentAnalyzer(c)
	got := sa.Analyze("text")
	if got.Neg != 1 {
		t.Errorf("expected Neg=1, got %+v", got)
	}
}

func TestScoreStruct(t *testing.T) {
	s := Score{Neg: 0.2, Neu: 0.5, Pos: 0.3, Comp: 0.5}
	if s.Neg != 0.2 || s.Neu != 0.5 || s.Pos != 0.3 || s.Comp != 0.5 {
		t.Errorf("unexpected Score: %+v", s)
	}
}

func TestAnalyzerInterface(t *testing.T) {
	var _ Analyzer = NewVaderAnalyzer()
}

func TestNewVaderAnalyzer(t *testing.T) {
	v := NewVaderAnalyzer()
	s := v.Analyze("good")
	_ = s
}

func TestNewSentimentAnalyzer(t *testing.T) {
	c := classify.TrainMajority([]classify.LabeledFeatures{
		{Features: classify.Features{"text": "x"}, Label: "pos"},
	})
	sa := NewSentimentAnalyzer(c)
	if sa == nil {
		t.Fatal("NewSentimentAnalyzer returned nil")
	}
}

func TestVaderExtremeBooster(t *testing.T) {
	very := NewVaderAnalyzer().Analyze("very good")
	extremely := NewVaderAnalyzer().Analyze("extremely good")
	if extremely.Comp <= very.Comp {
		t.Errorf("expected 'extremely' to amplify more than 'very', got %f vs %f", extremely.Comp, very.Comp)
	}
}

func TestVaderScoreComponents(t *testing.T) {
	s := NewVaderAnalyzer().Analyze("love this")
	total := s.Pos + s.Neg + s.Neu
	if total < 0 || total > 2 {
		t.Errorf("score components seem invalid: pos=%f neg=%f neu=%f", s.Pos, s.Neg, s.Neu)
	}
}

func BenchmarkVaderAnalyzer(b *testing.B) {
	v := NewVaderAnalyzer()
	text := "I really love this amazing product, it's the best thing ever"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Analyze(text)
	}
}
