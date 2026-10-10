package inference

import "github.com/ateto3/gopherml/go-nlp/sem"

type Command struct {
	Goal        sem.Expr
	Assumptions []sem.Expr
	Prover      Prover
	Result      bool
	Proof       string
}

func NewCommand(goal sem.Expr, assumptions []sem.Expr) *Command {
	return &Command{Goal: goal, Assumptions: assumptions, Prover: NewResolutionProver()}
}

func (c *Command) Prove() bool {
	ok, proof := c.Prover.Prove(c.Goal, c.Assumptions)
	c.Result, c.Proof = ok, proof
	return ok
}
