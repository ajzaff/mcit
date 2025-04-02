package mcit

import (
	"math"
)

// Stat defines a structure for statistics used in the multi-armed bandit process.
//
// Stat is embedded inside a Node.
type Stat struct {
	Action   string
	Priority float32
	Prior    float32
	Runs     float32
	Value    float32
}

func (n *Node) AddValueRuns(i int, val, runs float32) {
	if n.Minimize {
		// Negate minimizing nodes (min(a,b) = -max(-a,-b)).
		val = -val
	}
	n.Bandits[i].Value += val
	n.Bandits[i].Runs += runs
}

func (n Stat) Score() float32 {
	if n.Runs == 0 {
		return float32(math.Inf(-1))
	}
	return n.Value / n.Runs
}

func (n *Node) RecomputePriority(i int, exploreFactor float32) {
	n.Bandits[i].Priority = n.Bandits[i].ComputePriority(exploreFactor)
}

func (n Stat) ComputePriority(exploreFactor float32) float32 {
	if n.Runs == 0 {
		return float32(math.Inf(+1))
	}
	return (n.Value + n.Prior*exploreFactor) / n.Runs
}

func (s *Stat) Reset() {
	s.Action = ""
	s.Priority = float32(math.Inf(-1))
	s.Runs = 0
}
