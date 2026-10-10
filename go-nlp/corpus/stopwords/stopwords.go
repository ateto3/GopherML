// Package stopwords provides optimized stopword dictionaries for multiple languages.
package stopwords

// Stopwords is a type for checking stopwords in a specific language.
type Stopwords struct {
	lang  string
	words map[string]bool
}

// defaultDicts stores stopwords for each language as a set.
var defaultDicts = map[string]map[string]bool{
	"en": {
		"a": true, "an": true, "the": true, "and": true, "or": true,
		"but": true, "in": true, "on": true, "at": true, "to": true,
		"for": true, "of": true, "with": true, "by": true, "from": true,
		"is": true, "are": true, "was": true, "were": true, "be": true,
		"have": true, "has": true, "had": true, "do": true, "does": true,
		"did": true, "will": true, "would": true, "could": true, "should": true,
		"may": true, "might": true, "must": true, "shall": true,
		"this": true, "that": true, "these": true, "those": true,
		"i": true, "you": true, "he": true, "she": true, "it": true,
		"we": true, "they": true, "what": true, "which": true,
		"who": true, "whom": true, "its": true, "their": true,
		"his": true, "her": true, "our": true, "your": true, "my": true,
		"if": true, "then": true, "than": true, "because": true,
		"as": true, "until": true, "while": true, "about": true,
		"against": true, "between": true, "into": true, "through": true,
		"during": true, "before": true, "after": true, "above": true,
		"below": true, "over": true, "under": true, "again": true,
		"further": true, "once": true, "here": true, "there": true,
		"when": true, "where": true, "why": true, "how": true,
		"all": true, "each": true, "few": true, "more": true, "most": true,
		"other": true, "some": true, "such": true, "no": true, "nor": true,
		"not": true, "only": true, "own": true, "same": true, "so": true,
		"two": true, "three": true, "one": true,
	},
	"es": {
		"el": true, "la": true, "los": true, "las": true, "de": true,
		"que": true, "y": true, "a": true, "en": true, "un": true,
		"ser": true, "se": true, "no": true, "le": true, "por": true,
		"para": true, "con": true, "su": true, "al": true, "del": true,
		"como": true, "pero": true, "porque": true, "cuando": true,
		"donde": true, "muy": true, "sin": true, "si": true, "ya": true,
		"esta": true, "este": true, "esto": true, "me": true, "mi": true,
		"mis": true, "tu": true, "te": true, "ti": true, "tus": true,
		"nunca": true, "tambien": true, "también": true, "solo": true,
		"tengo": true, "tenemos": true, "tienes": true, "tienen": true,
		"quiero": true, "quiere": true, "queremos": true, "quieren": true,
		"puedo": true, "puede": true, "podemos": true, "pueden": true,
	},
	"pt": {
		"o": true, "a": true, "os": true, "as": true, "de": true,
		"que": true, "e": true, "do": true, "da": true,
		"em": true, "um": true, "para": true, "com": true, "são": true,
		"foi": true, "se": true, "na": true, "no": true, "das": true,
		"dos": true, "por": true, "uma": true,
		"umas": true, "uns": true, "entre": true, "mais": true, "menos": true,
		"bem": true, "muito": true, "sempre": true, "nunca": true,
		"tambem": true, "também": true, "fui": true, "fomos": true,
		"sou": true, "somos": true, "eres": true,
	},
	"fr": {
		"le": true, "la": true, "les": true, "de": true, "et": true,
		"à": true, "en": true, "un": true, "une": true, "du": true,
		"des": true, "est": true, "sont": true, "ce": true, "cet": true,
		"cette": true, "ici": true, "il": true, "elle": true, "on": true,
		"nous": true, "vous": true, "ils": true, "elles": true,
		"que": true, "qui": true, "quoi": true, "dont": true, "où": true,
		"quand": true, "comment": true, "pourquoi": true, "combien": true,
		"non": true, "pas": true, "plus": true, "moins": true, "tout": true,
		"toute": true, "tous": true, "toutes": true, "même": true, "mes": true,
		"mon": true, "ma": true, "tes": true, "ton": true, "ta": true,
	},
	"de": {
		"der": true, "die": true, "das": true, "und": true, "in": true,
		"den": true, "von": true, "zu": true, "des": true,
		"dem": true, "ein": true, "eine": true, "einem": true, "einen": true,
		"einer": true, "eines": true, "ist": true, "sind": true, "war": true,
		"waren": true, "sein": true, "hat": true, "haben": true, "habe": true,
		"hatte": true, "hatten": true, "mit": true, "auf": true,
		"für": true, "über": true, "nicht": true, "kein": true, "keine": true,
	},
	"it": {
		"il": true, "la": true, "lo": true, "i": true, "gli": true,
		"le": true, "di": true, "e": true, "a": true, "da": true,
		"in": true, "un": true, "una": true, "uno": true, "per": true,
		"con": true, "su": true, "tra": true, "fra": true,
		"è": true, "sono": true, "fui": true, "fosti": true,
		"siamo": true, "siete": true, "siano": true, "sarà": true,
		"ho": true, "hai": true, "ha": true, "abbiamo": true, "avete": true,
		"hanno": true, "abbia": true, "abbiate": true, "abbiano": true,
		"non": true, "no": true, "più": true, "ancora": true,
	},
	"nl": {
		"de": true, "het": true, "een": true, "en": true, "van": true,
		"ik": true, "te": true, "dat": true, "in": true,
		"aan": true, "op": true, "met": true, "zijn": true,
		"is": true, "was": true, "waren": true, "heeft": true, "hebben": true,
		"niet": true, "voor": true, "bij": true, "ook": true, "om": true,
		"wie": true, "wat": true, "welke": true, "zij": true,
	},
	"ru": {
		"в": true, "и": true, "а": true, "то": true, "на": true,
		"с": true, "по": true, "к": true, "у": true, "о": true,
		"за": true, "из": true, "но": true, "что": true, "бы": true,
		"же": true, "как": true, "это": true, "его": true, "её": true,
		"их": true, "я": true, "ты": true, "он": true, "она": true,
		"мы": true, "вы": true, "они": true, "кто": true,
		"какой": true, "какая": true, "какое": true, "какие": true,
		"не": true, "ни": true, "для": true, "до": true, "после": true,
	},
	"zh": {
		"的": true, "一": true, "是": true, "在": true, "不": true,
		"了": true, "有": true, "和": true, "人": true, "这": true,
		"中": true, "大": true, "为": true, "上": true, "个": true,
		"国": true, "我": true, "以": true, "要": true, "他": true,
		"时": true, "来": true, "用": true, "们": true, "生": true,
		"到": true, "作": true, "天": true, "能": true, "对": true,
		"小": true, "多": true, "然": true, "于": true, "心": true,
		"学": true, "说": true, "过": true, "会": true, "可": true,
	},
	"ja": {
		"的": true, "一": true, "是": true, "在": true, "不": true,
		"了": true, "有": true, "和": true, "人": true, "这": true,
		"中": true, "大": true, "为": true, "上": true, "个": true,
		"国": true, "我": true, "以": true, "要": true, "他": true,
		"时": true, "来": true, "用": true, "们": true, "生": true,
		"到": true, "作": true, "天": true, "能": true, "对": true,
		"小": true, "多": true, "然": true, "于": true, "心": true,
		"学": true, "说": true, "过": true, "会": true, "可": true,
	},
}

// NewStopwords creates a Stopwords instance for the given language.
func NewStopwords(lang string) *Stopwords {
	words, ok := defaultDicts[lang]
	if !ok {
		words = defaultDicts["en"]
	}
	return &Stopwords{lang: lang, words: words}
}

// IsStopword checks if a word is a stopword.
func (s *Stopwords) IsStopword(word string) bool {
	return s.words[word]
}

// Words returns the underlying stopword map.
func (s *Stopwords) Words() map[string]bool {
	return s.words
}

// Count returns the number of stopwords.
func (s *Stopwords) Count() int {
	return len(s.words)
}

// Language returns the language code.
func (s *Stopwords) Language() string {
	return s.lang
}

// Default returns the Stopwords for the specified language.
func Default(lang string) *Stopwords {
	return NewStopwords(lang)
}

// Contains checks if a word is a stopword for the specified language.
func Contains(lang, word string) bool {
	return Default(lang).IsStopword(word)
}

// Remove removes a word from the stopword list for the specified language.
func Remove(lang, word string) {
	if words, ok := defaultDicts[lang]; ok {
		delete(words, word)
	}
}

// Add adds a word to the stopword list for the specified language.
func Add(lang, word string) {
	if words, ok := defaultDicts[lang]; ok {
		words[word] = true
	}
}

// Languages returns a list of supported languages.
func Languages() []string {
	return []string{"en", "es", "pt", "fr", "de", "it", "nl", "ru", "zh", "ja"}
}

// Count returns the number of stopwords for the specified language.
func Count(lang string) int {
	return len(Default(lang).Words())
}
