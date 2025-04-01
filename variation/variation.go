package visitor

import (
	"math/rand/v2"

	"github.com/ajzaff/mcit"
)

func getSelectLine(root *mcit.Node, selectFn func(*mcit.Node) *mcit.Node) *mcit.Node {
	for len(root.Children) > 0 {
		next := selectFn(root)
		if next == nil {
			break
		}
		root = next
	}
	return root
}

func selectChildFunc(r *rand.Rand, cmpFn func(a, b mcit.Stat) int) func(*mcit.Node) *mcit.Node {
	return func(root *mcit.Node) *mcit.Node {
		// Create an equivalence slice for implementing fair random choice
		// To tie break between equivalent children according to cmpFn.
		equal := []mcit.Stat{{}}
		for _, b := range root.Bandits {
			a := equal[0]
			switch c := cmpFn(a, b); {
			case c > 0:
				// Swap out for a better node, and reset the equivalence slice.
				equal[0] = b
				equal = equal[:1]
			case c == 0:
				// When equal, add it to the equivalence slice.
				equal = append(equal, b)
			}
		}
		if len(equal) > 1 {
			// Fair random choice between equivalent children.
			i := r.IntN(len(equal))
			equal[0], equal[i] = equal[i], equal[0]
		}
		// When nil, no action was selected.
		return root.Children[equal[0].Action]
	}
}

func compareMaxStat(a, b mcit.Stat) int {
	as, bs := a.Score(), b.Score()
	if as < bs {
		return +1
	}
	if as > bs {
		return -1
	}
	return 0
}

func compareMinStat(a, b mcit.Stat) int {
	as, bs := a.Score(), b.Score()
	if as < bs {
		return -1
	}
	if as > bs {
		return +1
	}
	return 0
}

func compareStatPopularity(a, b mcit.Stat) int {
	ar, br := a.Runs, b.Runs
	if ar < br {
		return +1
	}
	if ar > br {
		return -1
	}
	return 0
}

func MaxVariation(root *mcit.Node, r *rand.Rand) *mcit.Node {
	return getSelectLine(root, selectChildFunc(r, compareMaxStat))
}
func MinVariation(root *mcit.Node, r *rand.Rand) *mcit.Node {
	return getSelectLine(root, selectChildFunc(r, compareMinStat))
}
func MostPopularVariation(root *mcit.Node, r *rand.Rand) *mcit.Node {
	return getSelectLine(root, selectChildFunc(r, compareStatPopularity))
}
