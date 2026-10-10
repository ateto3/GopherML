package preprocessing

import (
	"fmt"
	"sort"
)

// LabelEncoder encodes target labels with value between 0 and n_classes-1.
type LabelEncoder struct {
	Classes []string
	Mapping map[string]int
}

// NewLabelEncoder creates a new LabelEncoder.
func NewLabelEncoder() *LabelEncoder {
	return &LabelEncoder{
		Mapping: make(map[string]int),
	}
}

// Fit fits the label encoder.
func (s *LabelEncoder) Fit(y []string) {
	classMap := make(map[string]bool)
	for _, val := range y {
		classMap[val] = true
	}

	s.Classes = make([]string, 0, len(classMap))
	for class := range classMap {
		s.Classes = append(s.Classes, class)
	}
	sort.Strings(s.Classes)

	s.Mapping = make(map[string]int, len(s.Classes))
	for i, class := range s.Classes {
		s.Mapping[class] = i
	}
}

// Transform transforms labels to normalized encoding.
// Returns an error if any label was not seen during Fit.
func (s *LabelEncoder) Transform(y []string) ([]float64, error) {
	if len(s.Classes) == 0 {
		return nil, fmt.Errorf("LabelEncoder is not fitted — call Fit before Transform")
	}
	result := make([]float64, len(y))
	for i, val := range y {
		idx, ok := s.Mapping[val]
		if !ok {
			return nil, fmt.Errorf("LabelEncoder: unseen label %q at index %d", val, i)
		}
		result[i] = float64(idx)
	}
	return result, nil
}

// InverseTransform transforms labels back to original encoding.
// Returns an error if any encoded value is out of range.
func (s *LabelEncoder) InverseTransform(y []float64) ([]string, error) {
	if len(s.Classes) == 0 {
		return nil, fmt.Errorf("LabelEncoder is not fitted — call Fit before InverseTransform")
	}
	result := make([]string, len(y))
	for i, val := range y {
		idx := int(val)
		if idx < 0 || idx >= len(s.Classes) {
			return nil, fmt.Errorf("LabelEncoder: encoded value %d at index %d is out of range [0, %d)", idx, i, len(s.Classes))
		}
		result[i] = s.Classes[idx]
	}
	return result, nil
}
