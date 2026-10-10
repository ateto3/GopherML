package sem

func Evaluate(expr Expr, m Model) bool { return expr.Eval(m, map[string]string{}) }
