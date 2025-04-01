// Package perft provide statistical and analytical performance tests.
package perft

import (
	"github.com/ajzaff/mcit"
)

type SearchStats struct {
	NodeCount      int64
	LeafCount      int64
	ExhaustedCount int64
	MaxHeight      int64
	MaxHeightRun   int64
}

func DetailedSearchStats(root *mcit.Node) SearchStats {
	var results SearchStats
	visitNodes(root, func(n *mcit.Node) bool {
		// NodeCount
		results.NodeCount++
		// LeafCount
		if len(n.Children) == 0 {
			results.LeafCount++
		}
		// ExhaustedCount
		if n.Exhausted {
			results.ExhaustedCount++
		}
		// MaxHeight
		if results.MaxHeight < int64(n.Height) {
			results.MaxHeight = int64(n.Height)
		}
		for stat := range n.StatSeq() {
			// MaxHeightRun
			if stat.Runs > 0 && results.MaxHeightRun < int64(n.Height) {
				results.MaxHeightRun = int64(n.Height)
			}
		}
		return true
	})
	return results
}
