package variation

import (
	"math/rand/v2"

	"github.com/ajzaff/lazyq"
	"github.com/ajzaff/mcts"
)

func getSelectLine(root *mcts.Node, selectFn func(*mcts.Node) *mcts.Node) *mcts.Node {
	for root.Queue.Len() > 0 {
		next := selectFn(root)
		if next == nil {
			break
		}
		root = next
	}
	return root
}

func selectChildFunc(r *rand.Rand, cmpFn func(a, b mcts.Stat) int) func(*mcts.Node) *mcts.Node {
	return func(root *mcts.Node) *mcts.Node {
		// Create an equivalence slice for implementing fair random choice
		// To tie break between equivalent children according to cmpFn.
		equal := []mcts.Child{{}}
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

func compareMaxStat(a, b mcts.Stat) int {
	as, bs := a.Score(), b.Score()
	if as < bs {
		return +1
	}
	if as > bs {
		return -1
	}
	return 0
}

func compareMinStat(a, b mcts.Stat) int {
	as, bs := a.Score(), b.Score()
	if as < bs {
		return -1
	}
	if as > bs {
		return +1
	}
	return 0
}

func compareStatPopularity(a, b mcts.Stat) int {
	ar, br := a.Runs, b.Runs
	if ar < br {
		return +1
	}
	if ar > br {
		return -1
	}
	return 0
}

func MaxVariation(root *mcts.Node, r *rand.Rand) *mcts.Node {
	return getSelectLine(root, selectChildFunc(r, compareMaxStat))
}
func MinVariation(root *mcts.Node, r *rand.Rand) *mcts.Node {
	return getSelectLine(root, selectChildFunc(r, compareMinStat))
}
func MostPopularVariation(root *mcts.Node, r *rand.Rand) *mcts.Node {
	return getSelectLine(root, selectChildFunc(r, compareStatPopularity))
}

// Variation returns the node accessed from root by the given line or nil.
func Variation(root *mcts.Node, line ...string) *mcts.Node {
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
func Stat(root *mcts.Node, line ...string) mcts.Stat {
	n := Variation(root, line...)
	s := n.Stat()
	if s == nil {
		return mcts.Stat{}
	}
	return s.Stat
}
