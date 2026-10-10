package model_selection

import (
	"testing"

	"github.com/ateto3/gopherml/linear_model"
)

func TestGridSearchCVParallel(t *testing.T) {
	X := [][]float64{{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}, {6, 6}}
	y := []float64{0, 0, 0, 1, 1, 1}

	paramGrid := []map[string]interface{}{
		{"lr": 0.01, "iter": 100},
		{"lr": 0.1, "iter": 200},
		{"lr": 0.5, "iter": 300},
	}

	factory := func(params map[string]interface{}) Estimator {
		lr := linear_model.NewLogisticRegression()
		lr.LearningRate = params["lr"].(float64)
		lr.MaxIter = params["iter"].(int)
		return lr
	}

	gs := &GridSearchCV{
		EstimatorFactory: factory,
		ParamGrid:        paramGrid,
		CV:               KFold{NSplits: 2, Shuffle: true, RandomState: 42},
		NJobs:            -1, // Parallel!
	}

	err := gs.Fit(X, y)
	if err != nil {
		t.Fatalf("GridSearchCV failed: %v", err)
	}

	if gs.BestScore == 0 {
		t.Errorf("Expected best score > 0, got %f", gs.BestScore)
	}

	t.Logf("Best params: %v, Best score: %f", gs.BestParams, gs.BestScore)
}
