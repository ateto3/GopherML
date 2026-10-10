package featstruct

import "fmt"

func Unify(a, b FeatStruct) (FeatStruct, error) {
	out := New(a)
	for k, v := range b {
		if cur, ok := out[k]; ok {
			mv, err := merge(cur, v)
			if err != nil { return nil, err }
			out[k] = mv
		} else {
			out[k] = clone(v)
		}
	}
	return out, nil
}

func merge(a, b any) (any, error) {
	am, aok := a.(FeatStruct); bm, bok := b.(FeatStruct)
	if aok && bok { return Unify(am, bm) }
	if a == b { return clone(a), nil }
	return nil, fmt.Errorf("feature clash: %v != %v", a, b)
}

func clone(v any) any {
	if m, ok := v.(FeatStruct); ok { return New(m) }
	return v
}
