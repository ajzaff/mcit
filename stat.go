package mcit

import (
	"math"
	"slices"
)

type NodeStat struct {
	Parent   *NodeStat
	Children map[string]*NodeStat

	Action   string
	Height   int
	Priority float64
	Prior    float64
	Runs     float64
	Value    float64
	Score    float64
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

func newNodeStat(n *node) *NodeStat {
	s := new(NodeStat)
	s.reset(n)
	for _, child := range n.children {
		s.newChild(child)
	}
	return s
}

func (parent *NodeStat) newChild(n *node) *NodeStat {
	s := new(NodeStat)
	s.Parent = parent
	s.reset(n)
	for _, child := range n.children {
		s.newChild(child)
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
	s.Runs = n.results.Count
	s.Value = n.results.Value
	s.Score = n.score()
}

func (s *NodeStat) Reset() {
	s.Parent = nil
	s.Action = ""
	s.Height = -1
	s.Priority = math.Inf(-1)
	s.Runs = 0
	s.Score = math.Inf(-1)
}

func (s *NodeStat) Line() []string {
	buf := make([]string, 0, 1+s.Height)
	return s.AppendLine(buf)
}

func (s *NodeStat) AppendLine(buf []string) []string {
	i := len(buf)
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
}
