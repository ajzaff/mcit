package perft

import (
	"iter"

	"github.com/ajzaff/mcit"
)

func visitNodes(root *mcit.Node, visitFn func(*mcit.Node) bool) {
	if !visitFn(root) {
		return
	}
	for e := range root.StatSeq() {
		visitNodes(root.Children[e.Action], visitFn)
	}
}

// NodeSeq returns an iterator over all nodes under root recursively in descending priority order.
func NodeSeq(root *mcit.Node) iter.Seq[*mcit.Node] {
	return func(yield func(*mcit.Node) bool) { visitNodes(root, yield) }
}
