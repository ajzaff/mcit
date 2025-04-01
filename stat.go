package mcit

import (
	"math"
)

// Stat defines a structure for statistics used in the multi-armed bandit process.
//
// Stat is embedded inside a Node.
type Stat struct {
	Action   string
	Priority float64
	Prior    float64
	Runs     float64
	Value    float64
}

func (n *Node) AddValueRuns(i int, val, runs float64) {
	if n.Minimize {
		// Negate minimizing nodes (min(a,b) = -max(-a,-b)).
		val = -val
	}
	n.Bandits[i].Value += val
	n.Bandits[i].Runs += runs
}

func (n Stat) Score() float64 {
	if n.Runs == 0 {
		return math.Inf(-1)
	}
	return n.Value / n.Runs
}

func (n *Node) RecomputePriority(i int) { n.Bandits[i].Priority = n.Bandits[i].ComputePriority() }

func (n Stat) ComputePriority() float64 {
	if n.Runs == 0 {
		return math.Inf(+1)
	}
	return (n.Value + n.Prior*exploreTerm) / n.Runs
}

func (s *Stat) Reset() {
	s.Action = ""
	s.Priority = math.Inf(-1)
	s.Runs = 0
}
