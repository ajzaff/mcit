package mcit

import (
	"iter"
	"math"
	"slices"
)

const exploreTerm = math.Pi

type node struct {
	height      int
	action      string
	parent      *node
	runs        float64
	value       float64
	ucb1        float64
	prior       float64
	frontierIdx int
	children    map[string]*node
}

func newRoot() *node { return &node{frontierIdx: -1, prior: 1} }

func (parent *node) newChild(action string) *node {
	n := &node{
		parent:      parent,
		height:      parent.height + 1,
		frontierIdx: -1,
		prior:       1,
		action:      action,
	}
	if parent.children == nil {
		parent.children = make(map[string]*node)
	}
	parent.children[action] = n
	return n
}

func (n *node) recomputeUCB1() { n.ucb1 = n.computeUCB1() }

func (n *node) computeUCB1() float64 {
	if n.runs == 0 {
		return math.Inf(+1)
	}
	return (n.value + n.prior*exploreTerm) / n.runs
}

func (n *node) score() float64 {
	if n.runs == 0 {
		return math.Inf(-1)
	}
	return n.value / n.runs
}

func (n *node) appendLine(buf []string) []string {
	if n == nil {
		return buf
	}
	i := len(buf)
	buf = slices.Grow(buf[i:], 1+n.height)[:i+1+n.height]
	for i := len(buf) - 1; n.parent != nil; i, n = i-1, n.parent {
		buf[i] = n.action
	}
	return buf
}

type byUCB1 []*node

func (a byUCB1) Len() int { return len(a) }
func (a byUCB1) Swap(i, j int) {
	a[i].frontierIdx, a[j].frontierIdx = a[j].frontierIdx, a[i].frontierIdx
	a[i], a[j] = a[j], a[i]
}
func (a byUCB1) Less(i, j int) bool {
	if ui, uj := a[i].ucb1, a[j].ucb1; ui != uj {
		return ui < uj
	}
	// When priorities are equal (often +∞), we prioritize nodes closer to root.
	// This effectively implements BFS and results in more balanced search performance.
	return a[i].height < a[j].height
}
func (a *byUCB1) Push(x any) {
	n := len(*a)
	e := x.(*node)
	*a = append(*a, e)
	e.recomputeUCB1()
	e.frontierIdx = n
}
func (a *byUCB1) Pop() any {
	n := len(*a) - 1
	x := (*a)[n]
	*a = (*a)[:n]
	x.frontierIdx = -1
	return x
}

func maxNode(a, b *node) *node {
	if a == nil || a.score() < b.score() {
		return b
	}
	return a
}

func minNode(a, b *node) *node {
	if a == nil || b.score() < a.score() {
		return b
	}
	return a
}

func mostPopularNode(a, b *node) *node {
	if a == nil || a.runs < b.runs {
		return b
	}
	return a
}

func selectChildFunc(selectFn func(a, b *node) *node) func(*node) *node {
	return func(root *node) *node {
		var selectNode *node
		for _, child := range root.children {
			selectNode = selectFn(selectNode, child)
		}
		return selectNode
	}
}

func getSelectLine(root *node, selectFn func(*node) *node) *node {
	for len(root.children) > 0 {
		next := selectFn(root)
		if next == nil {
			break
		}
		root = next
	}
	return root
}

func chooseNode(chosen, root *node, chooseFn func(a, b *node) *node) *node {
	n := chooseFn(chosen, root)
	for _, child := range root.children {
		n = chooseNode(n, child, chooseFn)
	}
	return n
}

func nodeIter(root *node) iter.Seq[*node] {
	return func(yield func(*node) bool) { visitNodes(root, yield) }
}

func visitNodes(root *node, visitFn func(*node) bool) {
	if !visitFn(root) {
		return
	}
	for _, child := range root.children {
		visitNodes(child, visitFn)
	}
}
