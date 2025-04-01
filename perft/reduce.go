package perft

import (
	"github.com/ajzaff/mcit"
)

// Reduce a series of measures in bulk on Nodes.
func Reduce[T int | int64 | float64 | []int | []int64 | []float64](root *mcit.Node, v0 T, countFn func(*mcit.Node, T) T) T {
	v := v0
	for n := range NodeSeq(root) {
		v = countFn(n, v0)
	}
	return v
}

// ReduceStats reduces a series of measures in bulk on Stats.
func ReduceStats[T int | int64 | float64 | []int | []int64 | []float64](root *mcit.Node, v0 T, countFn func(mcit.Stat, T) T) T {
	v := v0
	for n := range NodeSeq(root) {
		for stat := range n.StatSeq() {
			v = countFn(stat, v)
		}
	}
	return v
}
