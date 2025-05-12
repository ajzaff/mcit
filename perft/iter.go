package perft

import (
	"iter"

	"github.com/ajzaff/lazyq"
	"github.com/ajzaff/mcit"
)

func visitNodes(root *mcit.Node, depth int, visitFn func(n *mcit.Node, depth int) bool) {
	if root == nil || !visitFn(root, depth) {
		return
	}
	for e := range lazyq.Payloads(root.Queue) {
		visitNodes(e.Node, depth+1, visitFn)
	}
}

// NodeSeq returns an iterator over all nodes under root recursively in descending priority order.
func NodeSeq(root *mcit.Node) iter.Seq[*mcit.Node] {
	return func(yield func(*mcit.Node) bool) {
		visitNodes(root, 0, func(n *mcit.Node, _ int) bool { return yield(n) })
	}
}
