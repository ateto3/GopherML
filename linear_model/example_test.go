package linear_model_test

import (
	"fmt"
	"log"

	"github.com/ateto3/gopherml/linear_model"
)

func ExampleLinearRegression() {
	// Training data: y = 2*x + 1
	X := [][]float64{
		{1},
		{2},
		{3},
		{4},
		{5},
	}
	y := []float64{3, 5, 7, 9, 11}

	// Create and train model
	model := linear_model.NewLinearRegression()
	if err := model.Fit(X, y); err != nil {
		log.Fatal(err)
	}

	// Make predictions
	predictions, err := model.Predict([][]float64{{6}, {7}})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Coefficient: %.2f\n", model.Coef[0])
	fmt.Printf("Intercept: %.2f\n", model.Intercept)
	fmt.Printf("Predictions: %.1f, %.1f\n", predictions[0], predictions[1])

	// Output:
	// Coefficient: 2.00
	// Intercept: 1.00
	// Predictions: 13.0, 15.0
}

func ExampleRidge() {
	// Training data with some noise
	X := [][]float64{
		{1},
		{2},
		{3},
		{4},
		{5},
	}
	y := []float64{3.1, 4.9, 7.2, 8.8, 11.1}

	// Ridge regression with regularization
	model := linear_model.NewRidgeWithParams(linear_model.RidgeParams{
		Alpha:        0.5,
		FitIntercept: true,
	})

	if err := model.Fit(X, y); err != nil {
		log.Fatal(err)
	}

	score, _ := model.Score(X, y)
	fmt.Printf("R² Score: %.3f\n", score)
}

func ExampleLasso() {
	// Training data with irrelevant feature
	X := [][]float64{
		{1, 0.1},
		{2, 0.3},
		{3, 0.1},
		{4, 0.3},
		{5, 0.2},
	}
	y := []float64{3, 5, 7, 9, 11} // y depends mostly on first feature

	// Lasso regression - should zero out second feature
	model := linear_model.NewLassoWithParams(linear_model.LassoParams{
		Alpha:   0.1,
		MaxIter: 1000,
	})

	if err := model.Fit(X, y); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Coefficients: %.2f, %.2f\n", model.Coef[0], model.Coef[1])
	fmt.Printf("Iterations: %d\n", model.NIter)
}
