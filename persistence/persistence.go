package persistence

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
)

// Validator is implemented by models that can verify their own integrity after loading.
// Implement this interface on any model to get automatic post-load validation.
type Validator interface {
	ValidateLoaded() error
}

// init registers common types for gob serialization.
// Any concrete type stored in an interface{} field must be registered here.
func init() {
	// Generic container types
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})

	// Concrete types used in MaxFeatures (interface{}) fields across tree/ensemble packages.
	// Without these registrations, gob will fail to encode/decode models that use
	// integer or float MaxFeatures values.
	gob.Register(int(0))
	gob.Register(float64(0))
	gob.Register(string(""))
}

// Save serializes the model and saves it to a file.
// Equivalent to joblib.dump(model, filename) in Python.
func Save(filename string, model interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	return Dump(model, file)
}

// Dump serializes the model to an io.Writer.
func Dump(model interface{}, w io.Writer) error {
	encoder := gob.NewEncoder(w)
	if err := encoder.Encode(model); err != nil {
		return fmt.Errorf("failed to encode model: %w", err)
	}
	return nil
}

// Load loads a serialized model from a file into the provided pointer.
// Equivalent to model = joblib.load(filename) in Python.
//
// Usage:
//
//	var rf ensemble.RandomForestClassifier
//	err := persistence.Load("model.gob", &rf)
//
// If the model implements the Validator interface, ValidateLoaded() is called
// automatically to detect corruption or partially written files.
func Load(filename string, model interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return LoadFrom(file, model)
}

// LoadFrom loads a serialized model from an io.Reader.
// Calls ValidateLoaded() on the model if it implements Validator.
func LoadFrom(r io.Reader, model interface{}) error {
	decoder := gob.NewDecoder(r)
	if err := decoder.Decode(model); err != nil {
		return fmt.Errorf("failed to decode model: %w", err)
	}

	// Post-load integrity check
	if v, ok := model.(Validator); ok {
		if err := v.ValidateLoaded(); err != nil {
			return fmt.Errorf("model integrity check failed after loading: %w", err)
		}
	}

	return nil
}
