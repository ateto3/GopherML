package probability

type ConditionalFreqDist map[string]FreqDist

func NewConditionalFreqDist() ConditionalFreqDist { return ConditionalFreqDist{} }

func (c ConditionalFreqDist) Inc(cond, sample string) {
	if c[cond] == nil {
		c[cond] = FreqDist{}
	}
	c[cond].Inc(sample)
}
