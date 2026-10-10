package ensemble

import (
	"fmt"
	"sync"

	"github.com/ateto3/gopherml/metrics"
	"github.com/ateto3/gopherml/tree"
	"github.com/ateto3/gopherml/utils"
)

// RandomForestClassifier implements a random forest classifier.
type RandomForestClassifier struct {
	BaseForest
	Estimators []*tree.DecisionTreeClassifier
	// Classes holds the sorted unique class labels seen during Fit.
	Classes  []int
	NClasses int
}

// NewRandomForestClassifier creates a new RandomForestClassifier with sane defaults.
func NewRandomForestClassifier() *RandomForestClassifier {
	return &RandomForestClassifier{
		BaseForest: BaseForest{
			NEstimators:     100,
			MaxDepth:        10,
			MinSamplesSplit: 2,
			MaxFeatures:     "sqrt",
			Bootstrap:       true,
			RandomState:     -1,
		},
	}
}

// Fit trains the forest on (X, y). Validates for NaN/Inf before training.
func (rf *RandomForestClassifier) Fit(X [][]float64, y []float64) error {
	if err := utils.CheckXy(X, y); err != nil {
		return err
	}

	nSamples := len(X)
	rf.NFeatures = len(X[0])
	rf.Estimators = make([]*tree.DecisionTreeClassifier, rf.NEstimators)

	// Collect errors from goroutines via a buffered channel.
	errCh := make(chan error, rf.NEstimators)

	// Mutex to safely capture class metadata from the first completed tree.
	var mu sync.Mutex
	classesCaptured := false

	var wg sync.WaitGroup
	for i := 0; i < rf.NEstimators; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			clf := tree.NewDecisionTreeClassifier()
			clf.MaxDepth = rf.MaxDepth
			clf.MinSamplesSplit = rf.MinSamplesSplit
			clf.MaxFeatures = rf.MaxFeatures

			rng := utils.NewRNG(rf.RandomState + int64(idx))
			indices := make([]int, nSamples)
			if !rf.Bootstrap {
				for j := range indices {
					indices[j] = j
				}
			} else {
				for j := 0; j < nSamples; j++ {
					indices[j] = rng.Intn(nSamples)
				}
			}

			XSub := make([][]float64, nSamples)
			ySub := make([]float64, nSamples)
			for j, rowIdx := range indices {
				XSub[j] = X[rowIdx]
				ySub[j] = y[rowIdx]
			}

			if err := clf.Fit(XSub, ySub, nil); err != nil {
				errCh <- fmt.Errorf("estimator[%d]: %w", idx, err)
				return
			}
			rf.Estimators[idx] = clf

			// Capture class metadata from first tree that finishes.
			mu.Lock()
			if !classesCaptured && clf.IsFitted {
				rf.Classes = clf.Classes
				rf.NClasses = clf.NClasses
				classesCaptured = true
			}
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	close(errCh)

	// Return the first error encountered, if any.
	if err := <-errCh; err != nil {
		return err
	}

	rf.IsFitted = true
	return nil
}

// Predict returns the majority-vote class label for each sample in X.
func (rf *RandomForestClassifier) Predict(X [][]float64) ([]int, error) {
	if err := utils.CheckIsFitted(rf.IsFitted, "RandomForestClassifier"); err != nil {
		return nil, err
	}

	nSamples := len(X)
	finalPreds := make([]int, nSamples)

	allPreds := make([][]int, rf.NEstimators)
	for i := 0; i < rf.NEstimators; i++ {
		allPreds[i], _ = rf.Estimators[i].Predict(X)
	}

	for i := 0; i < nSamples; i++ {
		counts := make(map[int]int)
		for j := 0; j < rf.NEstimators; j++ {
			counts[allPreds[j][i]]++
		}

		maxCount := -1
		bestLabel := -1
		for label, count := range counts {
			if count > maxCount {
				maxCount = count
				bestLabel = label
			}
		}
		finalPreds[i] = bestLabel
	}

	return finalPreds, nil
}

// PredictProba returns per-class probability estimates for each sample in X.
// Each row sums to 1.0. Essential for confidence-based trade entry thresholds.
//
// Example crypto bot usage:
//
//	proba, _ := rf.PredictProba([][]float64{features})
//	if proba[0][1] > 0.70 { // 70% confidence to go long
//	    executeTrade()
//	}
func (rf *RandomForestClassifier) PredictProba(X [][]float64) ([][]float64, error) {
	if err := utils.CheckIsFitted(rf.IsFitted, "RandomForestClassifier"); err != nil {
		return nil, err
	}

	nSamples := len(X)
	nClasses := rf.NClasses
	if nClasses == 0 && len(rf.Estimators) > 0 {
		nClasses = rf.Estimators[0].NClasses
	}
	if nClasses == 0 {
		return nil, fmt.Errorf("RandomForestClassifier: NClasses is 0, model may be corrupted")
	}

	// Accumulate vote fractions from each tree
	proba := make([][]float64, nSamples)
	for i := range proba {
		proba[i] = make([]float64, nClasses)
	}

	for _, est := range rf.Estimators {
		treeProba, err := est.PredictProba(X)
		if err != nil {
			return nil, err
		}
		for i := 0; i < nSamples; i++ {
			for j := 0; j < len(treeProba[i]) && j < nClasses; j++ {
				proba[i][j] += treeProba[i][j]
			}
		}
	}

	// Normalize each row by its own sum to guarantee sum=1.0
	// (some bootstrap trees may not have seen all classes, so dividing by
	// NEstimators would give sum < 1 for those samples)
	for i := range proba {
		rowSum := 0.0
		for _, v := range proba[i] {
			rowSum += v
		}
		if rowSum > 0 {
			for j := range proba[i] {
				proba[i][j] /= rowSum
			}
		}
	}

	return proba, nil
}

// PredictSingle predicts the class for a single sample without batch allocation overhead.
// Use this in high-frequency bot loops (e.g., on every new mempool event or price tick).
func (rf *RandomForestClassifier) PredictSingle(x []float64) (int, error) {
	if err := utils.CheckIsFitted(rf.IsFitted, "RandomForestClassifier"); err != nil {
		return -1, err
	}
	if err := utils.CheckFiniteSingle(x); err != nil {
		return -1, err
	}

	counts := make(map[int]int, 4)
	xWrapped := [][]float64{x}
	for _, est := range rf.Estimators {
		preds, _ := est.Predict(xWrapped)
		if len(preds) > 0 {
			counts[preds[0]]++
		}
	}

	maxCount := -1
	bestLabel := -1
	for label, count := range counts {
		if count > maxCount {
			maxCount = count
			bestLabel = label
		}
	}
	return bestLabel, nil
}

// PredictProbaSingle returns class probabilities for a single feature vector.
// Returns a []float64 of length NClasses, summing to 1.0.
// Use for confidence-based signal filtering without batch slice overhead.
func (rf *RandomForestClassifier) PredictProbaSingle(x []float64) ([]float64, error) {
	if err := utils.CheckIsFitted(rf.IsFitted, "RandomForestClassifier"); err != nil {
		return nil, err
	}
	if err := utils.CheckFiniteSingle(x); err != nil {
		return nil, err
	}

	nClasses := rf.NClasses
	if nClasses == 0 && len(rf.Estimators) > 0 {
		nClasses = rf.Estimators[0].NClasses
	}

	proba := make([]float64, nClasses)
	xWrapped := [][]float64{x}

	for _, est := range rf.Estimators {
		treeProba, err := est.PredictProba(xWrapped)
		if err != nil {
			return nil, err
		}
		if len(treeProba) > 0 {
			for j := 0; j < len(treeProba[0]) && j < nClasses; j++ {
				proba[j] += treeProba[0][j]
			}
		}
	}

	rowSum := 0.0
	for _, v := range proba {
		rowSum += v
	}
	if rowSum > 0 {
		for j := range proba {
			proba[j] /= rowSum
		}
	}
	return proba, nil
}

// Score returns the accuracy on the given test data and labels.
func (rf *RandomForestClassifier) Score(X [][]float64, y []float64) (float64, error) {
	preds, err := rf.Predict(X)
	if err != nil {
		return 0, err
	}
	yPred := make([]float64, len(preds))
	for i, v := range preds {
		yPred[i] = float64(v)
	}
	return metrics.AccuracyScore(y, yPred)
}

// ValidateLoaded checks that a model loaded from disk is fully valid.
// Called automatically by persistence.Load when the model implements Validator.
func (rf *RandomForestClassifier) ValidateLoaded() error {
	if !rf.IsFitted {
		return fmt.Errorf("RandomForestClassifier loaded from disk is not fitted")
	}
	if len(rf.Estimators) == 0 {
		return fmt.Errorf("RandomForestClassifier loaded from disk has 0 estimators")
	}
	for i, est := range rf.Estimators {
		if est == nil || est.Tree == nil {
			return fmt.Errorf("RandomForestClassifier estimator[%d] is nil or has nil tree", i)
		}
	}
	return nil
}

// RandomForestRegressor implements a random forest regressor.
type RandomForestRegressor struct {
	BaseForest
	Estimators []*tree.DecisionTreeRegressor
}

// NewRandomForestRegressor creates a new RandomForestRegressor with sane defaults.
func NewRandomForestRegressor() *RandomForestRegressor {
	return &RandomForestRegressor{
		BaseForest: BaseForest{
			NEstimators:     100,
			MaxDepth:        10,
			MinSamplesSplit: 2,
			MaxFeatures:     "sqrt",
			Bootstrap:       true,
			RandomState:     -1,
		},
	}
}

// Fit trains the regressor forest on (X, y). Validates for NaN/Inf before training.
func (rf *RandomForestRegressor) Fit(X [][]float64, y []float64) error {
	if err := utils.CheckXy(X, y); err != nil {
		return err
	}

	nSamples := len(X)
	rf.NFeatures = len(X[0])
	rf.Estimators = make([]*tree.DecisionTreeRegressor, rf.NEstimators)

	// Collect errors from goroutines via a buffered channel.
	errCh := make(chan error, rf.NEstimators)

	var wg sync.WaitGroup
	for i := 0; i < rf.NEstimators; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			reg := tree.NewDecisionTreeRegressor()
			reg.MaxDepth = rf.MaxDepth
			reg.MinSamplesSplit = rf.MinSamplesSplit
			reg.MaxFeatures = rf.MaxFeatures

			rng := utils.NewRNG(rf.RandomState + int64(idx))
			indices := make([]int, nSamples)
			if !rf.Bootstrap {
				for j := range indices {
					indices[j] = j
				}
			} else {
				for j := 0; j < nSamples; j++ {
					indices[j] = rng.Intn(nSamples)
				}
			}

			XSub := make([][]float64, nSamples)
			ySub := make([]float64, nSamples)
			for j, rowIdx := range indices {
				XSub[j] = X[rowIdx]
				ySub[j] = y[rowIdx]
			}

			if err := reg.Fit(XSub, ySub, nil); err != nil {
				errCh <- fmt.Errorf("estimator[%d]: %w", idx, err)
				return
			}
			rf.Estimators[idx] = reg
		}(i)
	}
	wg.Wait()
	close(errCh)

	// Return the first error encountered, if any.
	if err := <-errCh; err != nil {
		return err
	}

	rf.IsFitted = true
	return nil
}

// Predict returns the averaged prediction for each sample in X.
func (rf *RandomForestRegressor) Predict(X [][]float64) ([]float64, error) {
	if err := utils.CheckIsFitted(rf.IsFitted, "RandomForestRegressor"); err != nil {
		return nil, err
	}

	nSamples := len(X)
	finalPreds := make([]float64, nSamples)

	for i := 0; i < rf.NEstimators; i++ {
		preds, _ := rf.Estimators[i].Predict(X)
		for j := 0; j < nSamples; j++ {
			finalPreds[j] += preds[j]
		}
	}

	for j := 0; j < nSamples; j++ {
		finalPreds[j] /= float64(rf.NEstimators)
	}

	return finalPreds, nil
}

// PredictSingle predicts for a single sample without batch allocation overhead.
// Ideal for high-frequency bot loops (per-tick price prediction).
func (rf *RandomForestRegressor) PredictSingle(x []float64) (float64, error) {
	if err := utils.CheckIsFitted(rf.IsFitted, "RandomForestRegressor"); err != nil {
		return 0, err
	}
	if err := utils.CheckFiniteSingle(x); err != nil {
		return 0, err
	}

	xWrapped := [][]float64{x}
	sum := 0.0
	for _, est := range rf.Estimators {
		preds, _ := est.Predict(xWrapped)
		if len(preds) > 0 {
			sum += preds[0]
		}
	}
	return sum / float64(rf.NEstimators), nil
}

// Score returns the R² score on the given test data and labels.
func (rf *RandomForestRegressor) Score(X [][]float64, y []float64) (float64, error) {
	preds, err := rf.Predict(X)
	if err != nil {
		return 0, err
	}
	return metrics.R2Score(y, preds)
}

// ValidateLoaded checks that a model loaded from disk is fully valid.
func (rf *RandomForestRegressor) ValidateLoaded() error {
	if !rf.IsFitted {
		return fmt.Errorf("RandomForestRegressor loaded from disk is not fitted")
	}
	if len(rf.Estimators) == 0 {
		return fmt.Errorf("RandomForestRegressor loaded from disk has 0 estimators")
	}
	return nil
}
