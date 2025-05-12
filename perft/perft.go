// Package perft provide statistical and analytical performance tests.
package perft

import (
	"github.com/ajzaff/lazyq"
	"github.com/ajzaff/mcit"
)

type SearchStats struct {
	NodeCount      int64
	LeafCount      int64
	ExhaustedCount int64
	Height         int64
	DeepestRun     int64
}

func DetailedSearchStats(root *mcit.Node) SearchStats {
	var results SearchStats
	visitNodes(root, 0, func(n *mcit.Node, depth int) bool {
		// NodeCount
		results.NodeCount++
		// LeafCount
		if n.Queue.Len() == 0 {
			results.LeafCount++
		}
		// ExhaustedCount
		if n.Exhausted() {
			results.ExhaustedCount++
		}
		// MaxHeight
		if results.Height < int64(depth) {
			results.Height = int64(depth)
		}
		for stat := range lazyq.Payloads(n.Queue) {
			// MaxHeightRun
			if stat.Runs > 0 && results.DeepestRun < int64(depth) {
				results.DeepestRun = int64(depth)
			}
		}
		return true
	})
	return results
}
