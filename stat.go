package mcit

import (
	"math"

	"github.com/ajzaff/fastlog"
	"github.com/ajzaff/lazyq"
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

// Clear zeroes the statistics from the Stat as if it were never run.
//
// It does not clear Action or Prior values.
func (s *Stat) Clear() {
	s.Priority = float32(math.Inf(+1))
	s.Runs = 0
	s.Value = 0
}

// addValueRuns adds val and runs to the top bandit Stat
// and updates the node's Trials counter.
//
// addValueRuns correctly handles the node's Minimize flag.
//
// We expect to call recomputePriority afterwards.
func (n *Node) addValueRuns(val, runs float32) {
	if n.Minimize {
		// Negate minimizing nodes (min(a,b) = -max(-a,-b)).
		val = -val
	}
	e := lazyq.First(n.Queue)
	e.Value += val
	e.Runs += runs
	lazyq.ReplacePayload(n.Queue, 0, e)
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
// Maintains the queue's heap property.
//
// recomputePriority should only be called when runs > 0, otherwise it returns NaN.
func (n *Node) recomputePriority(exploreFactor float32) {
	bandit := lazyq.First(n.Queue)
	n.Queue.Decrease(computePriority(bandit.Value, bandit.Prior, bandit.Runs, n.Trials, exploreFactor))
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
