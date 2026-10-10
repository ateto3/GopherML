package crypto

import (
	"math"
	"os"
	"testing"
)

// --- OHLCV Feature Tests ---

func TestToFeatureVector_Basic(t *testing.T) {
	candles := []OHLCV{
		{Open: 100, High: 105, Low: 98, Close: 103, Volume: 1000},
		{Open: 103, High: 108, Low: 101, Close: 106, Volume: 1200},
		{Open: 106, High: 110, Low: 104, Close: 108, Volume: 900},
		{Open: 108, High: 112, Low: 106, Close: 107, Volume: 1100},
		{Open: 107, High: 109, Low: 103, Close: 105, Volume: 800},
	}

	cfg := NewFeatureConfig()
	fv, err := ToFeatureVector(candles, cfg)
	if err != nil {
		t.Fatalf("ToFeatureVector failed: %v", err)
	}

	expectedSize := FeatureSize(len(candles), cfg)
	if len(fv) != expectedSize {
		t.Errorf("Expected %d features, got %d", expectedSize, len(fv))
	}

	// All features must be finite
	for i, v := range fv {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			t.Errorf("Feature[%d] is NaN or Inf", i)
		}
	}
}

func TestToFeatureVector_TooFewCandles(t *testing.T) {
	candles := []OHLCV{
		{Open: 100, High: 105, Low: 98, Close: 103, Volume: 1000},
		{Open: 103, High: 108, Low: 101, Close: 106, Volume: 1200},
	}
	_, err := ToFeatureVector(candles, NewFeatureConfig())
	if err == nil {
		t.Error("Expected error for fewer than MinWindowSize candles, got nil")
	}
}

func TestToFeatureVector_InvalidClose(t *testing.T) {
	candles := []OHLCV{
		{Open: 100, High: 105, Low: 98, Close: 0, Volume: 1000}, // Close=0 is invalid
		{Open: 103, High: 108, Low: 101, Close: 106, Volume: 1200},
		{Open: 106, High: 110, Low: 104, Close: 108, Volume: 900},
	}
	_, err := ToFeatureVector(candles, NewFeatureConfig())
	if err == nil {
		t.Error("Expected error for Close=0, got nil")
	}
}

func TestFeatureSize(t *testing.T) {
	cfg := NewFeatureConfig()
	n := 10
	size := FeatureSize(n, cfg)
	// returns: n-1, log_returns: n-1, price_range: n, volume_ratio: n, momentum: 1, vwap: 1
	expected := (n-1)*2 + n*2 + 2
	if size != expected {
		t.Errorf("Expected FeatureSize=%d, got %d", expected, size)
	}
}

func TestBuildFeatureMatrix(t *testing.T) {
	n := 50
	candles := makeSyntheticCandles(n)

	windowSize := 10
	cfg := NewFeatureConfig()
	X, err := BuildFeatureMatrix(candles, windowSize, cfg)
	if err != nil {
		t.Fatalf("BuildFeatureMatrix failed: %v", err)
	}

	expectedRows := n - windowSize + 1
	if len(X) != expectedRows {
		t.Errorf("Expected %d rows, got %d", expectedRows, len(X))
	}

	expectedCols := FeatureSize(windowSize, cfg)
	for i, row := range X {
		if len(row) != expectedCols {
			t.Errorf("Row[%d] has %d cols, expected %d", i, len(row), expectedCols)
		}
	}
}

// --- CryptoPredictor Tests ---

func TestCryptoPredictorFitPredict(t *testing.T) {
	cfg := DefaultPredictorConfig()
	cfg.NEstimators = 10 // Fast for tests
	cfg.WindowSize = 5

	p := NewCryptoPredictor(cfg)

	// Build synthetic data: 100 candles
	candles := makeSyntheticCandles(100)

	// Build labels for sliding window: nCandles - windowSize + 1 labels
	nLabels := len(candles) - cfg.WindowSize + 1
	labels := make([]float64, nLabels)
	for i := range labels {
		if i%3 == 0 {
			labels[i] = 1 // LONG
		} else if i%3 == 1 {
			labels[i] = 2 // SHORT
		}
		// else 0 = HOLD
	}

	if err := p.Fit(candles, labels); err != nil {
		t.Fatalf("Fit failed: %v", err)
	}

	if !p.IsFitted() {
		t.Error("Expected IsFitted=true after Fit")
	}

	// PredictSignal on the last window
	lastWindow := candles[len(candles)-cfg.WindowSize:]
	signal, confidence, err := p.PredictSignal(lastWindow)
	if err != nil {
		t.Fatalf("PredictSignal failed: %v", err)
	}

	if confidence < 0 || confidence > 1 {
		t.Errorf("Confidence out of [0,1]: %f", confidence)
	}
	_ = signal // Signal may be HOLD — that's valid behavior
}

func TestCryptoPredictorSaveLoad(t *testing.T) {
	cfg := DefaultPredictorConfig()
	cfg.NEstimators = 5
	cfg.WindowSize = 5

	p := NewCryptoPredictor(cfg)
	candles := makeSyntheticCandles(60)
	nLabels := len(candles) - cfg.WindowSize + 1
	labels := make([]float64, nLabels)
	for i := range labels {
		labels[i] = float64(i % 3)
	}

	if err := p.Fit(candles, labels); err != nil {
		t.Fatalf("Fit failed: %v", err)
	}

	// Save
	tmpFile := "test_crypto_predictor.gob"
	defer os.Remove(tmpFile)
	if err := p.SaveModel(tmpFile); err != nil {
		t.Fatalf("SaveModel failed: %v", err)
	}

	// Load into a fresh predictor
	p2 := NewCryptoPredictor(cfg)
	if err := p2.LoadModel(tmpFile); err != nil {
		t.Fatalf("LoadModel failed: %v", err)
	}

	if !p2.IsFitted() {
		t.Error("Expected loaded predictor to be fitted")
	}

	// Both predictors should produce the same signal
	lastWindow := candles[len(candles)-cfg.WindowSize:]
	s1, c1, _ := p.PredictSignal(lastWindow)
	s2, c2, _ := p2.PredictSignal(lastWindow)

	if s1 != s2 {
		t.Errorf("Signals differ after save/load: original=%v loaded=%v", s1, s2)
	}
	if math.Abs(c1-c2) > 1e-9 {
		t.Errorf("Confidence differs after save/load: original=%f loaded=%f", c1, c2)
	}
}

func TestCryptoPredictorNotFittedError(t *testing.T) {
	p := NewCryptoPredictor(DefaultPredictorConfig())
	candles := makeSyntheticCandles(20)
	_, _, err := p.PredictSignal(candles[:20])
	if err == nil {
		t.Error("Expected error from unfitted predictor, got nil")
	}
}

func TestCryptoPredictorWrongWindowSize(t *testing.T) {
	cfg := DefaultPredictorConfig()
	cfg.NEstimators = 5
	cfg.WindowSize = 10

	p := NewCryptoPredictor(cfg)
	candles := makeSyntheticCandles(50)
	nLabels := len(candles) - cfg.WindowSize + 1
	labels := make([]float64, nLabels)
	p.Fit(candles, labels)

	// Pass wrong number of candles
	_, _, err := p.PredictSignal(candles[:5]) // only 5, expected 10
	if err == nil {
		t.Error("Expected error for wrong window size, got nil")
	}
}

func TestSignalString(t *testing.T) {
	if SignalHold.String() != "HOLD" {
		t.Errorf("Expected HOLD, got %s", SignalHold.String())
	}
	if SignalLong.String() != "LONG" {
		t.Errorf("Expected LONG, got %s", SignalLong.String())
	}
	if SignalShort.String() != "SHORT" {
		t.Errorf("Expected SHORT, got %s", SignalShort.String())
	}
}

// --- Helpers ---

// makeSyntheticCandles generates a realistic-looking synthetic OHLCV time series.
func makeSyntheticCandles(n int) []OHLCV {
	candles := make([]OHLCV, n)
	price := 50000.0 // BTC-like starting price
	for i := 0; i < n; i++ {
		// Simple pseudo-random walk (deterministic)
		delta := math.Sin(float64(i)*0.3) * 500
		open := price
		close_ := price + delta
		high := math.Max(open, close_) + math.Abs(delta)*0.2
		low := math.Min(open, close_) - math.Abs(delta)*0.2
		vol := 100.0 + math.Abs(delta)*0.1

		if close_ <= 0 {
			close_ = 100
		}
		if high < close_ {
			high = close_ + 1
		}
		if low > close_ {
			low = close_ - 1
		}
		if low <= 0 {
			low = 1
		}

		candles[i] = OHLCV{
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close_,
			Volume: vol,
		}
		price = close_
	}
	return candles
}
