package featstruct

import "strings"

func Parse(s string) FeatStruct {
	s = strings.TrimSpace(strings.Trim(s, "[]"))
	out := FeatStruct{}
	if s == "" { return out }
	for _, part := range strings.Fields(s) {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 { out[kv[0]] = strings.Trim(kv[1], "[]") }
	}
	return out
}
