package variation

import (
	"math/rand/v2"

	"github.com/ajzaff/lazyq"
	"github.com/ajzaff/mcit"
)

func getSelectLine(root *mcit.Node, selectFn func(*mcit.Node) *mcit.Node) *mcit.Node {
	for root.Queue.Len() > 0 {
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
		equal := []mcit.Child{{}}
		for b := range lazyq.Payloads(root.Queue) {
			a := equal[0]
			switch c := cmpFn(a.Stat, b.Stat); {
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
		e, _ := LookupElem(root, equal[0].Action)
		return e.E.Node
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

// Variation returns the node accessed from root by the given line or nil.
func Variation(root *mcit.Node, line ...string) *mcit.Node {
	for _, a := range line {
		if root == nil {
			return nil
		}
		e, ok := LookupElem(root, a)
		if !ok {
			return nil
		}
		root = e.E.Node
	}
	return root
}

// Stat returns the stat accessed from root by the given line or empty.
func Stat(root *mcit.Node, line ...string) mcit.Stat {
	n := Variation(root, line...)
	s := n.Stat()
	if s == nil {
		return mcit.Stat{}
	}
	return s.Stat
}
