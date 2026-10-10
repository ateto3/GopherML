package crypto

import (
	"fmt"
	"sync"

	"github.com/ateto3/gopherml/ensemble"
	"github.com/ateto3/gopherml/persistence"
	"github.com/ateto3/gopherml/preprocessing"
)

// Signal represents a trading signal produced by the predictor.
type Signal int

const (
	// SignalHold indicates no clear directional bias — stay out.
	SignalHold Signal = 0
	// SignalLong indicates the model predicts price will rise — go long.
	SignalLong Signal = 1
	// SignalShort indicates the model predicts price will fall — go short.
	SignalShort Signal = 2
)

// String returns a human-readable signal name.
func (s Signal) String() string {
	switch s {
	case SignalLong:
		return "LONG"
	case SignalShort:
		return "SHORT"
	default:
		return "HOLD"
	}
}

// PredictorConfig holds configuration for a CryptoPredictor.
type PredictorConfig struct {
	// NEstimators is the number of trees in the random forest (default: 100).
	NEstimators int
	// MaxDepth is the maximum depth of each tree (default: 10).
	MaxDepth int
	// WindowSize is the number of candles per feature window (default: 20).
	WindowSize int
	// MinConfidence is the minimum probability to emit a non-HOLD signal (default: 0.55).
	// Set higher (e.g., 0.70) for more conservative, fewer-but-higher-quality signals.
	MinConfidence float64
	// Features controls which feature groups are computed (default: all enabled).
	Features FeatureConfig
}

// DefaultPredictorConfig returns production-ready defaults for crypto bots.
func DefaultPredictorConfig() PredictorConfig {
	return PredictorConfig{
		NEstimators:   100,
		MaxDepth:      10,
		WindowSize:    20,
		MinConfidence: 0.55,
		Features:      NewFeatureConfig(),
	}
}

// CryptoPredictor is a production-ready, self-contained ML pipeline for crypto trading signals.
// It wraps a RandomForestClassifier with a StandardScaler and OHLCV feature engineering.
//
// Typical usage in a bot:
//
//	// At startup — load a pre-trained model:
//	p := crypto.NewCryptoPredictor(crypto.DefaultPredictorConfig())
//	err := p.LoadModel("btcusdt_rf.gob")
//
//	// In the trading loop — get a signal on each new candle:
//	signal, confidence, err := p.PredictSignal(recentCandles)
//	if err == nil && confidence > 0.65 {
//	    executeOrder(signal)
//	}
type CryptoPredictor struct {
	mu     sync.RWMutex
	cfg    PredictorConfig
	rf     *ensemble.RandomForestClassifier
	scaler *preprocessing.StandardScaler
	fitted bool
}

// NewCryptoPredictor creates a new CryptoPredictor with the given configuration.
func NewCryptoPredictor(cfg PredictorConfig) *CryptoPredictor {
	if cfg.NEstimators <= 0 {
		cfg.NEstimators = 100
	}
	if cfg.MaxDepth <= 0 {
		cfg.MaxDepth = 10
	}
	if cfg.WindowSize < MinWindowSize() {
		cfg.WindowSize = 20
	}
	if cfg.MinConfidence <= 0 {
		cfg.MinConfidence = 0.55
	}

	rf := ensemble.NewRandomForestClassifier()
	rf.NEstimators = cfg.NEstimators
	rf.MaxDepth = cfg.MaxDepth

	return &CryptoPredictor{
		cfg:    cfg,
		rf:     rf,
		scaler: preprocessing.NewStandardScaler(),
	}
}

// Fit trains the predictor on a labelled dataset of OHLCV candles.
//
// Parameters:
//   - candles: full time series of OHLCV candles (must have len >= windowSize+1)
//   - labels:  one label per sliding window: 0=hold, 1=long, 2=short.
//     len(labels) must equal len(candles) - windowSize + 1
//
// The method performs:
//  1. Feature extraction via BuildFeatureMatrix
//  2. StandardScaler fit+transform
//  3. RandomForest training
func (p *CryptoPredictor) Fit(candles []OHLCV, labels []float64) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	X, err := BuildFeatureMatrix(candles, p.cfg.WindowSize, p.cfg.Features)
	if err != nil {
		return fmt.Errorf("feature extraction failed: %w", err)
	}

	if len(X) != len(labels) {
		return fmt.Errorf("feature matrix has %d rows but %d labels provided (expected equal)", len(X), len(labels))
	}

	// Fit scaler and transform
	XScaled, err := p.scaler.FitTransform(X)
	if err != nil {
		return fmt.Errorf("scaler failed: %w", err)
	}

	// Train forest
	if err := p.rf.Fit(XScaled, labels); err != nil {
		return fmt.Errorf("forest training failed: %w", err)
	}

	p.fitted = true
	return nil
}

// PredictSignal predicts a trading signal from a window of recent OHLCV candles.
// The window must be exactly WindowSize candles long.
//
// Returns:
//   - signal: SignalHold, SignalLong, or SignalShort
//   - confidence: probability of the predicted class in [0, 1]
//   - err: non-nil if the model is not fitted or input is invalid
//
// A signal is only non-HOLD if confidence >= cfg.MinConfidence.
func (p *CryptoPredictor) PredictSignal(candles []OHLCV) (Signal, float64, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.fitted {
		return SignalHold, 0, fmt.Errorf("CryptoPredictor is not fitted — call Fit or LoadModel first")
	}
	if len(candles) != p.cfg.WindowSize {
		return SignalHold, 0, fmt.Errorf("expected %d candles, got %d", p.cfg.WindowSize, len(candles))
	}

	// Extract features
	fv, err := ToFeatureVector(candles, p.cfg.Features)
	if err != nil {
		return SignalHold, 0, fmt.Errorf("feature extraction: %w", err)
	}

	// Scale using fitted scaler
	scaled, err := p.scaler.Transform([][]float64{fv})
	if err != nil {
		return SignalHold, 0, fmt.Errorf("scaler: %w", err)
	}

	// Get class probabilities
	proba, err := p.rf.PredictProbaSingle(scaled[0])
	if err != nil {
		return SignalHold, 0, fmt.Errorf("prediction: %w", err)
	}

	// Find most confident class
	bestClass := 0
	bestProba := proba[0]
	for i := 1; i < len(proba); i++ {
		if proba[i] > bestProba {
			bestProba = proba[i]
			bestClass = i
		}
	}

	// Apply confidence threshold
	if bestProba < p.cfg.MinConfidence {
		return SignalHold, bestProba, nil
	}

	return Signal(bestClass), bestProba, nil
}

// predictorGob is the on-disk serialization struct for CryptoPredictor.
type predictorGob struct {
	RF     *ensemble.RandomForestClassifier
	Scaler *preprocessing.StandardScaler
	Cfg    PredictorConfig
	Fitted bool
}

// SaveModel saves the trained predictor to a .gob file.
// Both the scaler and forest are serialized atomically.
func (p *CryptoPredictor) SaveModel(path string) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.fitted {
		return fmt.Errorf("cannot save an unfitted CryptoPredictor")
	}

	data := predictorGob{
		RF:     p.rf,
		Scaler: p.scaler,
		Cfg:    p.cfg,
		Fitted: p.fitted,
	}
	return persistence.Save(path, &data)
}

// LoadModel loads a previously saved predictor from a .gob file.
// After loading, PredictSignal can be called immediately.
func (p *CryptoPredictor) LoadModel(path string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var data predictorGob
	if err := persistence.Load(path, &data); err != nil {
		return fmt.Errorf("failed to load model from %s: %w", path, err)
	}

	if data.RF == nil || data.Scaler == nil {
		return fmt.Errorf("loaded model from %s is corrupt (nil RF or Scaler)", path)
	}
	if !data.Fitted {
		return fmt.Errorf("loaded model from %s was not fitted", path)
	}

	p.rf = data.RF
	p.scaler = data.Scaler
	p.cfg = data.Cfg
	p.fitted = data.Fitted
	return nil
}

// IsFitted returns true if the predictor has been trained or loaded.
// Thread-safe.
func (p *CryptoPredictor) IsFitted() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.fitted
}

// Config returns the predictor configuration (read-only copy).
func (p *CryptoPredictor) Config() PredictorConfig {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.cfg
}
