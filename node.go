package mcit

import (
	"iter"
	"math"
)

const exploreTerm = math.Pi

type byUCB1 []*NodeStat

func (a byUCB1) Len() int { return len(a) }
func (a byUCB1) Swap(i, j int) {
	a[i].frontierIdx, a[j].frontierIdx = a[j].frontierIdx, a[i].frontierIdx
	a[i], a[j] = a[j], a[i]
}
func (a byUCB1) Less(i, j int) bool {
	if ui, uj := a[i].Priority, a[j].Priority; ui != uj {
		return ui < uj
	}
	if hi, hj := a[i].Height, a[j].Height; hi != hj {
		// When priorities are equal (often +∞), we prioritize nodes closer to root.
		// This effectively implements BFS and results in more balanced search performance.
		return hi < hj
	}
	// Fall back to prior comparison.
	return a[i].Prior > a[j].Prior
}
func (a *byUCB1) Push(x any) {
	n := len(*a)
	e := x.(*NodeStat)
	*a = append(*a, e)
	e.RecomputePriority()
	e.frontierIdx = n
}
func (a *byUCB1) Pop() any {
	n := len(*a) - 1
	x := (*a)[n]
	*a = (*a)[:n]
	x.frontierIdx = -1
	return x
}

func maxNode(a, b *NodeStat) *NodeStat {
	if a == nil || a.Score() < b.Score() {
		return b
	}
	return a
}

func minNode(a, b *NodeStat) *NodeStat {
	if a == nil || b.Score() < a.Score() {
		return b
	}
	return a
}

func mostPopularNode(a, b *NodeStat) *NodeStat {
	if a == nil || a.Runs < b.Runs {
		return b
	}
	return a
}

func selectChildFunc(selectFn func(a, b *NodeStat) *NodeStat) func(*NodeStat) *NodeStat {
	return func(root *NodeStat) *NodeStat {
		var selectNode *NodeStat
		for _, child := range root.Children {
			selectNode = selectFn(selectNode, child)
		}
		return selectNode
	}
}

func getSelectLine(root *NodeStat, selectFn func(*NodeStat) *NodeStat) *NodeStat {
	for len(root.Children) > 0 {
		next := selectFn(root)
		if next == nil {
			break
		}
		root = next
	}
	return root
}

func chooseNode(chosen, root *NodeStat, chooseFn func(a, b *NodeStat) *NodeStat) *NodeStat {
	n := chooseFn(chosen, root)
	for _, child := range root.Children {
		n = chooseNode(n, child, chooseFn)
	}
	return n
}

func nodeIter(root *NodeStat) iter.Seq[*NodeStat] {
	return func(yield func(*NodeStat) bool) { visitNodes(root, yield) }
}

func visitNodes(root *NodeStat, visitFn func(*NodeStat) bool) {
	if !visitFn(root) {
		return
	}
	for _, child := range root.Children {
		visitNodes(child, visitFn)
	}
}
