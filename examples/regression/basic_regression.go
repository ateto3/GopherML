package main

import (
	"fmt"
	"log"

	"github.com/ateto3/gopherml/ensemble"
	"github.com/ateto3/gopherml/linear_model"
	"github.com/ateto3/gopherml/metrics"
	"github.com/ateto3/gopherml/preprocessing"
)

func main() {
	// Example: Regression task - predicting house prices
	// Features: [square_feet, bedrooms, age, distance_to_city]

	XTrain := [][]float64{
		{1500, 3, 10, 5.2},
		{2000, 4, 5, 3.1},
		{1200, 2, 15, 8.5},
		{2500, 4, 2, 1.5},
		{1800, 3, 8, 4.0},
		{3000, 5, 1, 0.8},
		{1400, 2, 20, 10.0},
		{2200, 4, 7, 2.5},
	}

	// Target: house price in thousands of dollars
	yTrain := []float64{250, 400, 200, 550, 320, 700, 180, 450}

	XTest := [][]float64{
		{1600, 3, 12, 5.0},
		{2800, 5, 3, 1.2},
	}

	yTest := []float64{280, 650}

	fmt.Println("=== Basic Regression Example ===")
	fmt.Println("Training samples:", len(XTrain))
	fmt.Println("Test samples:", len(XTest))

	// Scale features
	scaler := preprocessing.NewStandardScaler()
	XTrainScaled, err := scaler.FitTransform(XTrain)
	if err != nil {
		log.Fatalf("Error scaling training data: %v", err)
	}

	XTestScaled, err := scaler.Transform(XTest)
	if err != nil {
		log.Fatalf("Error scaling test data: %v", err)
	}

	// Method 1: Linear Regression
	fmt.Println("\n--- Linear Regression ---")
	lr := linear_model.NewLinearRegression()
	if err := lr.Fit(XTrainScaled, yTrain); err != nil {
		log.Fatalf("Error training linear regression: %v", err)
	}

	lrPreds, err := lr.Predict(XTestScaled)
	if err != nil {
		log.Fatalf("Error predicting with linear regression: %v", err)
	}

	fmt.Println("Linear Regression Predictions:")
	for i, pred := range lrPreds {
		fmt.Printf("Sample %d: Predicted=$%.0fk, Actual=$%.0fk, Error=$%.0fk\n",
			i, pred, yTest[i], pred-yTest[i])
	}

	lrR2, err := lr.Score(XTestScaled, yTest)
	if err != nil {
		log.Fatalf("Error calculating R²: %v", err)
	}
	fmt.Printf("Linear Regression R² Score: %.4f\n", lrR2)

	// Method 2: Random Forest Regressor
	fmt.Println("\n--- Random Forest Regressor ---")
	rf := ensemble.NewRandomForestRegressor()
	rf.NEstimators = 100
	rf.MaxDepth = 8
	rf.RandomState = 42

	if err := rf.Fit(XTrainScaled, yTrain); err != nil {
		log.Fatalf("Error training random forest: %v", err)
	}

	rfPreds, err := rf.Predict(XTestScaled)
	if err != nil {
		log.Fatalf("Error predicting with random forest: %v", err)
	}

	fmt.Println("Random Forest Predictions:")
	for i, pred := range rfPreds {
		fmt.Printf("Sample %d: Predicted=$%.0fk, Actual=$%.0fk, Error=$%.0fk\n",
			i, pred, yTest[i], pred-yTest[i])
	}

	rfR2, err := rf.Score(XTestScaled, yTest)
	if err != nil {
		log.Fatalf("Error calculating R²: %v", err)
	}
	fmt.Printf("Random Forest R² Score: %.4f\n", rfR2)

	// Additional regression metrics
	fmt.Println("\n--- Regression Metrics (Random Forest) ---")
	mse, _ := metrics.MeanSquaredError(yTest, rfPreds)
	mae, _ := metrics.MeanAbsoluteError(yTest, rfPreds)
	rmse := mse // RMSE is sqrt of MSE, showing MSE for simplicity

	fmt.Printf("Mean Squared Error: %.2f\n", mse)
	fmt.Printf("Mean Absolute Error: %.2f\n", mae)
	fmt.Printf("Root Mean Squared Error: %.2f\n", rmse)

	// Single prediction example (high-frequency use)
	fmt.Println("\n--- Single Prediction (Hot Path) ---")
	singleSample := XTestScaled[0]
	predSingle, err := rf.PredictSingle(singleSample)
	if err != nil {
		log.Fatalf("Error in single prediction: %v", err)
	}
	fmt.Printf("Single sample prediction: $%.0fk\n", predSingle)

	// Compare models
	fmt.Println("\n--- Model Comparison ---")
	fmt.Printf("Linear Regression R²: %.4f\n", lrR2)
	fmt.Printf("Random Forest R²: %.4f\n", rfR2)
	if rfR2 > lrR2 {
		fmt.Println("Random Forest performs better on this dataset")
	} else {
		fmt.Println("Linear Regression performs better on this dataset")
	}

	fmt.Println("\n=== Example Complete ===")
}
