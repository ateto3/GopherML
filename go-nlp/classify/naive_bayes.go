package classify

import (
	"fmt"
	"math"
	"sort"
)

// NaiveBayesClassifier is a lightweight multinomial NB classifier.
type NaiveBayesClassifier struct {
	labels   []string
	logPrior map[string]float64
	logCond  map[string]map[string]map[string]float64
	defaults map[string]float64
}

// TrainNaiveBayes trains a Naive Bayes classifier from labeled feature sets.
func TrainNaiveBayes(data []LabeledFeatures) *NaiveBayesClassifier {
	labelCounts := map[string]int{}
	featCounts := map[string]map[string]map[string]int{}
	labelVals := map[string]map[string]struct{}{}
	for _, ex := range data {
		labelCounts[ex.Label]++
		if featCounts[ex.Label] == nil {
			featCounts[ex.Label] = map[string]map[string]int{}
			labelVals[ex.Label] = map[string]struct{}{}
		}
		for k, v := range ex.Features {
			sv := stringify(v)
			if featCounts[ex.Label][k] == nil {
				featCounts[ex.Label][k] = map[string]int{}
			}
			featCounts[ex.Label][k][sv]++
			labelVals[ex.Label][sv] = struct{}{}
		}
	}
	labels := make([]string, 0, len(labelCounts))
	for l := range labelCounts {
		labels = append(labels, l)
	}
	sort.Strings(labels)
	total := float64(len(data))
	c := &NaiveBayesClassifier{
		labels:   labels,
		logPrior: map[string]float64{},
		logCond:  map[string]map[string]map[string]float64{},
		defaults: map[string]float64{},
	}
	for _, l := range labels {
		c.logPrior[l] = math.Log(float64(labelCounts[l]) / total)
		if c.logCond[l] == nil {
			c.logCond[l] = map[string]map[string]float64{}
		}
		for feat, vals := range featCounts[l] {
			bins := float64(len(vals) + 1)
			if c.logCond[l][feat] == nil {
				c.logCond[l][feat] = map[string]float64{}
			}
			sum := 0
			for _, cnt := range vals {
				sum += cnt
			}
			c.defaults[feat] = math.Log(1.0 / (float64(sum) + bins))
			for val, cnt := range vals {
				c.logCond[l][feat][val] = math.Log((float64(cnt) + 1.0) / (float64(sum) + bins))
			}
		}
	}
	return c
}

func (c *NaiveBayesClassifier) Labels() []string { return c.labels }

func (c *NaiveBayesClassifier) Classify(features Features) string {
	best := ""
	bestScore := math.Inf(-1)
	for _, l := range c.labels {
		s := c.logPrior[l]
		for k, v := range features {
			sv := stringify(v)
			if cond := c.logCond[l][k]; cond != nil {
				if p, ok := cond[sv]; ok {
					s += p
				} else {
					s += c.defaults[k]
				}
			}
		}
		if s > bestScore {
			bestScore, best = s, l
		}
	}
	return best
}

func (c *NaiveBayesClassifier) ClassifyMany(featuresets []Features) []string {
	out := make([]string, len(featuresets))
	for i := range featuresets {
		out[i] = c.Classify(featuresets[i])
	}
	return out
}

func stringify(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case bool:
		if x {
			return "true"
		}
		return "false"
	default:
		return fmtAny(x)
	}
}

func fmtAny(v any) string {
	return fmt.Sprintf("%v", v)
}
