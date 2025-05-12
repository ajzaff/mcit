package variation

import (
	"github.com/ajzaff/lazyq"
	"github.com/ajzaff/mcit"
)

// Depth calculates the number of nodes between n and root.
func Depth(n *mcit.Node) int {
	if n.Parent == nil {
		return 0
	}
	return Depth(n.Parent) + 1
}

// Detatched returns a clone of the stat object detatched from patents and children
// without modifying the original stat object.
func Detatched(n *mcit.Node) *mcit.Node {
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
func LookupElem(n *mcit.Node, action string) (lazyq.Elem[mcit.Child], bool) {
	for _, e := range lazyq.ElementIndices(n.Queue) {
		if e.E.Action == action {
			return e, true
		}
	}
	return lazyq.Elem[mcit.Child]{}, false
}
