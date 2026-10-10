package wsd

import "strings"

func ContextWords(text string) []string { return strings.Fields(strings.ToLower(text)) }
