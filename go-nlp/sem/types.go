package sem

type Expr interface{ Eval(Model, map[string]string) bool }
type Model struct{ Dom map[string]struct{}; Val map[string]any }
