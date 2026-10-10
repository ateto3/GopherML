// Package tree provides tree building algorithms.
package tree

import (
"container/list"
"errors"
"fmt"
)

// DepthFirstTreeBuilder builds a tree using depth-first strategy.
type DepthFirstTreeBuilder struct {
Splitter              Splitter
Criterion             Criterion
MaxDepth              int
MinSamplesSplit       int
MinSamplesLeaf        int
MinWeightFractionLeaf float64
MinImpurityDecrease   float64
}

// NewDepthFirstTreeBuilder creates a new builder.
func NewDepthFirstTreeBuilder(splitter Splitter, criterion Criterion, maxDepth, minSamplesSplit, minSamplesLeaf int, minImpurityDecrease float64) *DepthFirstTreeBuilder {
return &DepthFirstTreeBuilder{
Splitter:            splitter,
Criterion:           criterion,
MaxDepth:            maxDepth,
MinSamplesSplit:     minSamplesSplit,
MinSamplesLeaf:      minSamplesLeaf,
MinImpurityDecrease: minImpurityDecrease,
}
}

// StackRecord represents a record on the stack for depth-first building.
type StackRecord struct {
Start    int
End      int
Depth    int
Parent   int
IsLeft   bool
Impurity float64
}

// Build constructs the decision tree.
func (b *DepthFirstTreeBuilder) Build(X [][]float64, y [][]float64, sampleWeight []float64, nClasses int) (*Tree, error) {
if len(X) == 0 || len(y) == 0 {
return nil, errors.New("empty training data")
}
if len(X) != len(y) {
return nil, errors.New("X and y have different lengths")
}

nSamples := len(X)
nFeatures := len(X[0])
nOutputs := len(y[0])

// Initialize indices
indices := make([]int, nSamples)
for i := 0; i < nSamples; i++ {
indices[i] = i
}

// Initialize splitter and criterion
b.Splitter.Init(X, y, sampleWeight, indices)
b.Criterion.Init(y, sampleWeight, indices, 0, nSamples)

// Determine tree type
var classes []int
if nClasses > 0 {
classes = make([]int, nClasses)
for i := 0; i < nClasses; i++ {
classes[i] = i
}
}

tree := NewTree(nFeatures, nOutputs, classes)

// Stack for depth-first building
stack := list.New()

// Push root node
rootRecord := StackRecord{
Start:    0,
End:      nSamples,
Depth:    0,
Parent:   -1,
IsLeft:   false,
Impurity: b.Criterion.NodeImpurity(),
}
stack.PushBack(rootRecord)

maxDepthSeen := 0

for stack.Len() > 0 {
// Pop from stack
elem := stack.Back()
record := elem.Value.(StackRecord)
stack.Remove(elem)

start := record.Start
end := record.End
depth := record.Depth
parent := record.Parent
isLeft := record.IsLeft
impurity := record.Impurity

nNodeSamples := end - start

// Check if we should split
isLeaf := b.shouldBeLeaf(depth, nNodeSamples, impurity)

if isLeaf {
// Create leaf node
b.Criterion.Reset(start, end)
nodeValue := b.Criterion.NodeValue()
tree.AddLeaf(parent, isLeft, nNodeSamples, impurity, nodeValue)
} else {
// Find best split
split, hasSplit := b.Splitter.NodeSplit(start, end, b.Criterion)

if !hasSplit {
// No valid split found, make leaf
b.Criterion.Reset(start, end)
nodeValue := b.Criterion.NodeValue()
tree.AddLeaf(parent, isLeft, nNodeSamples, impurity, nodeValue)
} else {
// Partition indices based on split
splitPos := b.partition(indices, X, start, end, split)

// Check if split actually partitioned the data
if splitPos == start || splitPos == end {
b.Criterion.Reset(start, end)
nodeValue := b.Criterion.NodeValue()
tree.AddLeaf(parent, isLeft, nNodeSamples, impurity, nodeValue)
continue
}

// Check impurity decrease
b.Criterion.Reset(start, end)
impurityDecrease := b.Criterion.NodeImpurity() - split.Impurity

if impurityDecrease < b.MinImpurityDecrease {
// Impurity decrease too small, make leaf
b.Criterion.Reset(start, end)
nodeValue := b.Criterion.NodeValue()
tree.AddLeaf(parent, isLeft, nNodeSamples, impurity, nodeValue)
} else {
// Create internal node
nodeIdx := tree.AddNode(parent, isLeft, split.Feature, split.Threshold,
nNodeSamples, impurity, []float64{})

if depth > maxDepthSeen {
maxDepthSeen = depth
}

// Push children (right first so left is processed first)
if b.MaxDepth == -1 || depth+1 < b.MaxDepth {
// Right child
rightRecord := StackRecord{
Start:    splitPos,
End:      end,
Depth:    depth + 1,
Parent:   nodeIdx,
IsLeft:   false,
Impurity: split.Impurity,
}
stack.PushBack(rightRecord)

// Left child
leftRecord := StackRecord{
Start:    start,
End:      splitPos,
Depth:    depth + 1,
Parent:   nodeIdx,
IsLeft:   true,
Impurity: split.Impurity,
}
stack.PushBack(leftRecord)
} else {
// At max depth, make both children leaves
b.Criterion.Reset(start, splitPos)
leftValue := b.Criterion.NodeValue()
tree.AddLeaf(nodeIdx, true, splitPos-start, split.Impurity, leftValue)

b.Criterion.Reset(splitPos, end)
rightValue := b.Criterion.NodeValue()
tree.AddLeaf(nodeIdx, false, end-splitPos, split.Impurity, rightValue)
}
}
}
}
}

tree.MaxDepth = maxDepthSeen
return tree, nil
}

// partition reorders indices[start:end] based on the split and returns the split position.
func (b *DepthFirstTreeBuilder) partition(indices []int, X [][]float64, start, end int, split SplitRecord) int {
i := start
j := end - 1
for i <= j {
if X[indices[i]][split.Feature] <= split.Threshold {
i++
} else {
indices[i], indices[j] = indices[j], indices[i]
j--
}
}
return i
}

// shouldBeLeaf determines if a node should be a leaf.
func (b *DepthFirstTreeBuilder) shouldBeLeaf(depth int, nSamples int, impurity float64) bool {
if b.MaxDepth != -1 && depth >= b.MaxDepth {
return true
}
if nSamples < b.MinSamplesSplit {
return true
}
if nSamples < 2*b.MinSamplesLeaf {
return true
}
if impurity <= 0 {
return true
}
return false
}

// TreeParams contains parameters for tree building.
type TreeParams struct {
MaxDepth            int
MinSamplesSplit     int
MinSamplesLeaf      int
MinImpurityDecrease float64
NClasses            int
}

// DefaultTreeParams returns default tree parameters.
func DefaultTreeParams() TreeParams {
return TreeParams{
MaxDepth:            -1,
MinSamplesSplit:     2,
MinSamplesLeaf:      1,
MinImpurityDecrease: 0.0,
NClasses:            0,
}
}

// BuildTree is a convenience function to build a tree.
func BuildTree(X [][]float64, y [][]float64, sampleWeight []float64, criterion Criterion, splitter Splitter, params TreeParams) (*Tree, error) {
builder := NewDepthFirstTreeBuilder(
splitter,
criterion,
params.MaxDepth,
params.MinSamplesSplit,
params.MinSamplesLeaf,
params.MinImpurityDecrease,
)

nClasses := 0
if params.NClasses > 0 {
nClasses = params.NClasses
}

return builder.Build(X, y, sampleWeight, nClasses)
}

// PrintTree prints a visual representation of the tree for debugging.
func PrintTree(t *Tree, featureNames []string) {
if len(t.Nodes) == 0 {
fmt.Println("Empty tree")
return
}
printNode(t, 0, "", featureNames, true)
}

func printNode(t *Tree, nodeIdx int, prefix string, featureNames []string, isLast bool) {
node := &t.Nodes[nodeIdx]

connector := "├── "
if isLast {
connector = "└── "
}

if node.IsLeaf() {
fmt.Printf("%s%sLeaf: value=%v, samples=%d\n", prefix, connector, node.Value, node.NSamples)
} else {
featureName := fmt.Sprintf("feature_%d", node.Feature)
if featureNames != nil && node.Feature < len(featureNames) {
featureName = featureNames[node.Feature]
}
fmt.Printf("%s%s%s <= %.3f? (impurity=%.4f, samples=%d)\n",
prefix, connector, featureName, node.Threshold, node.Impurity, node.NSamples)
}

if !node.IsLeaf() {
newPrefix := prefix + "│   "
if isLast {
newPrefix = prefix + "    "
}

if node.Left != -1 {
printNode(t, node.Left, newPrefix, featureNames, node.Right == -1)
}
if node.Right != -1 {
printNode(t, node.Right, newPrefix, featureNames, true)
}
}
}
