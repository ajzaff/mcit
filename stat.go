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
	Minimize    bool
}

func newRootStat() *NodeStat { return &NodeStat{frontierIdx: -1, Prior: 1} }

func (parent *NodeStat) NewChild(action string) (child *NodeStat, created bool) {
	if child, found := parent.Children[action]; found {
		return child, false
	}
	child = &NodeStat{
		Parent:      parent,
		Height:      parent.Height + 1,
		frontierIdx: -1,
		Prior:       1,
		Action:      action,
	}
	if parent.Children == nil {
		parent.Children = map[string]*NodeStat{}
	}
	parent.Children[action] = child
	return child, true
}

// Detatched returns a shallow clone of the stat object detatched from patents, children, and the frontier
// without modifying the original stat object.
func (s *NodeStat) Detatched() *NodeStat {
	var copy NodeStat
	copy = *s
	copy.Parent = nil
	copy.Children = nil
	copy.frontierIdx = -1
	return &copy
}

func (n *NodeStat) Exhausted() bool { return n.frontierIdx == -1 }

func (n *NodeStat) Score() float64 {
	if n.Runs == 0 {
		return math.Inf(-1)
	}
	return n.Value / n.Runs
}

// ConvexScore maps all finite scores v to v*v.
func (n *NodeStat) ConvexScore() float64 {
	v := n.Score()
	if math.IsInf(v, 0) {
		return v
	}
	return v * v
}

func (n *NodeStat) RecomputePriority() { n.Priority = n.ComputePriority() }

func (n *NodeStat) ComputePriority() float64 {
	if n.Runs == 0 {
		return math.Inf(+1)
	}
	value := n.Value
	if n.Minimize { // Negate minimizing nodes (min(a,b) = -max(-a,-b)).
		value = -value
	}
	return (value + n.Prior*exploreTerm) / n.Runs
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
