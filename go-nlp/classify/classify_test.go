package classify

import "testing"

func TestNaiveBayes(t *testing.T) {
	data := []LabeledFeatures{
		{Features: Features{"contains_fast": true, "len": 3}, Label: "tech"},
		{Features: Features{"contains_fast": true, "len": 4}, Label: "tech"},
		{Features: Features{"contains_fast": false, "len": 2}, Label: "sport"},
	}
	c := TrainNaiveBayes(data)
	if got := c.Classify(Features{"contains_fast": true, "len": 3}); got != "tech" {
		t.Fatalf("expected tech, got %s", got)
	}
}

func TestApplyFeatures(t *testing.T) {
	samples := []string{"a", "bb"}
	out := ApplyFeatures(samples, func(s string) Features { return Features{"len": len(s)} }, func(s string) string { return s })
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}
