package tokenizer

// MacIntyreContractions holds regex patterns for English contractions,
// ported from nltk.tokenize.destructive.MacIntyreContractions.
type MacIntyreContractions struct {
	// CONTRACTIONS2 patterns: 2-part contractions
	CONTRACTIONS2 []string
	// CONTRACTIONS3 patterns: 3-part contractions
	CONTRACTIONS3 []string
}

// NewMacIntyreContractions returns a standard set of contraction patterns.
func NewMacIntyreContractions() *MacIntyreContractions {
	return &MacIntyreContractions{
		CONTRACTIONS2: []string{
			`(?i)\b(can)(not)\b`,
			`(?i)\b(d)('ye)\b`,
			`(?i)\b(gim)(me)\b`,
			`(?i)\b(gon)(na)\b`,
			`(?i)\b(got)(ta)\b`,
			`(?i)\b(lem)(me)\b`,
			`(?i)\b(mor)('n)\b`,
			`(?i)\b(wan)(na)\b`,
			`(?i)\b(o)('clock)\b`,
		},
		CONTRACTIONS3: []string{
			`(?i)('ll)(.*)`,
			`(?i)('re)(.*)`,
			`(?i)('ve)(.*)`,
			`(?i)(n't)(.*)`,
			`(?i)('LL)(.*)`,
			`(?i)('RE)(.*)`,
			`(?i)('VE)(.*)`,
			`(?i)(N'T)(.*)`,
			`(?i)('[sSmMdD])(.*)`,
			`(?i)('S')(.*)`,
			`(?i)('M)(.*)`,
			`(?i)('D)(.*)`,
		},
	}
}
