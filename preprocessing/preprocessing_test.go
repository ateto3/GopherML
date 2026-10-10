package preprocessing

import (
	"math"
	"testing"
)

func TestStandardScaler(t *testing.T) {
	X := [][]float64{
		{1, 10},
		{2, 20},
		{3, 30},
	}
	// Mean: [2, 20], Std: [sqrt(2/3), sqrt(200/3)] approx [0.816, 8.16]

	scaler := NewStandardScaler()
	if err := scaler.Fit(X); err != nil {
		t.Fatalf("Fit error: %v", err)
	}

	if math.Abs(scaler.Mean[0]-2.0) > 1e-6 || math.Abs(scaler.Mean[1]-20.0) > 1e-6 {
		t.Errorf("Mean mismatch: %v", scaler.Mean)
	}

	XTrans, err := scaler.Transform(X)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	// First row: (1-2)/0.816 = -1.22, (10-20)/8.16 = -1.22
	if math.Abs(XTrans[0][0]+1.224744) > 1e-6 {
		t.Errorf("Transform mismatch: %f", XTrans[0][0])
	}
}

func TestStandardScalerFitTransform(t *testing.T) {
	X := [][]float64{{1, 2}, {3, 4}, {5, 6}}
	scaler := NewStandardScaler()
	XTrans, err := scaler.FitTransform(X)
	if err != nil {
		t.Fatalf("FitTransform error: %v", err)
	}
	if len(XTrans) != 3 || len(XTrans[0]) != 2 {
		t.Errorf("Unexpected shape: %v", XTrans)
	}
	// Middle row should be ~0 after standard scaling
	if math.Abs(XTrans[1][0]) > 1e-6 {
		t.Errorf("Expected center row to be ~0, got %f", XTrans[1][0])
	}
}

func TestStandardScalerTransformBeforeFit(t *testing.T) {
	scaler := NewStandardScaler()
	_, err := scaler.Transform([][]float64{{1.0, 2.0}})
	if err == nil {
		t.Error("Expected error when Transform called before Fit")
	}
}

func TestStandardScalerInverseTransform(t *testing.T) {
	X := [][]float64{{1.0, 10.0}, {2.0, 20.0}, {3.0, 30.0}}
	scaler := NewStandardScaler()
	XTrans, _ := scaler.FitTransform(X)
	XBack, err := scaler.InverseTransform(XTrans)
	if err != nil {
		t.Fatalf("InverseTransform error: %v", err)
	}
	if math.Abs(XBack[0][0]-1.0) > 1e-6 || math.Abs(XBack[2][1]-30.0) > 1e-6 {
		t.Errorf("InverseTransform mismatch: %v", XBack)
	}
}

func TestMinMaxScaler(t *testing.T) {
	X := [][]float64{
		{1, 10},
		{2, 20},
		{3, 30},
	}

	scaler := NewMinMaxScaler()
	if err := scaler.Fit(X); err != nil {
		t.Fatalf("Fit error: %v", err)
	}

	XTrans, err := scaler.Transform(X)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if XTrans[0][0] != 0.0 || XTrans[2][0] != 1.0 {
		t.Errorf("MinMax Transform mismatch: %v", XTrans)
	}
}

func TestMinMaxScalerFitTransform(t *testing.T) {
	X := [][]float64{{0.0, 0.0}, {5.0, 10.0}, {10.0, 20.0}}
	scaler := NewMinMaxScaler()
	XTrans, err := scaler.FitTransform(X)
	if err != nil {
		t.Fatalf("FitTransform error: %v", err)
	}
	if XTrans[0][0] != 0.0 || XTrans[2][0] != 1.0 {
		t.Errorf("Unexpected result: %v", XTrans)
	}
}

func TestMinMaxScalerTransformBeforeFit(t *testing.T) {
	scaler := NewMinMaxScaler()
	_, err := scaler.Transform([][]float64{{1.0, 2.0}})
	if err == nil {
		t.Error("Expected error when Transform called before Fit")
	}
}

func TestLabelEncoder(t *testing.T) {
	y := []string{"cat", "dog", "cat", "bird"}
	// Sorted: bird, cat, dog -> 0, 1, 2

	le := NewLabelEncoder()
	le.Fit(y)

	if len(le.Classes) != 3 || le.Classes[0] != "bird" {
		t.Errorf("Classes mismatch: %v", le.Classes)
	}

	yTrans, err := le.Transform(y)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if yTrans[0] != 1.0 || yTrans[1] != 2.0 || yTrans[3] != 0.0 {
		t.Errorf("Transform mismatch: %v", yTrans)
	}

	yInv, err := le.InverseTransform(yTrans)
	if err != nil {
		t.Fatalf("InverseTransform error: %v", err)
	}
	if yInv[0] != "cat" || yInv[3] != "bird" {
		t.Errorf("InverseTransform mismatch: %v", yInv)
	}
}

func TestLabelEncoderUnseenLabel(t *testing.T) {
	le := NewLabelEncoder()
	le.Fit([]string{"a", "b", "c"})

	_, err := le.Transform([]string{"a", "unknown"})
	if err == nil {
		t.Error("Expected error for unseen label, got nil")
	}
}

func TestLabelEncoderInverseTransformOutOfRange(t *testing.T) {
	le := NewLabelEncoder()
	le.Fit([]string{"a", "b"})

	_, err := le.InverseTransform([]float64{0, 5}) // 5 is out of range
	if err == nil {
		t.Error("Expected error for out-of-range index, got nil")
	}
}
