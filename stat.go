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
	ExploreFactor float32
	Runs          float32
	Value         float32
}

// Clear zeroes the statistics from the Stat as if it were never run.
//
// It does not clear Action or Prior values.
func (s *Stat) Clear() {
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
	if n.Minimize() {
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
	if n.Minimize() {
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

// logTrials returns an approximation of Log(Trials) for n.
//
// logTrials behavior is undefined when return n.Trials <= 0.
func (n *Node) logTrials() float32 { return fastlog.Log(n.Trials + 1) }

// computePriority computes the PUCT formula on the inputs.
func (s Stat) computePriority(logTrials float32) float32 {
	runFactor := 1 / (s.Runs + 1)
	exploit := s.Value * runFactor
	explore := s.ExploreFactor * float32(math.Sqrt(float64(logTrials*runFactor)))
	return exploit + explore
}
