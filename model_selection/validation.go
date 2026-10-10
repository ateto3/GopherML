package model_selection

// Estimator is the interface for models that can be evaluated with CrossValScore.
type Estimator interface {
	Fit(X [][]float64, y []float64) error
	Score(X [][]float64, y []float64) (float64, error)
}

// CrossValScore evaluates a score by cross-validation.
func CrossValScore(model Estimator, X [][]float64, y []float64, cv KFold) ([]float64, error) {
	nSamples := len(X)
	folds := cv.Split(nSamples)
	scores := make([]float64, len(folds))

	for i, fold := range folds {
		// Prepare data for this fold
		XTrain := make([][]float64, len(fold.TrainIndices))
		yTrain := make([]float64, len(fold.TrainIndices))
		for j, idx := range fold.TrainIndices {
			XTrain[j] = X[idx]
			yTrain[j] = y[idx]
		}

		XTest := make([][]float64, len(fold.TestIndices))
		yTest := make([]float64, len(fold.TestIndices))
		for j, idx := range fold.TestIndices {
			XTest[j] = X[idx]
			yTest[j] = y[idx]
		}

		// Fit and Score
		err := model.Fit(XTrain, yTrain)
		if err != nil {
			return nil, err
		}

		score, err := model.Score(XTest, yTest)
		if err != nil {
			return nil, err
		}
		scores[i] = score
	}

	return scores, nil
}
