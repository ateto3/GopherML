// Package linear_model provides linear models for regression and classification.
//
// This package implements the following models:
//   - LinearRegression: Ordinary Least Squares
//   - Ridge: L2 regularized linear regression
//   - Lasso: L1 regularized linear regression
//   - ElasticNet: Combined L1 and L2 regularized linear regression
//
// All models follow the scikit-learn API pattern with Fit(), Predict(), and Score() methods.
//
// Example usage:
//
//	// Create and train a linear regression model
//	model := linear_model.NewLinearRegression()
//	err := model.Fit(X_train, y_train)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Make predictions
//	predictions, err := model.Predict(X_test)
//
//	// Evaluate the model
//	r2, err := model.Score(X_test, y_test)
//
// The package is optimized for performance using Go's native concurrency and
// minimal memory allocations.
package linear_model
