package main

import (
	"fmt"
	"log"

	"github.com/ateto3/gopherml/linear_model"
	"github.com/ateto3/gopherml/model_selection"
)

func main() {
	fmt.Println("=== Hyperparameter Tuning (GridSearchCV) ===")

	// Dataset binario linealmente separable
	X := [][]float64{
		{1.0, 1.1}, {1.4, 0.9}, {1.2, 1.3}, {0.9, 1.0},
		{1.1, 1.2}, {1.3, 1.0}, {0.8, 1.4}, {1.2, 0.8},
		{4.8, 5.1}, {5.2, 4.9}, {5.0, 5.3}, {4.7, 5.0},
		{5.1, 4.8}, {4.9, 5.2}, {5.3, 5.0}, {4.6, 5.1},
	}
	y := []float64{0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1}

	// Fábrica: construye un modelo nuevo para cada combinación de parámetros
	factory := func(params map[string]interface{}) model_selection.Estimator {
		lr := linear_model.NewLogisticRegression()
		lr.LearningRate = params["lr"].(float64)
		lr.MaxIter = params["iter"].(int)
		return lr
	}

	gs := &model_selection.GridSearchCV{
		EstimatorFactory: factory,
		ParamGrid: []map[string]interface{}{
			{"lr": 0.01, "iter": 100},
			{"lr": 0.10, "iter": 200},
			{"lr": 0.50, "iter": 300},
		},
		CV: model_selection.KFold{
			NSplits:     4,
			Shuffle:     true,
			RandomState: 42,
		},
		NJobs: -1, // paralelo: todos los núcleos
	}

	if err := gs.Fit(X, y); err != nil {
		log.Fatalf("GridSearchCV: %v", err)
	}

	fmt.Println("\nResultados por combinación:")
	for _, r := range gs.CVResults {
		if r.Err != nil {
			fmt.Printf("  %v → error: %v\n", r.Params, r.Err)
			continue
		}
		fmt.Printf("  %-32v → score: %.4f\n", r.Params, r.MeanScore)
	}

	fmt.Printf("\nMejores parámetros: %v\n", gs.BestParams)
	fmt.Printf("Mejor score (CV):   %.4f\n", gs.BestScore)

	// Entrenar el modelo final con los mejores parámetros
	best := gs.EstimatorFactory(gs.BestParams)
	if err := best.Fit(X, y); err != nil {
		log.Fatalf("Fit final: %v", err)
	}
	score, err := best.Score(X, y)
	if err != nil {
		log.Fatalf("Score: %v", err)
	}
	fmt.Printf("Score del modelo final: %.4f\n", score)

	fmt.Println("\n=== Example Complete ===")
}
