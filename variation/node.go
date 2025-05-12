package variation

import (
	"github.com/ajzaff/lazyq"
	"github.com/ajzaff/mcts"
)

// Depth calculates the number of nodes between n and root.
func Depth(n *mcts.Node) int {
	if n.Parent == nil {
		return 0
	}
	return Depth(n.Parent) + 1
}

// Detatched returns a clone of the stat object detatched from patents and children
// without modifying the original stat object.
func Detatched(n *mcts.Node) *mcts.Node {
	copy := *n
	copy.Parent = nil
	copy.Queue = lazyq.Clone(copy.Queue)
	for i, e := range lazyq.ElementIndices(copy.Queue) {
		e.E.Node = nil
		lazyq.ReplacePayload(copy.Queue, i, e.E)
	}
	return &copy
}

// LookupElem searches over all children of n and returns the first marked with action.
func LookupElem(n *mcts.Node, action string) (lazyq.Elem[mcts.Child], bool) {
	for _, e := range lazyq.ElementIndices(n.Queue) {
		if e.E.Action == action {
			return e, true
		}
	}
	return lazyq.Elem[mcts.Child]{}, false
}
