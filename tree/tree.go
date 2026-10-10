// Package tree provides decision tree-based models for classification and regression.
// Based on scikit-learn's tree module, optimized for Go with native data structures
// and goroutine-based parallelization where applicable.
package tree

// Node represents a node in the decision tree.
// Using a flat slice-based tree structure for cache efficiency.
type Node struct {
	// Feature index used for splitting (-1 for leaf nodes)
	Feature int
	// Threshold value for the split (NaN for leaf nodes)
	Threshold float64
	// Left child index (-1 for leaf)
	Left int
	// Right child index (-1 for leaf)
	Right int
	// Index of the parent node (-1 for root)
	Parent int
	// Number of samples at this node
	NSamples int
	// Impurity at this node
	Impurity float64
	// Value stored at this node (for regression: mean, for classification: class counts)
	Value []float64
}

// IsLeaf returns true if this node is a leaf node.
func (n *Node) IsLeaf() bool {
	return n.Left == -1 && n.Right == -1
}

// Tree stores the decision tree structure using a flat array of nodes.
// This is more memory-efficient than a pointer-based tree.
type Tree struct {
	// Nodes stored in pre-order (root, left subtree, right subtree)
	Nodes []Node
	// Number of features in the input data
	NFeatures int
	// Number of outputs (1 for single output, >1 for multi-output)
	NOutputs int
	// Maximum depth of the tree
	MaxDepth int
	// Number of classes for classification (nil for regression)
	Classes []int
}

// NewTree creates a new empty tree.
func NewTree(nFeatures, nOutputs int, classes []int) *Tree {
	return &Tree{
		Nodes:     make([]Node, 0, 1024),
		NFeatures: nFeatures,
		NOutputs:  nOutputs,
		MaxDepth:  0,
		Classes:   classes,
	}
}

// AddNode adds a new node to the tree and returns its index.
func (t *Tree) AddNode(parent int, isLeft bool, feature int, threshold float64,
	nSamples int, impurity float64, value []float64) int {

	nodeIdx := len(t.Nodes)
	node := Node{
		Feature:   feature,
		Threshold: threshold,
		Left:      -1,
		Right:     -1,
		Parent:    parent,
		NSamples:  nSamples,
		Impurity:  impurity,
		Value:     make([]float64, len(value)),
	}
	copy(node.Value, value)

	t.Nodes = append(t.Nodes, node)

	// Update parent's child reference
	if parent != -1 {
		if isLeft {
			t.Nodes[parent].Left = nodeIdx
		} else {
			t.Nodes[parent].Right = nodeIdx
		}
	}

	return nodeIdx
}

// AddLeaf adds a leaf node to the tree.
func (t *Tree) AddLeaf(parent int, isLeft bool, nSamples int, impurity float64, value []float64) int {
	return t.AddNode(parent, isLeft, -1, 0, nSamples, impurity, value)
}

// GetDepth returns the depth of a specific node.
func (t *Tree) GetDepth(nodeIdx int) int {
	depth := 0
	for t.Nodes[nodeIdx].Parent != -1 {
		nodeIdx = t.Nodes[nodeIdx].Parent
		depth++
	}
	return depth
}

// NLeaves returns the number of leaf nodes.
func (t *Tree) NLeaves() int {
	count := 0
	for i := range t.Nodes {
		if t.Nodes[i].IsLeaf() {
			count++
		}
	}
	return count
}

// Predict traverses the tree for a single sample and returns the node value.
func (t *Tree) Predict(x []float64) []float64 {
	if len(t.Nodes) == 0 {
		return nil
	}

	nodeIdx := 0
	for !t.Nodes[nodeIdx].IsLeaf() {
		node := &t.Nodes[nodeIdx]
		var nextIdx int
		if x[node.Feature] <= node.Threshold {
			nextIdx = node.Left
		} else {
			nextIdx = node.Right
		}
		// Guard against corrupt tree structure.
		if nextIdx < 0 || nextIdx >= len(t.Nodes) {
			break
		}
		nodeIdx = nextIdx
	}

	result := make([]float64, len(t.Nodes[nodeIdx].Value))
	copy(result, t.Nodes[nodeIdx].Value)
	return result
}

// PredictBatch predicts for multiple samples.
func (t *Tree) PredictBatch(X [][]float64) [][]float64 {
	results := make([][]float64, len(X))
	for i := range X {
		results[i] = t.Predict(X[i])
	}
	return results
}

// TreeBuilder interface for different tree building strategies.
type TreeBuilder interface {
	Build(X [][]float64, y [][]float64, sampleWeight []float64) (*Tree, error)
}
