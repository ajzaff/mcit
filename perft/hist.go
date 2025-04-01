package perft

import (
	"math"
	"slices"

	"github.com/ajzaff/mcit"
)

const exploreTerm = 2 * math.Pi

var (
	runBins      = []float64{0, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, math.Inf(+1)}
	scoreBins    = []float64{-math.MaxFloat64, -1, -0.8, -0.6, -0.4, -0.2, 0, 0.2, 0.4, 0.6, 0.8, 1, math.Inf(+1)}
	priorityBins = []float64{
		-math.MaxFloat64, -1 * exploreTerm, -0.8 * exploreTerm, -0.6 * exploreTerm, -0.4 * exploreTerm, -0.2 * exploreTerm,
		0,
		0.2 * exploreTerm, 0.4 * exploreTerm, 0.6 * exploreTerm, 0.8 * exploreTerm, 1 * exploreTerm, math.Inf(+1),
	}
)

func DefaultRunBins() []float64      { return slices.Clone(runBins) }
func DefaultScoreBins() []float64    { return slices.Clone(scoreBins) }
func DefaultPriorityBins() []float64 { return slices.Clone(priorityBins) }

type HistBin[T int64 | float32] struct {
	Max   T
	Count int64
}

type Hist[T int64 | float32] struct {
	Bins []HistBin[T]
}

func MakeHist[T int64 | float32](bins []T) Hist[T] {
	b := make([]HistBin[T], len(bins))
	for i, v := range bins {
		b[i].Max = v
	}
	return Hist[T]{Bins: b}
}

func Fill[T int64 | float32](root *mcit.Node, hist Hist[T], valueFn func(mcit.Stat) T) {
	for n := range NodeSeq(root) {
		for e := range n.StatSeq() {
			x := valueFn(e)
			hist.Insert(x)
		}
	}
}

func (h Hist[T]) Insert(x T) {
	i, _ := slices.BinarySearchFunc(h.Bins, HistBin[T]{Max: x}, func(a, b HistBin[T]) int {
		if a == b {
			return 0
		}
		if a.Max < b.Max {
			return -1
		}
		return +1
	})
	h.Bins[i].Count++
}

func (h Hist[T]) Remove(x T) {
	i, _ := slices.BinarySearchFunc(h.Bins, HistBin[T]{Max: x}, func(a, b HistBin[T]) int {
		if a == b {
			return 0
		}
		if a.Max < b.Max {
			return -1
		}
		return +1
	})
	if h.Bins[i].Count > 0 {
		h.Bins[i].Count--
	}
}
