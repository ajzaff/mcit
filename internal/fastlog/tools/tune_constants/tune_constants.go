package main

import (
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/ajzaff/mcit"
	"github.com/ajzaff/mcit/internal/fastlog"
	"github.com/ajzaff/mcit/perft"
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

func (s *search) Reset(actions []string) {
	s.k = constRange{kLo, kHi}

	// Apply actions.
	for _, a := range actions {
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

	result := mcit.Search(func(sel mcit.NodeSelector) (results mcit.RunResults) {
		var s search
		s.Reset(sel.Actions)

		// Calculate suite MSE.
		k := s.k.Mid()
		mse := fastlog.CalculateSuiteMSE(1+k, -k)

		results.Count = 1
		results.Value = -mse // Maximize(-MSE).
		results.Expand = []string{
			"lo",
			"hi",
		}

		return
	}, mcit.RandSource(r), mcit.MaxIters(1_000_000))

	var maxNode *mcit.Node
	maxNMSE := perft.Reduce(result.Root, math.Inf(-1), func(n *mcit.Node, bestNMSE float64) float64 {
		var bestAction string
		for stat := range n.StatSeq() {
			if nmse := stat.Score(); bestNMSE < nmse {
				bestNMSE = nmse
				bestAction = stat.Action
			}
		}
		if bestNode := n.Children[bestAction]; bestNode != nil {
			maxNode = bestNode
		}
		return bestNMSE
	})

	line := maxNode.Line()

	fmt.Println(maxNMSE)
	fmt.Println(line)
	fmt.Println()

	var s search
	s.Reset(line)

	k := s.k.Mid()
	mse := fastlog.CalculateSuiteMSE(1+k, -k)
	stats := perft.DetailedSearchStats(result.Root)

	fmt.Println("k:              ", k)
	fmt.Println("c0, c1:         ", 1+k, -k)
	fmt.Println("suiteMSE:       ", mse)
	fmt.Println("iterations:     ", result.Iterations)
	fmt.Println("duration:       ", result.Duration)
	fmt.Println("node_count:     ", stats.NodeCount)
	fmt.Println("leaf_count:     ", stats.LeafCount)
	fmt.Println("max_height:     ", stats.MaxHeight)
	fmt.Println("max_height_run: ", stats.MaxHeightRun)

	bestMSE := fastlog.CalculateSuiteMSE(1+fastlog.K, -fastlog.K)
	if mse < bestMSE {
		fmt.Println()
		fmt.Println("New c0, c1 MSE beats the best known k !!")
	}
}
