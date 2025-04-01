package perft

import (
	"github.com/ajzaff/mcit"
)

// Reduce a series of measures in bulk on Nodes.
func Reduce[T int64 | float64 | []int64 | []float64 | Hist[int64] | Hist[float64]](root *mcit.Node, v0 T, reduceFn func(*mcit.Node, T) T) T {
	v := v0
	for n := range NodeSeq(root) {
		v = reduceFn(n, v)
	}
	return v
}
