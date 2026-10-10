// Package crypto provides high-level ML utilities tailored for cryptocurrency trading bots.
// It bridges raw OHLCV market data to GopherML feature vectors and prediction pipelines.
package crypto

import (
	"fmt"
	"math"
)

// OHLCV represents a single candlestick (Open, High, Low, Close, Volume).
// This is the universal input format used by most crypto exchanges.
type OHLCV struct {
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

// FeatureConfig controls which features are computed from an OHLCV window.
// All fields default to true when using NewFeatureConfig().
type FeatureConfig struct {
	// Returns: relative price changes between closes
	Returns bool
	// LogReturns: log(close[i]/close[i-1]) — more stationary than raw returns
	LogReturns bool
	// PriceRange: (high-low)/close — normalized volatility per candle
	PriceRange bool
	// VolumeRatio: volume[i]/avgVolume — relative volume vs window average
	VolumeRatio bool
	// Momentum: RSI-like directional momentum over the window
	Momentum bool
	// VWAP: Volume-Weighted Average Price deviation from close
	VWAP bool
}

// NewFeatureConfig returns the default feature configuration with all features enabled.
func NewFeatureConfig() FeatureConfig {
	return FeatureConfig{
		Returns:     true,
		LogReturns:  true,
		PriceRange:  true,
		VolumeRatio: true,
		Momentum:    true,
		VWAP:        true,
	}
}

// MinWindowSize returns the minimum number of candles required for the given config.
func MinWindowSize() int {
	return 3 // Need at least 3 candles for meaningful returns
}

// ToFeatureVector converts a sliding window of OHLCV candles into a flat feature vector
// suitable for GopherML models. The window should be at least MinWindowSize() candles.
//
// Feature layout (with defaults):
//   - returns[0..n-2]:      relative returns (n-1 values)
//   - log_returns[0..n-2]:  log returns (n-1 values)
//   - price_range[0..n-1]:  (high-low)/close per candle (n values)
//   - volume_ratio[0..n-1]: volume/avg_volume per candle (n values)
//   - momentum:             1 value (net up-moves / n-1)
//   - vwap_dev:             1 value (last close deviation from VWAP)
func ToFeatureVector(candles []OHLCV, cfg FeatureConfig) ([]float64, error) {
	n := len(candles)
	if n < MinWindowSize() {
		return nil, fmt.Errorf("need at least %d candles, got %d", MinWindowSize(), n)
	}

	// Validate all candle values
	for i, c := range candles {
		if c.Close <= 0 {
			return nil, fmt.Errorf("candle[%d].Close = %f: must be > 0", i, c.Close)
		}
		if c.High < c.Low {
			return nil, fmt.Errorf("candle[%d]: High (%f) < Low (%f)", i, c.High, c.Low)
		}
		if c.Volume < 0 {
			return nil, fmt.Errorf("candle[%d].Volume = %f: must be >= 0", i, c.Volume)
		}
		if math.IsNaN(c.Close) || math.IsInf(c.Close, 0) {
			return nil, fmt.Errorf("candle[%d].Close is NaN or Inf", i)
		}
	}

	var features []float64

	// --- Returns ---
	if cfg.Returns {
		for i := 1; i < n; i++ {
			ret := (candles[i].Close - candles[i-1].Close) / candles[i-1].Close
			features = append(features, ret)
		}
	}

	// --- Log Returns ---
	if cfg.LogReturns {
		for i := 1; i < n; i++ {
			lr := math.Log(candles[i].Close / candles[i-1].Close)
			features = append(features, lr)
		}
	}

	// --- Price Range (normalized volatility per candle) ---
	if cfg.PriceRange {
		for _, c := range candles {
			pr := (c.High - c.Low) / c.Close
			features = append(features, pr)
		}
	}

	// --- Volume Ratio ---
	if cfg.VolumeRatio {
		avgVol := 0.0
		for _, c := range candles {
			avgVol += c.Volume
		}
		avgVol /= float64(n)
		for _, c := range candles {
			ratio := 0.0
			if avgVol > 0 {
				ratio = c.Volume / avgVol
			}
			features = append(features, ratio)
		}
	}

	// --- Momentum (fraction of up-moves in the window) ---
	if cfg.Momentum {
		upMoves := 0.0
		for i := 1; i < n; i++ {
			if candles[i].Close > candles[i-1].Close {
				upMoves++
			}
		}
		features = append(features, upMoves/float64(n-1))
	}

	// --- VWAP Deviation ---
	if cfg.VWAP {
		totalVWAP := 0.0
		totalVol := 0.0
		for _, c := range candles {
			typicalPrice := (c.High + c.Low + c.Close) / 3.0
			totalVWAP += typicalPrice * c.Volume
			totalVol += c.Volume
		}
		vwapDev := 0.0
		if totalVol > 0 {
			vwap := totalVWAP / totalVol
			lastClose := candles[n-1].Close
			vwapDev = (lastClose - vwap) / vwap
		}
		features = append(features, vwapDev)
	}

	return features, nil
}

// FeatureSize returns the number of features produced by ToFeatureVector for a given window size.
// Use this to pre-allocate feature matrices.
func FeatureSize(windowSize int, cfg FeatureConfig) int {
	n := windowSize
	size := 0
	if cfg.Returns {
		size += n - 1
	}
	if cfg.LogReturns {
		size += n - 1
	}
	if cfg.PriceRange {
		size += n
	}
	if cfg.VolumeRatio {
		size += n
	}
	if cfg.Momentum {
		size += 1
	}
	if cfg.VWAP {
		size += 1
	}
	return size
}

// BuildFeatureMatrix converts a time series of OHLCV candles into a feature matrix
// using a sliding window approach. The resulting X matrix has shape:
//
//	[len(candles) - windowSize + 1][FeatureSize(windowSize, cfg)]
//
// This is the recommended way to prepare training data from exchange data.
func BuildFeatureMatrix(candles []OHLCV, windowSize int, cfg FeatureConfig) ([][]float64, error) {
	if len(candles) < windowSize {
		return nil, fmt.Errorf("need at least %d candles, got %d", windowSize, len(candles))
	}
	if windowSize < MinWindowSize() {
		return nil, fmt.Errorf("windowSize must be >= %d", MinWindowSize())
	}

	nWindows := len(candles) - windowSize + 1
	X := make([][]float64, nWindows)

	for i := 0; i < nWindows; i++ {
		window := candles[i : i+windowSize]
		fv, err := ToFeatureVector(window, cfg)
		if err != nil {
			return nil, fmt.Errorf("window[%d]: %w", i, err)
		}
		X[i] = fv
	}

	return X, nil
}
