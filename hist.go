package mcit

import (
	"math"
	"slices"
)

var (
	runBins      = []float64{0, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, math.Inf(+1)}
	scoreBins    = []float64{-math.MaxFloat64, -1, -0.8, -0.6, -0.4, -0.2, 0, 0.2, 0.4, 0.6, 0.8, 1, math.Inf(+1)}
	priorityBins = []float64{
		-math.MaxFloat64, -1 * exploreTerm, -0.8 * exploreTerm, -0.6 * exploreTerm, -0.4 * exploreTerm, -0.2 * exploreTerm,
		0,
		0.2 * exploreTerm, 0.4 * exploreTerm, 0.6 * exploreTerm, 0.8 * exploreTerm, 1 * exploreTerm, math.Inf(+1),
	}
)

type HistBin struct {
	Max   float64
	Count int64
}

type Hist struct {
	Bins []HistBin
}

func makeHist(maxBinValues []float64) Hist {
	bins := make([]HistBin, len(maxBinValues))
	for i, v := range maxBinValues {
		bins[i].Max = v
	}
	return Hist{bins}
}

func (h Hist) Insert(x float64) {
	i, _ := slices.BinarySearchFunc(h.Bins, HistBin{Max: x}, func(a, b HistBin) int {
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

func (h Hist) Remove(x float64) {
	i, _ := slices.BinarySearchFunc(h.Bins, HistBin{Max: x}, func(a, b HistBin) int {
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
