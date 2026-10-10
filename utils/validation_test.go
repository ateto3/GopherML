package utils

import (
	"math"
	"testing"
)

func TestCheckFinite_Valid(t *testing.T) {
	X := [][]float64{{1.0, 2.0}, {3.0, 4.0}}
	if err := CheckFinite(X); err != nil {
		t.Errorf("expected no error for valid data, got: %v", err)
	}
}

func TestCheckFinite_NaN(t *testing.T) {
	X := [][]float64{{1.0, math.NaN()}}
	if err := CheckFinite(X); err == nil {
		t.Error("expected error for NaN value, got nil")
	}
}

func TestCheckFinite_Inf(t *testing.T) {
	X := [][]float64{{math.Inf(1), 2.0}}
	if err := CheckFinite(X); err == nil {
		t.Error("expected error for +Inf value, got nil")
	}
	X2 := [][]float64{{math.Inf(-1), 2.0}}
	if err := CheckFinite(X2); err == nil {
		t.Error("expected error for -Inf value, got nil")
	}
}

func TestCheckFiniteVec_Valid(t *testing.T) {
	y := []float64{0.0, 1.0, -1.0, 100.0}
	if err := CheckFiniteVec(y); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestCheckFiniteVec_NaN(t *testing.T) {
	y := []float64{1.0, math.NaN(), 3.0}
	if err := CheckFiniteVec(y); err == nil {
		t.Error("expected error for NaN in y, got nil")
	}
}

func TestCheckFiniteSingle_Valid(t *testing.T) {
	x := []float64{0.1, 0.5, 99.9}
	if err := CheckFiniteSingle(x); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestCheckFiniteSingle_NaN(t *testing.T) {
	x := []float64{1.0, math.NaN()}
	if err := CheckFiniteSingle(x); err == nil {
		t.Error("expected error for NaN in single vector, got nil")
	}
}

func TestCheckXy_RejectsNaN(t *testing.T) {
	X := [][]float64{{1.0, math.NaN()}, {2.0, 3.0}}
	y := []float64{0.0, 1.0}
	if err := CheckXy(X, y); err == nil {
		t.Error("expected CheckXy to reject NaN in X, got nil")
	}
}

func TestCheckXy_RejectsNaNInY(t *testing.T) {
	X := [][]float64{{1.0, 2.0}, {3.0, 4.0}}
	y := []float64{0.0, math.NaN()}
	if err := CheckXy(X, y); err == nil {
		t.Error("expected CheckXy to reject NaN in y, got nil")
	}
}
