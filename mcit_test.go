package mcit

import (
	"math/rand"
	"testing"
)

func TestSearch(t *testing.T) {
	const maxIters = 100
	iters := 0

	results := Search(func(actions []string) (expand []string, priors []float64, replace bool, results RunResults, done bool) {
		if len(actions) < 10 {
			expand = []string{
				"a",
				"b",
			}
			priors = []float64{
				1,
				1,
			}

			// Random shuffle.
			if rand.Int()&1 == 0 {
				expand[0], expand[1] = expand[1], expand[0]
			}
		}

		// Provide a slight incentive to picking "a" over "b".
		lambda := 1.0
		for _, a := range actions {
			if a == "a" {
				lambda *= 0.99
			} else if a == "b" {
				lambda *= 1.01
			}
		}

		if iters++; iters > maxIters {
			done = true
		}

		const experiments = 100

		value := 0.0
		for range experiments {
			value += rand.ExpFloat64() / lambda
		}

		return expand, priors, true,
			RunResults{
				Count: experiments,
				Value: value,
			}, done
	})

	t.Logf("%#v\n", results)
}
