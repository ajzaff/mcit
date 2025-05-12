package main

import (
	"fmt"
	"iter"
	"math/rand/v2"
	"slices"

	"github.com/ajzaff/fastlog"
	"github.com/ajzaff/fastlog/suite"
	mcts "github.com/ajzaff/mcit"
	"github.com/ajzaff/mcit/perft"
	"github.com/ajzaff/mcit/variation"
)

const kLo, kHi = -2, 2

type constRange struct {
	lo float32
	hi float32
}

func (c constRange) Range() float32 { return c.hi - c.lo }
func (c constRange) Mid() float32   { return c.Range()/2 + c.lo }

func (c constRange) Lo() (lo constRange) {
	mid := c.Mid()
	return constRange{c.lo, mid}
}

func (c constRange) Hi() (hi constRange) {
	mid := c.Mid()
	return constRange{mid, c.hi}
}

func (c constRange) String() string { return fmt.Sprintf("[%.16f,%.16f)", c.lo, c.hi) }

type search struct {
	k constRange
}

func (s *search) Reset(actions iter.Seq[string]) {
	s.k = constRange{kLo, kHi}

	// Apply actions.
	for a := range actions {
		switch a {
		case "lo":
			s.k = s.k.Lo()
		case "hi":
			s.k = s.k.Hi()
		}
	}
}

func main() {
	r := rand.New(rand.NewPCG(1336, 1338))

	result := mcts.Search(func(c *mcts.Context) {
		var s search
		s.Reset(c.Actions())

		// Calculate suite MSE.
		k := s.k.Mid()
		mse := suite.CalculateLog2MSE(1+k, -k)

		// Minimize MSE
		c.Minimize()
		c.SetResultValue(mse)
		c.Expand("lo", "hi")
	}, mcts.RandSource(r), mcts.MaxIters(1_000))

	maxNode := perft.Max(result.Root, func(n *mcts.Node, stat mcts.Stat) float32 {
		return stat.Score()
	})

	line := variation.Line(maxNode)

	fmt.Println(variation.Stat(maxNode))
	fmt.Println(line)
	fmt.Println()

	var s search
	s.Reset(slices.Values(line))

	k := s.k.Mid()
	mse := suite.CalculateLog2MSE(1+k, -k)
	stats := perft.DetailedSearchStats(result.Root)

	fmt.Println("k:           ", k)
	fmt.Println("c0, c1:      ", 1+k, -k)
	fmt.Println("suiteMSE:    ", mse)
	fmt.Println("iterations:  ", result.Iterations)
	fmt.Println("duration:    ", result.Duration)
	fmt.Println("node_count:  ", stats.NodeCount)
	fmt.Println("leaf_count:  ", stats.LeafCount)
	fmt.Println("height:      ", stats.Height)
	fmt.Println("deepest_run: ", stats.DeepestRun)

	bestMSE := suite.CalculateLog2MSE(1+fastlog.K, -fastlog.K)
	if mse < bestMSE {
		fmt.Println()
		fmt.Println("New c0, c1 MSE beats the best known k !!")
	}
}
