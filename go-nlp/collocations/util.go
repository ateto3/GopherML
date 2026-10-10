package collocations

import (
	"math"
	"strings"
)

func split2(s string) []string  { return strings.SplitN(s, "|", 2) }
func split3(s string) []string { p := strings.SplitN(s, "|", 3); return p }
func join(v []string) string   { return strings.Join(v, "|") }
func log2(x float64) float64    { return math.Log2(x) }
