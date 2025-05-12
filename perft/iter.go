package perft

import (
	"iter"

	"github.com/ajzaff/lazyq"
	"github.com/ajzaff/mcts"
)

func visitNodes(root *mcts.Node, depth int, visitFn func(n *mcts.Node, depth int) bool) {
	if root == nil || !visitFn(root, depth) {
		return
	}
	for e := range lazyq.Payloads(root.Queue) {
		visitNodes(e.Node, depth+1, visitFn)
	}
}

// NodeSeq returns an iterator over all nodes under root recursively in descending priority order.
func NodeSeq(root *mcts.Node) iter.Seq[*mcts.Node] {
	return func(yield func(*mcts.Node) bool) {
		visitNodes(root, 0, func(n *mcts.Node, _ int) bool { return yield(n) })
	}
}
