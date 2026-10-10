package inference

import "github.com/ateto3/gopherml/go-nlp/sem"

type Prover interface {
	Prove(goal sem.Expr, assumptions []sem.Expr) (bool, string)
}
