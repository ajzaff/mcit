package mcit

import (
	"math/rand/v2"
	"testing"
)

func TestSearch(t *testing.T) {
	results := Search(func(selector NodeSelector) (results RunResults) {
		if len(selector.Actions) < 10 {
			results.Expand = []string{
				"a",
				"b",
			}
			results.Priors = []float32{
				1,
				1,
			}

			// Random shuffle.
			if rand.Int()&1 == 0 {
				results.Expand[0], results.Expand[1] = results.Expand[1], results.Expand[0]
			}
		}

		// Provide a slight incentive to picking "a" over "b".
		lambda := 1.0
		for _, a := range selector.Actions {
			if a == "a" {
				lambda *= 0.99
			} else if a == "b" {
				lambda *= 1.01
			}
		}

		const experiments = 100

		value := float32(0)
		for range experiments {
			value += float32(rand.ExpFloat64() / lambda)
		}

		results.Replace = len(selector.Actions) > 0
		results.Count = experiments
		results.Value = value

		return
	}, MaxIters(100))

	t.Log(results.Root.lazyQueue.Bandits)
	t.Log(results.Iterations)
	t.Log(results.Duration)
}

func TestSearchFloatRange(t *testing.T) {
	const maxIters = 1000

	x := func(actions []string, loCmd, hiCmd string) float32 {
		lo, hi := -100., 100.

		for _, a := range actions {
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
	results := Search(func(selector NodeSelector) (results RunResults) {
		a := x(selector.Actions, "lo_a", "hi_a")
		b := x(selector.Actions, "lo_b", "hi_b")

		got := objective(a, b)

		results.Value = -loss(got)
		results.Count = 1

		if loss(got) <= epsilon {
			bestA = a
			bestB = b
			results.Done = true
			return
		}

		results.Expand = []string{"lo_a", "hi_a", "hi_b", "lo_b"}

		return
	}, MaxIters(10))

	t.Log("2a^2 + 2b - 100 = 0")
	t.Log("a =", bestA, "b =", bestB, "loss =", loss(objective(bestA, bestB)))
	t.Log(results.Iterations, "iterations")
	t.Log(results.Duration)
}
