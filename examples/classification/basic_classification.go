package main

import (
	"fmt"
	"log"

	"github.com/ateto3/gopherml/ensemble"
	"github.com/ateto3/gopherml/metrics"
	"github.com/ateto3/gopherml/model_selection"
	"github.com/ateto3/gopherml/preprocessing"
)

func main() {
	// Example: Binary classification using Random Forest
	// This simulates a simple classification task (e.g., predicting if a transaction is fraudulent)

	// Sample data: [feature1, feature2, feature3]
	XTrain := [][]float64{
		{1.2, 2.3, 0.5},
		{1.5, 2.1, 0.6},
		{2.1, 1.8, 0.4},
		{2.3, 1.5, 0.3},
		{3.1, 0.8, 0.2},
		{3.5, 0.5, 0.1},
		{0.8, 3.1, 0.7},
		{0.5, 3.5, 0.8},
	}

	// Labels: 0 = not fraudulent, 1 = fraudulent
	yTrain := []float64{0, 0, 0, 0, 1, 1, 0, 0}

	XTest := [][]float64{
		{1.8, 1.9, 0.45},
		{3.2, 0.6, 0.15},
		{0.9, 3.0, 0.75},
	}

	yTest := []float64{0, 1, 0}

	fmt.Println("=== Basic Classification Example ===")
	fmt.Println("Training samples:", len(XTrain))
	fmt.Println("Test samples:", len(XTest))

	// Scale features for better performance
	scaler := preprocessing.NewStandardScaler()
	XTrainScaled, err := scaler.FitTransform(XTrain)
	if err != nil {
		log.Fatalf("Error scaling training data: %v", err)
	}

	XTestScaled, err := scaler.Transform(XTest)
	if err != nil {
		log.Fatalf("Error scaling test data: %v", err)
	}

	// Create and train Random Forest Classifier
	rf := ensemble.NewRandomForestClassifier()
	rf.NEstimators = 50
	rf.MaxDepth = 5
	rf.RandomState = 42

	fmt.Println("\nTraining Random Forest Classifier...")
	if err := rf.Fit(XTrainScaled, yTrain); err != nil {
		log.Fatalf("Error training model: %v", err)
	}

	// Make predictions
	predictions, err := rf.Predict(XTestScaled)
	if err != nil {
		log.Fatalf("Error making predictions: %v", err)
	}

	fmt.Println("\nPredictions:")
	for i, pred := range predictions {
		fmt.Printf("Sample %d: Predicted=%d, Actual=%.0f\n", i, pred, yTest[i])
	}

	// Get prediction probabilities (confidence scores)
	proba, err := rf.PredictProba(XTestScaled)
	if err != nil {
		log.Fatalf("Error getting probabilities: %v", err)
	}

	fmt.Println("\nPrediction Probabilities:")
	for i, p := range proba {
		fmt.Printf("Sample %d: Class 0=%.2f%%, Class 1=%.2f%%\n", i, p[0]*100, p[1]*100)
	}

	// Evaluate model
	accuracy, err := rf.Score(XTestScaled, yTest)
	if err != nil {
		log.Fatalf("Error calculating score: %v", err)
	}
	fmt.Printf("\nAccuracy: %.2f%%\n", accuracy*100)

	// Cross-validation
	fmt.Println("\n=== Cross-Validation ===")
	kf := model_selection.KFold{NSplits: 3, Shuffle: true, RandomState: 42}
	scores, err := model_selection.CrossValScore(rf, XTrainScaled, yTrain, kf)
	if err != nil {
		log.Fatalf("Error in cross-validation: %v", err)
	}

	fmt.Printf("Cross-validation scores: %v\n", scores)

	// Calculate mean CV score
	sum := 0.0
	for _, s := range scores {
		sum += s
	}
	fmt.Printf("Mean CV score: %.2f%%\n", (sum/float64(len(scores)))*100)

	// Additional metrics
	yPred := make([]float64, len(predictions))
	for i, p := range predictions {
		yPred[i] = float64(p)
	}

	precision, _ := metrics.PrecisionScore(yTest, yPred, "binary")
	recall, _ := metrics.RecallScore(yTest, yPred, "binary")
	f1, _ := metrics.F1Score(yTest, yPred, "binary")

	fmt.Println("\n=== Detailed Metrics ===")
	fmt.Printf("Precision: %.2f%%\n", precision*100)
	fmt.Printf("Recall: %.2f%%\n", recall*100)
	fmt.Printf("F1-Score: %.2f%%\n", f1*100)

	// Single prediction example (for high-frequency use)
	fmt.Println("\n=== Single Prediction (Hot Path) ===")
	singleSample := XTestScaled[0]
	label, err := rf.PredictSingle(singleSample)
	if err != nil {
		log.Fatalf("Error in single prediction: %v", err)
	}
	probaSingle, err := rf.PredictProbaSingle(singleSample)
	if err != nil {
		log.Fatalf("Error in single probability: %v", err)
	}

	fmt.Printf("Single sample prediction: %d\n", label)
	fmt.Printf("Confidence: Class 0=%.2f%%, Class 1=%.2f%%\n", probaSingle[0]*100, probaSingle[1]*100)

	fmt.Println("\n=== Example Complete ===")
}
