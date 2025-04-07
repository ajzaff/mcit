package mcit

import (
	"math"

	"github.com/ajzaff/mcit/internal/fastlog"
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

// Remove the ith Bandit from the Node.
func (n *Node) Remove(i int) {
	n.Trials -= n.Bandits[i].Runs
	n.remove(i) // Remove the stat while considering lazy elements.
}

// Clear zeroes the statistics of the ith bandit in the Node
// and decrements the Node's trial counter.
func (n *Node) Clear(i int) {
	n.Trials -= n.Bandits[i].Runs
	n.clear(i) // Fix the lazy heap index.
}

// Clear zeroes the statistics from the Stat as if it were never run.
//
// It does not clear Action or Prior values.
func (s *Stat) Clear() {
	s.Priority = float32(math.Inf(+1))
	s.Runs = 0
	s.Value = 0
}

// AddValueRuns adds val and runs to the ith bandit statistics
// and the nodes Trials counter.
//
// AddValueRuns correctly handles the node's Minimize flag.
func (n *Node) AddValueRuns(i int, val, runs float32) {
	if n.Minimize {
		// Negate minimizing nodes (min(a,b) = -max(-a,-b)).
		val = -val
	}
	n.Bandits[i].Value += val
	n.Bandits[i].Runs += runs
	n.Trials += runs
}

// Score the stat on the node taking into account the Minimize flag.
func (n *Node) Score(stat Stat) float32 {
	v := stat.Score()
	if n.Minimize {
		return -v
	}
	return v
}

// Score returns the raw score statistic on the Stat.
//
// It does not take into account the Node's Minimize flag.
func (s Stat) Score() float32 {
	if s.Runs == 0 {
		return float32(math.Inf(-1))
	}
	return s.Value / s.Runs
}

// recomputePriority updates the PUCT policy value for the ith bandit with the current statistics.
//
// recomputePriority should only be called when runs > 0, otherwise it returns NaN.
func (n *Node) recomputePriority(i int, exploreFactor float32) {
	bandit := n.Bandits[i]
	n.Bandits[i].Priority = computePriority(bandit.Value, bandit.Prior, bandit.Runs, n.Trials, exploreFactor)
}

// computePriority computes the PUCT formula on the inputs.
func computePriority(value, prior, runs, trials, exploreFactor float32) float32 {
	runFactor := runs + 1
	runFactor = 1 / runFactor
	exploit := value * runFactor
	explore := fastlog.Log(trials) * runFactor
	explore = float32(math.Sqrt(float64(explore)))
	explore *= prior * exploreFactor
	return exploit + explore
}
