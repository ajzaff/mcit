package mcts

import (
	"iter"
	"math/rand/v2"
	"testing"

	"github.com/ajzaff/lazyq"
)

func TestSearch(t *testing.T) {
	src := rand.NewPCG(1337, 420)

	results := Search(func(c *Context) {
		if c.Len() < 10 {
			c.Expand("a", "b")
			c.Priors(1, 1)
		}

		// Provide a slight incentive to picking "a" over "b".
		lambda := 1.0
		for a := range c.Actions() {
			if a == "a" {
				lambda *= 0.99
			} else if a == "b" {
				lambda *= 1.01
			}
		}

		const experiments = 100

		for range experiments {
			v := float32(rand.ExpFloat64() / lambda)
			c.AddResultValue(v)
		}
	}, MaxIters(100), RandSource(src))

	t.Log("a score", extractStat(results.Root, "a").Score())
	t.Log("b score", extractStat(results.Root, "b").Score())

	t.Log(results.Root.Queue)
	t.Log(results.Iterations)
	t.Log(results.Duration)
}

func TestSearchFloatRange(t *testing.T) {
	const maxIters = 10

	x := func(actions iter.Seq[string], loCmd, hiCmd string) float32 {
		lo, hi := -100., 100.

		for a := range actions {
			if a == loCmd {
				lo += (hi - lo) / 2
			} else if a == hiCmd {
				hi -= (hi - lo) / 2
			}
		}

		n := lo + (hi-lo)/2
		return float32(n)
	}

	objective := func(a, b float32) float32 { return 2*a*a + 2*b - 100 }

	loss := func(objective float32) float32 { a := 0 - objective; return a * a }

	var (
		bestA float32
		bestB float32
	)

	// Attempts to solve the equation: 2a^2 + 2b - 100 = 0.
	epsilon := float32(0)
	results := Search(func(c *Context) {
		a := x(c.Actions(), "lo_a", "hi_a")
		b := x(c.Actions(), "lo_b", "hi_b")

		got := objective(a, b)

		c.SetResultValue(-loss(got))

		if loss(got) <= epsilon {
			bestA = a
			bestB = b
			c.Stop()
			return
		}

		c.Expand("lo_a", "hi_a", "hi_b", "lo_b")
	}, MaxIters(maxIters))

	t.Log("2a^2 + 2b - 100 = 0")
	t.Log("a =", bestA, "b =", bestB, "loss =", loss(objective(bestA, bestB)))
	t.Log(results.Iterations, "iterations")
	t.Log(results.Duration)
}

// extractVariation is a test helper that mirrors variation.Variation.
func extractVariation(root *Node, line ...string) *Node {
	for _, a := range line {
		if root == nil {
			return nil
		}
		idx, ok := root.indexChild(a)
		if !ok {
			return nil
		}
		s := lazyq.At(root.Queue, idx)
		if s.Node == nil {
			return nil
		}
		root = s.Node
	}
	return root
}

// extractStat is a test helper that mirrors variation.Stat.
func extractStat(root *Node, line ...string) Stat {
	n := extractVariation(root, line...)
	s := n.Stat()
	if s == nil {
		return Stat{}
	}
	return s.Stat
}
