package perft

import (
	"math"

	"github.com/ajzaff/lazyq"
	"github.com/ajzaff/mcts"
)

// Reduce a series of measures in bulk on Nodes.
func Reduce[T int64 | float32 | []int64 | []float32 | Hist[int64] | Hist[float32]](root *mcts.Node, v0 T, reduceFn func(*mcts.Node, T) T) T {
	v := v0
	for n := range NodeSeq(root) {
		v = reduceFn(n, v)
	}
	return v
}

// ReduceChild reduces a series of node children.
func ReduceChild[T int64 | float32 | []int64 | []float32 | Hist[int64] | Hist[float32]](root *mcts.Node, v0 T, reduceFn func(*mcts.Node, mcts.Child, T) T) T {
	v := v0
	for n := range NodeSeq(root) {
		for s := range lazyq.Payloads(n.Queue) {
			v = reduceFn(n, s, v)
		}
	}
	return v
}

func Min(root *mcts.Node, valueFn func(*mcts.Node, mcts.Stat) float32) *mcts.Node {
	var (
		v0      = float32(math.Inf(+1))
		minNode *mcts.Node
	)
	ReduceChild(root, v0, func(n *mcts.Node, c mcts.Child, minValue float32) float32 {
		v := valueFn(n, c.Stat)
		if v < minValue {
			minNode = n
			return v
		}
		return minValue
	})
	return minNode
}

func Max(root *mcts.Node, valueFn func(*mcts.Node, mcts.Stat) float32) *mcts.Node {
	var (
		v0      = float32(math.Inf(-1))
		maxNode *mcts.Node
	)
	ReduceChild(root, v0, func(n *mcts.Node, c mcts.Child, maxValue float32) float32 {
		v := valueFn(n, c.Stat)
		if maxValue < v {
			maxNode = n
			return v
		}
		return maxValue
	})
	return maxNode
}
