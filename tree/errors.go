package tree

import "errors"

// Common errors used throughout the tree package.
var (
	// ErrNotFitted indicates that the tree has not been fitted yet.
	ErrNotFitted = errors.New("tree is not fitted yet")

	// ErrEmptyData indicates that empty data was provided.
	ErrEmptyData = errors.New("empty data provided")

	// ErrDimensionMismatch indicates that dimensions of inputs don't match.
	ErrDimensionMismatch = errors.New("dimension mismatch")

	// ErrInvalidParameter indicates an invalid parameter value.
	ErrInvalidParameter = errors.New("invalid parameter value")

	// ErrUnknownCriterion indicates an unknown criterion was specified.
	ErrUnknownCriterion = errors.New("unknown criterion")

	// ErrUnknownSplitter indicates an unknown splitter was specified.
	ErrUnknownSplitter = errors.New("unknown splitter")
)
