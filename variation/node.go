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
//
// LookupElem returns false when no such entry exists.
func LookupElem(n *mcts.Node, action string) (lazyq.Elem[mcts.Child], bool) {
	for e := range lazyq.Elements(n.Queue) {
		if e.E.Action == action {
			return e, true
		}
	}
	return lazyq.Elem[mcts.Child]{}, false
}

// LookupSelf searches over parent's children to find n's child entry.
//
// LookupSelf returns false when no such entry exists.
func LookupSelf(n *mcts.Node) (lazyq.Elem[mcts.Child], bool) {
	if n == nil || n.Parent == nil {
		return lazyq.Elem[mcts.Child]{}, false
	}
	return LookupElem(n.Parent, n.Action)
}
