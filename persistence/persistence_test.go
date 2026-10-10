package persistence

import (
	"os"
	"testing"

	"github.com/ateto3/gopherml/linear_model"
	"github.com/ateto3/gopherml/tree"
)

func TestPersistenceLinearRegression(t *testing.T) {
	// 1. Create and "fit" a model (manually setting parameters for simplicity)
	lr := linear_model.NewLinearRegression()
	lr.Coef = []float64{1.5, -2.0}
	lr.Intercept = 0.5
	lr.NFeaturesIn = 2
	lr.IsFitted = true

	filename := "test_lr_model.gob"
	defer os.Remove(filename)

	// 2. Save the model
	if err := Save(filename, lr); err != nil {
		t.Fatalf("Failed to save model: %v", err)
	}

	// 3. Load the model into a new instance
	lrLoaded := linear_model.NewLinearRegression()
	if err := Load(filename, lrLoaded); err != nil {
		t.Fatalf("Failed to load model: %v", err)
	}

	// 4. Verify fields
	if lrLoaded.Intercept != lr.Intercept {
		t.Errorf("Expected intercept %v, got %v", lr.Intercept, lrLoaded.Intercept)
	}
	if len(lrLoaded.Coef) != len(lr.Coef) {
		t.Fatalf("Expected coef length %v, got %v", len(lr.Coef), len(lrLoaded.Coef))
	}
	for i := range lr.Coef {
		if lrLoaded.Coef[i] != lr.Coef[i] {
			t.Errorf("Expected coef[%d] %v, got %v", i, lr.Coef[i], lrLoaded.Coef[i])
		}
	}
}

func TestPersistenceDecisionTree(t *testing.T) {
	// 1. Create and "fit" a tree
	dt := tree.NewDecisionTreeClassifier()
	dt.NFeatures = 4
	dt.NClasses = 3
	dt.Classes = []int{0, 1, 2}
	dt.MaxFeatures = "sqrt" // Test interface{} field
	dt.IsFitted = true

	// Create a simple tree structure
	tr := tree.NewTree(4, 1, dt.Classes)
	tr.AddNode(-1, false, 0, 0.5, 100, 0.0, []float64{33, 33, 34})
	dt.Tree = tr

	filename := "test_dt_model.gob"
	defer os.Remove(filename)

	// 2. Save
	if err := Save(filename, dt); err != nil {
		t.Fatalf("Failed to save model: %v", err)
	}

	// 3. Load
	dtLoaded := tree.NewDecisionTreeClassifier()
	if err := Load(filename, dtLoaded); err != nil {
		t.Fatalf("Failed to load model: %v", err)
	}

	// 4. Verify
	if dtLoaded.NFeatures != dt.NFeatures {
		t.Errorf("Expected NFeatures %v, got %v", dt.NFeatures, dtLoaded.NFeatures)
	}
	if len(dtLoaded.Tree.Nodes) != len(dt.Tree.Nodes) {
		t.Errorf("Expected %d nodes, got %d", len(dt.Tree.Nodes), len(dtLoaded.Tree.Nodes))
	}
	if dtLoaded.MaxFeatures != dt.MaxFeatures {
		t.Errorf("Expected MaxFeatures %v, got %v", dt.MaxFeatures, dtLoaded.MaxFeatures)
	}
}
