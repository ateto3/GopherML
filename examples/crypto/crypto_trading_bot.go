package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/ateto3/gopherml/crypto"
)

const modelPath = "crypto_example_model.gob"

// syntheticCandles genera una serie OHLCV sintética con régimen cíclico (reproducible).
func syntheticCandles(n int) []crypto.OHLCV {
	rng := rand.New(rand.NewSource(42))
	candles := make([]crypto.OHLCV, n)
	price := 100.0
	for i := range candles {
		regime := math.Sin(float64(i) / 8.0) // ciclos alcistas/bajistas
		open := price
		close := open + regime*0.3 + (rng.Float64()-0.5)*0.6
		high := math.Max(open, close) + rng.Float64()*0.5
		low := math.Min(open, close) - rng.Float64()*0.5
		volume := 1000 + rng.Float64()*500
		candles[i] = crypto.OHLCV{Open: open, High: high, Low: low, Close: close, Volume: volume}
		price = close
	}
	return candles
}

func signalName(s crypto.Signal) string {
	switch s {
	case crypto.SignalLong:
		return "LONG"
	case crypto.SignalShort:
		return "SHORT"
	default:
		return "HOLD"
	}
}

func main() {
	fmt.Println("=== Crypto Trading Bot Example ===")

	// 1. Configurar el predictor
	cfg := crypto.DefaultPredictorConfig()
	cfg.WindowSize = 10
	cfg.NEstimators = 50
	cfg.MaxDepth = 6
	cfg.MinConfidence = 0.55

	// 2. Datos: velas sintéticas (en producción vendrían del exchange)
	candles := syntheticCandles(500)
	fmt.Println("Velas generadas:", len(candles))

	// 3. Etiquetas por ventana deslizante: 0=hold, 1=long, 2=short
	//    (según el movimiento de la vela siguiente a la ventana)
	nLabels := len(candles) - cfg.WindowSize + 1
	labels := make([]float64, nLabels)
	for i := 0; i < nLabels-1; i++ {
		last := candles[i+cfg.WindowSize-1]
		next := candles[i+cfg.WindowSize]
		switch {
		case next.Close > last.Close:
			labels[i] = 1
		case next.Close < last.Close:
			labels[i] = 2
		}
		// última ventana: sin vela siguiente → 0 (hold)
	}

	// 4. Entrenar
	p := crypto.NewCryptoPredictor(cfg)
	if err := p.Fit(candles, labels); err != nil {
		log.Fatalf("Fit: %v", err)
	}
	fmt.Println("Modelo entrenado (modelo + scaler interno)")

	// 5. Guardar y recargar (persistencia de todo el pipeline)
	if err := p.SaveModel(modelPath); err != nil {
		log.Fatalf("SaveModel: %v", err)
	}
	defer os.Remove(modelPath)

	reloaded := crypto.NewCryptoPredictor(cfg)
	if err := reloaded.LoadModel(modelPath); err != nil {
		log.Fatalf("LoadModel: %v", err)
	}
	fmt.Println("Modelo recargado desde", modelPath)

	// 6. Loop de trading: una predicción por nueva vela
	window := candles[len(candles)-cfg.WindowSize:]
	signal, confidence, err := reloaded.PredictSignal(window)
	if err != nil {
		log.Fatalf("PredictSignal: %v", err)
	}

	fmt.Printf("\nSeñal: %s | Confianza: %.1f%%\n", signalName(signal), confidence*100)
	switch signal {
	case crypto.SignalLong:
		fmt.Println("→ Acción: comprar")
	case crypto.SignalShort:
		fmt.Println("→ Acción: vender / abrir corto")
	default:
		fmt.Println("→ Acción: no operar (umbral de confianza no superado)")
	}

	fmt.Println("\n=== Example Complete ===")
}
