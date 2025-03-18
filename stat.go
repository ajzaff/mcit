package mcit

import (
	"math"
	"slices"
)

type NodeStat struct {
	Parent   *NodeStat
	Children map[string]*NodeStat

	Action      string
	Height      int
	Priority    float64
	Prior       float64
	Runs        float64
	Value       float64
	frontierIdx int
}

func newShallowNodeStat(n *node) *NodeStat {
	s := new(NodeStat)
	s.reset(n)
	return s
}

func newVariationStat(n *node) *NodeStat {
	s := new(NodeStat)
	s.reset(n)
	if n.parent != nil {
		s.Parent = newVariationStat(n.parent)
	}
	return s
}

func newFullNodeStat(n *node) *NodeStat {
	s := new(NodeStat)
	s.reset(n)
	for _, child := range n.children {
		s.newSubtree(child)
	}
	return s
}

func newRootStat() *NodeStat { return &NodeStat{frontierIdx: -1, Prior: 1} }

func (parent *NodeStat) NewChild(action string) *NodeStat {
	if child, found := parent.Children[action]; found {
		return child
	}
	n := &NodeStat{
		Parent:      parent,
		Height:      parent.Height + 1,
		frontierIdx: -1,
		Prior:       1,
		Action:      action,
	}
	if parent.Children == nil {
		parent.Children = map[string]*NodeStat{}
	}
	parent.Children[action] = n
	return n
}

func (parent *NodeStat) newSubtree(n *node) *NodeStat {
	s := new(NodeStat)
	s.Parent = parent
	s.reset(n)
	s.Height = parent.Height + 1
	for _, child := range n.children {
		s.newSubtree(child)
	}
	if parent.Children == nil {
		parent.Children = make(map[string]*NodeStat)
	}
	parent.Children[n.action] = s
	return s
}

func (s *NodeStat) reset(n *node) {
	s.Action = n.action
	s.Height = n.height
	s.Priority = n.ucb1
	s.Prior = n.prior
	s.Runs = n.runs
	s.Value = n.value
	s.frontierIdx = n.frontierIdx
}

func (n *NodeStat) Exhausted() bool { return n.frontierIdx == -1 }

func (n *NodeStat) Score() float64 {
	if n.Runs == 0 {
		return math.Inf(-1)
	}
	return n.Value / n.Runs
}

func (n *NodeStat) RecomputePriority() { n.Priority = n.ComputePriority() }

func (n *NodeStat) ComputePriority() float64 {
	if n.Runs == 0 {
		return math.Inf(+1)
	}
	return (n.Value + n.Prior*exploreTerm) / n.Runs
}

func (s *NodeStat) Reset() {
	s.Parent = nil
	s.Action = ""
	s.Height = -1
	s.Priority = math.Inf(-1)
	s.Runs = 0
}

func (s *NodeStat) Line() []string { return s.AppendLine(nil) }

func (s *NodeStat) AppendLine(buf []string) []string {
	i := len(buf)
	buf = slices.Grow(buf[i:], 1+s.Height)
	for ; s.Parent != nil; s = s.Parent {
		buf = append(buf, s.Action)
	}
	slices.Reverse(buf[i:])
	return buf
}

func (s *NodeStat) Hist(hist Hist, valueFn func(*NodeStat) float64) {
	x := valueFn(s)
	hist.Insert(x)
	for _, child := range s.Children {
		child.Hist(hist, valueFn)
	}
}

type SearchStats struct {
	MinSampledPriority float64
	Iterations         int64
	NodeCount          int64
	LeafCount          int64
	MaxFrontierSize    int64
	ExhaustedNodes     int64
	MaxDepthRun        int64
	MaxDepth           int64
}

func newSearchStats() *SearchStats {
	return &SearchStats{
		MinSampledPriority: math.Inf(-1),
	}
}
