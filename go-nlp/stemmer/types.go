package stemmer

// Stemmer reduces words to stems.
type Stemmer interface{ Stem(string) string }
