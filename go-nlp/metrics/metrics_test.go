package metrics

import (
	"math"
	"testing"
)

func TestMetrics(t *testing.T) {
	ref := []string{"A", "B", "A"}
	test := []string{"A", "A", "A"}
	if Accuracy(ref, test) <= 0 {
		t.Fatal("acc")
	}
	cm := NewConfusionMatrix(ref, test)
	if cm.Count("A", "A") != 2 {
		t.Fatal("cm")
	}
	if CohenKappa(ref, ref) != 1 {
		t.Fatal("kappa")
	}
}

func TestAccuracyPerfect(t *testing.T) {
	ref := []string{"A", "B", "A", "B"}
	test := []string{"A", "B", "A", "B"}
	got := Accuracy(ref, test)
	if got != 1.0 {
		t.Errorf("expected accuracy 1.0, got %f", got)
	}
}

func TestAccuracyZero(t *testing.T) {
	ref := []string{"A", "A", "A"}
	test := []string{"B", "B", "B"}
	got := Accuracy(ref, test)
	if got != 0.0 {
		t.Errorf("expected accuracy 0.0, got %f", got)
	}
}

func TestAccuracyPartial(t *testing.T) {
	ref := []string{"A", "B", "A", "B"}
	test := []string{"A", "A", "A", "A"}
	got := Accuracy(ref, test)
	if got != 0.5 {
		t.Errorf("expected accuracy 0.5, got %f", got)
	}
}

func TestAccuracyEmpty(t *testing.T) {
	got := Accuracy(nil, nil)
	if got != 0 {
		t.Errorf("expected 0 for empty inputs, got %f", got)
	}
}

func TestAccuracyLengthMismatch(t *testing.T) {
	ref := []string{"A", "B"}
	test := []string{"A"}
	got := Accuracy(ref, test)
	if got != 0 {
		t.Errorf("expected 0 for length mismatch, got %f", got)
	}
}

func TestPrecisionRecall(t *testing.T) {
	ref := map[string]struct{}{"a": {}, "b": {}, "c": {}}
	test := map[string]struct{}{"a": {}, "b": {}}
	p := Precision(ref, test)
	if p != 1.0 {
		t.Errorf("expected precision 1.0, got %f", p)
	}
	r := Recall(ref, test)
	expected := 2.0 / 3.0
	if r != expected {
		t.Errorf("expected recall %f, got %f", expected, r)
	}
}

func TestPrecisionPartial(t *testing.T) {
	ref := map[string]struct{}{"a": {}, "b": {}}
	test := map[string]struct{}{"a": {}, "c": {}, "d": {}}
	p := Precision(ref, test)
	expected := 1.0 / 3.0
	if p != expected {
		t.Errorf("expected precision %f, got %f", expected, p)
	}
}

func TestPrecisionEmpty(t *testing.T) {
	ref := map[string]struct{}{"a": {}}
	test := map[string]struct{}{}
	if Precision(ref, test) != 0 {
		t.Error("expected 0 precision for empty test")
	}
}

func TestRecallEmpty(t *testing.T) {
	ref := map[string]struct{}{}
	test := map[string]struct{}{"a": {}}
	if Recall(ref, test) != 0 {
		t.Error("expected 0 recall for empty ref")
	}
}

func TestFMeasure(t *testing.T) {
	ref := map[string]struct{}{"a": {}, "b": {}}
	test := map[string]struct{}{"a": {}, "b": {}}
	f := FMeasure(ref, test)
	if f != 1.0 {
		t.Errorf("expected F1 1.0, got %f", f)
	}
}

func TestFMeasureZero(t *testing.T) {
	ref := map[string]struct{}{"a": {}}
	test := map[string]struct{}{}
	f := FMeasure(ref, test)
	if f != 0 {
		t.Errorf("expected F1 0 for empty test, got %f", f)
	}
}

func TestFMeasurePartial(t *testing.T) {
	ref := map[string]struct{}{"a": {}, "b": {}}
	test := map[string]struct{}{"a": {}}
	f := FMeasure(ref, test)
	if math.Abs(f-2.0/3.0) > 1e-10 {
		t.Errorf("expected F1 ~0.667, got %f", f)
	}
}

func TestNewConfusionMatrix(t *testing.T) {
	ref := []string{"A", "B", "A"}
	test := []string{"A", "A", "A"}
	cm := NewConfusionMatrix(ref, test)
	if len(cm.Values) != 2 {
		t.Errorf("expected 2 unique values, got %d", len(cm.Values))
	}
	if cm.Count("A", "A") != 2 {
		t.Errorf("expected A->A count 2, got %d", cm.Count("A", "A"))
	}
	if cm.Count("B", "A") != 1 {
		t.Errorf("expected B->A count 1, got %d", cm.Count("B", "A"))
	}
}

func TestConfusionMatrixPrecision(t *testing.T) {
	ref := []string{"A", "A", "B"}
	test := []string{"A", "B", "A"}
	cm := NewConfusionMatrix(ref, test)
	p := cm.Precision("A")
	if p != 0.5 {
		t.Errorf("expected precision(A) 0.5, got %f", p)
	}
}

func TestConfusionMatrixRecall(t *testing.T) {
	ref := []string{"A", "A", "B"}
	test := []string{"A", "B", "A"}
	cm := NewConfusionMatrix(ref, test)
	r := cm.Recall("A")
	if r != 0.5 {
		t.Errorf("expected recall(A) 0.5, got %f", r)
	}
}

func TestConfusionMatrixFMeasure(t *testing.T) {
	ref := []string{"A", "B", "A", "B"}
	test := []string{"A", "A", "A", "A"}
	cm := NewConfusionMatrix(ref, test)
	f := cm.FMeasure("A")
	if math.Abs(f-2.0/3.0) > 1e-10 {
		t.Errorf("expected F1 ~0.667, got %f", f)
	}
}

func TestConfusionMatrixEmpty(t *testing.T) {
	cm := NewConfusionMatrix([]string{}, []string{})
	if len(cm.Values) != 0 {
		t.Errorf("expected 0 values, got %d", len(cm.Values))
	}
}

func TestCohenKappaPerfect(t *testing.T) {
	a := []string{"A", "B", "A", "B"}
	if CohenKappa(a, a) != 1.0 {
		t.Errorf("expected kappa 1.0 for identical inputs, got %f", CohenKappa(a, a))
	}
}

func TestCohenKappaEqualAgreement(t *testing.T) {
	a := []string{"A", "A", "A", "A"}
	b := []string{"B", "B", "B", "B"}
	// po=0, pe=0 -> kappa=0 (no better than chance)
	k := CohenKappa(a, b)
	if k != 0.0 {
		t.Errorf("expected kappa 0.0 for total disagreement, got %f", k)
	}
}

func TestCohenKappaEmpty(t *testing.T) {
	k := CohenKappa(nil, nil)
	if k != 0 {
		t.Errorf("expected 0 for empty, got %f", k)
	}
}

func TestCohenKappaLengthMismatch(t *testing.T) {
	a := []string{"A", "B"}
	b := []string{"A"}
	k := CohenKappa(a, b)
	if k != 0 {
		t.Errorf("expected 0 for mismatch, got %f", k)
	}
}

func TestConfusionMatrixSingleClass(t *testing.T) {
	cm := NewConfusionMatrix([]string{"A", "A"}, []string{"A", "A"})
	if cm.Precision("A") != 1.0 {
		t.Errorf("expected precision 1.0, got %f", cm.Precision("A"))
	}
	if cm.Recall("A") != 1.0 {
		t.Errorf("expected recall 1.0, got %f", cm.Recall("A"))
	}
}

func TestConfusionMatrixSortedValues(t *testing.T) {
	cm := NewConfusionMatrix([]string{"Z", "A"}, []string{"M", "A"})
	if len(cm.Values) != 3 {
		t.Errorf("expected 3 unique values, got %d", len(cm.Values))
	}
	if cm.Values[0] != "A" || cm.Values[1] != "M" || cm.Values[2] != "Z" {
		t.Errorf("expected sorted values [A, M, Z], got %v", cm.Values)
	}
}

func TestCohenKappaModerateAgreement(t *testing.T) {
	a := []string{"A", "A", "A", "B", "B"}
	b := []string{"A", "A", "B", "B", "B"}
	k := CohenKappa(a, b)
	if k <= 0 {
		t.Errorf("expected positive kappa, got %f", k)
	}
}

func BenchmarkAccuracy(b *testing.B) {
	ref := make([]string, 1000)
	test := make([]string, 1000)
	for i := range ref {
		ref[i] = "A"
		if i%2 == 0 {
			test[i] = "A"
		} else {
			test[i] = "B"
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Accuracy(ref, test)
	}
}

func BenchmarkCohenKappa(b *testing.B) {
	ref := make([]string, 100)
	test := make([]string, 100)
	for i := range ref {
		ref[i] = "A"
		test[i] = "A"
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CohenKappa(ref, test)
	}
}
