package probability

import (
	"testing"
)

func TestFreqDist(t *testing.T) {
	fd := NewFreqDist("a", "b", "a")
	if fd.Max() != "a" || fd.N() != 3 {
		t.Fatal(fd)
	}
}

func TestFreqDistEmpty(t *testing.T) {
	fd := NewFreqDist()
	if fd.N() != 0 {
		t.Errorf("expected N=0 for empty, got %d", fd.N())
	}
	if fd.B() != 0 {
		t.Errorf("expected B=0 for empty, got %d", fd.B())
	}
}

func TestFreqDistHapaxes(t *testing.T) {
	fd := NewFreqDist("a", "a", "b", "c", "c", "c", "d")
	hapaxes := fd.Hapaxes()
	if len(hapaxes) != 2 {
		t.Errorf("expected 2 hapaxes (b, d), got %d: %v", len(hapaxes), hapaxes)
	}
	// Verify b and d are in hapaxes
	hasB, hasD := false, false
	for _, h := range hapaxes {
		if h == "b" {
			hasB = true
		}
		if h == "d" {
			hasD = true
		}
	}
	if !hasB || !hasD {
		t.Errorf("expected 'b' and 'd' in hapaxes, got: %v", hapaxes)
	}
}

func TestFreqDistEmptyHapaxes(t *testing.T) {
	fd := NewFreqDist("a", "a", "b", "b")
	if len(fd.Hapaxes()) != 0 {
		t.Errorf("expected 0 hapaxes, got %d", len(fd.Hapaxes()))
	}
}

func TestFreqDistInc(t *testing.T) {
	fd := NewFreqDist("a")
	fd.Inc("a")
	fd.Inc("b")
	if fd.N() != 3 {
		t.Errorf("expected N=3, got %d", fd.N())
	}
	if fd["a"] != 2 {
		t.Errorf("expected 'a' count 2, got %d", fd["a"])
	}
}

func TestFreqDistB(t *testing.T) {
	fd := NewFreqDist("a", "a", "b", "c")
	if fd.B() != 3 {
		t.Errorf("expected B=3, got %d", fd.B())
	}
}

func TestFreqDistBEmpty(t *testing.T) {
	fd := NewFreqDist()
	if fd.B() != 0 {
		t.Errorf("expected B=0, got %d", fd.B())
	}
}

func TestFreqDistMax(t *testing.T) {
	fd := NewFreqDist("a", "a", "b", "c", "c")
	if fd.Max() != "a" {
		t.Errorf("expected 'a' (tied with 'c', first wins), got %s", fd.Max())
	}
}

func TestFreqDistMaxSingle(t *testing.T) {
	fd := NewFreqDist("hello")
	if fd.Max() != "hello" {
		t.Errorf("expected 'hello', got %s", fd.Max())
	}
}

func TestProbDist(t *testing.T) {
	fd := NewFreqDist("a", "a", "b")
	p := NewMLEProbDist(fd)
	if p.Prob("a") <= p.Prob("b") {
		t.Fatal("expected a > b")
	}
}

func TestMLEProbDist(t *testing.T) {
	fd := NewFreqDist("a", "a", "b", "b", "c")
	p := NewMLEProbDist(fd)
	if p.Prob("a") != 0.4 {
		t.Errorf("expected P(a)=0.4, got %f", p.Prob("a"))
	}
	if p.Prob("c") != 0.2 {
		t.Errorf("expected P(c)=0.2, got %f", p.Prob("c"))
	}
}

func TestMLEProbDistUnseen(t *testing.T) {
	fd := NewFreqDist("a", "b")
	p := NewMLEProbDist(fd)
	if p.Prob("z") != 0 {
		t.Errorf("expected 0 for unseen, got %f", p.Prob("z"))
	}
}

func TestMLEProbDistMax(t *testing.T) {
	fd := NewFreqDist("a", "a", "a", "b", "c")
	p := NewMLEProbDist(fd)
	if p.Max() != "a" {
		t.Errorf("expected max 'a', got %s", p.Max())
	}
}

func TestMLEProbDistEmpty(t *testing.T) {
	fd := NewFreqDist()
	p := NewMLEProbDist(fd)
	if p.Prob("a") != 0 {
		t.Errorf("expected 0 for empty, got %f", p.Prob("a"))
	}
}

func TestLidstoneProbDist(t *testing.T) {
	fd := NewFreqDist("a", "b")
	p := NewLidstoneProbDist(fd, 0.5, 2)
	// P(a) = (1+0.5)/(2+0.5*2) = 1.5/3 = 0.5
	if p.Prob("a") != 0.5 {
		t.Errorf("expected P(a)=0.5, got %f", p.Prob("a"))
	}
}

func TestLidstoneProbDistUnseen(t *testing.T) {
	fd := NewFreqDist("a", "b")
	p := NewLidstoneProbDist(fd, 0.5, 2)
	// P(unseen) = (0+0.5)/3 = 1/6
	expected := 1.0 / 6.0
	if p.Prob("z") != expected {
		t.Errorf("expected P(z)=%f, got %f", expected, p.Prob("z"))
	}
}

func TestLidstoneProbDistMax(t *testing.T) {
	fd := NewFreqDist("a", "a", "b")
	p := NewLidstoneProbDist(fd, 0.5, 2)
	if p.Max() != "a" {
		t.Errorf("expected max 'a', got %s", p.Max())
	}
}

func TestLidstoneProbDistDefaultBins(t *testing.T) {
	fd := NewFreqDist("a", "b")
	// bins=0 should default to fd.B()
	p := NewLidstoneProbDist(fd, 0.5, 0)
	if p.Prob("a") <= 0 {
		t.Error("expected positive probability")
	}
}

func TestProbDistInterface(t *testing.T) {
	var p ProbDist = NewMLEProbDist(NewFreqDist("a", "b"))
	if p.Prob("a") <= 0 {
		t.Error("expected positive probability through interface")
	}
	if p.Max() != "a" {
		t.Errorf("expected max 'a', got %s", p.Max())
	}
}

func TestNormalizeLogScores(t *testing.T) {
	scores := map[string]float64{"a": 0.0, "b": 0.0}
	out := NormalizeLogScores(scores)
	sum := 0.0
	for _, v := range out {
		sum += v
	}
	if sum < 0.99 || sum > 1.01 {
		t.Errorf("expected normalized scores to sum to ~1, got %f", sum)
	}
	// Equal scores should yield equal normalized values
	if out["a"] != out["b"] {
		t.Errorf("expected equal normalized values, got %f vs %f", out["a"], out["b"])
	}
}

func TestNormalizeLogScoresUnequal(t *testing.T) {
	scores := map[string]float64{"a": 0.0, "b": 1.0}
	out := NormalizeLogScores(scores)
	sum := 0.0
	for _, v := range out {
		sum += v
	}
	if sum < 0.99 || sum > 1.01 {
		t.Errorf("expected normalized scores to sum to ~1, got %f", sum)
	}
	// b should have higher probability than a
	if out["b"] <= out["a"] {
		t.Errorf("expected 'b' to have higher probability than 'a', got %f vs %f", out["b"], out["a"])
	}
}

func TestNormalizeLogScoresEmpty(t *testing.T) {
	out := NormalizeLogScores(map[string]float64{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestNewFreqDistConcurrent(t *testing.T) {
	samples := []string{"a", "a", "b", "b", "c", "c", "c"}
	fd := NewFreqDistConcurrent(samples)
	if fd.N() != len(samples) {
		t.Errorf("expected N=%d, got %d", len(samples), fd.N())
	}
	if fd["a"] != 2 || fd["b"] != 2 || fd["c"] != 3 {
		t.Errorf("unexpected counts: %v", fd)
	}
}

func TestNewFreqDistConcurrentEmpty(t *testing.T) {
	fd := NewFreqDistConcurrent(nil)
	if fd.N() != 0 {
		t.Errorf("expected N=0, got %d", fd.N())
	}
}

func TestFreqDistMapMutation(t *testing.T) {
	fd := NewFreqDist("a", "b", "c")
	if fd["a"] != 1 {
		t.Errorf("expected 'a' count 1, got %d", fd["a"])
	}
	fd["a"] = 10
	if fd["a"] != 10 {
		t.Errorf("expected 'a' count 10, got %d", fd["a"])
	}
}

func BenchmarkFreqDistInc(b *testing.B) {
	fd := NewFreqDist()
	words := []string{"a", "b", "c", "d", "e"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fd.Inc(words[i%len(words)])
	}
}

func BenchmarkNewMLEProbDist(b *testing.B) {
	fd := NewFreqDist("a", "a", "b", "b", "c", "c", "c")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewMLEProbDist(fd)
	}
}
