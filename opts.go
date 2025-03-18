package mcit

import (
	"math/rand/v2"
	"time"
)

type searchOptions struct {
	src           rand.Source
	maxIters      int64
	expandShuffle bool
	continuation  *Continuation
	exhaustable   bool
	root          *node
	done          <-chan struct{}
	searchStats   *SearchStats
}

func newSearchOptions() *searchOptions {
	return &searchOptions{
		src:           rand.NewPCG(1337, 0xBEEF),
		expandShuffle: true,
		searchStats:   newSearchStats(),
	}
}

type Option struct {
	preFn  func(opts *searchOptions)
	postFn func(opts *searchOptions)
}

func Done(done <-chan struct{}) Option {
	return Option{preFn: func(opts *searchOptions) { opts.done = done }}
}
func DoneAfter(d time.Duration) Option {
	return Option{preFn: func(opts *searchOptions) {
		done := make(chan struct{})
		go func() {
			<-time.After(d)
			done <- struct{}{}
		}()
		opts.done = done
	}}
}
func Exhaustable(exhaustable bool) Option {
	return Option{preFn: func(opts *searchOptions) { opts.exhaustable = exhaustable }}
}
func MaxIters(n int64) Option { return Option{preFn: func(opts *searchOptions) { opts.maxIters = n }} }
func UseContinuation(c *Continuation) Option {
	return Option{preFn: func(opts *searchOptions) { opts.continuation = c; opts.root = c.root }}
}
func RandSource(src rand.Source) Option {
	return Option{preFn: func(opts *searchOptions) { opts.src = src }}
}
func ExpandShuffle(expandShuffle bool) Option {
	return Option{preFn: func(opts *searchOptions) { opts.expandShuffle = expandShuffle }}
}
func DetailedSearchStats(stats *SearchStats) Option {
	return Option{preFn: func(opts *searchOptions) { opts.searchStats = stats }, postFn: func(opts *searchOptions) {
		visitNodes(opts.root, func(n *node) bool {
			opts.searchStats.NodeCount++
			if len(n.children) == 0 {
				opts.searchStats.LeafCount++
			}
			if opts.root.frontierIdx == -1 {
				opts.searchStats.ExhaustedNodes++
			}
			return true
		})
	}}
}
func BestVariation(stat *NodeStat) Option {
	return Option{postFn: func(opts *searchOptions) {
		if stat == nil {
			return
		}
		n := getSelectLine(opts.root, selectChildFunc(maxNode))
		*stat = *newVariationStat(n)
	}}
}

func WorstVariation(stat *NodeStat) Option {
	return Option{postFn: func(opts *searchOptions) {
		if stat == nil {
			return
		}
		n := getSelectLine(opts.root, selectChildFunc(minNode))
		*stat = *newVariationStat(n)
	}}
}
func MostPopularVariation(stat *NodeStat) Option {
	return Option{postFn: func(opts *searchOptions) {
		if stat == nil {
			return
		}
		n := getSelectLine(opts.root, selectChildFunc(mostPopularNode))
		*stat = *newVariationStat(n)
	}}
}
func MinVariation(stat *NodeStat) Option {
	return Option{postFn: func(opts *searchOptions) {
		if stat == nil {
			return
		}
		n := chooseNode(nil, opts.root, minNode)
		*stat = *newVariationStat(n)
	}}
}
func MaxVariation(stat *NodeStat) Option {
	return Option{postFn: func(opts *searchOptions) {
		if stat == nil {
			return
		}
		n := chooseNode(nil, opts.root, maxNode)
		*stat = *newVariationStat(n)
	}}
}
func SearchTreeShallow(stat *NodeStat) Option {
	return Option{postFn: func(opts *searchOptions) {
		if stat == nil {
			return
		}
		*stat = *newShallowNodeStat(opts.root)
	}}
}
func SearchTree(stat *NodeStat) Option {
	return Option{postFn: func(opts *searchOptions) {
		if stat == nil {
			return
		}
		*stat = *newFullNodeStat(opts.root)
	}}
}
func Histogram(hist Hist, valueFn func(*NodeStat) float64) Option {
	return Option{postFn: func(opts *searchOptions) {
		for node := range nodeIter(opts.root) {
			n := newShallowNodeStat(node)
			x := valueFn(n)
			hist.Insert(x)
		}
	}}
}
func Visit(visitFn func(*NodeStat) bool) Option {
	return Option{postFn: func(opts *searchOptions) {
		visitNodes(opts.root, func(n *node) bool {
			s := newShallowNodeStat(n)
			return visitFn(s)
		})
	}}
}
func Count(results []int64, countFns ...func(*NodeStat) int64) Option {
	return Option{postFn: func(opts *searchOptions) {
		visitNodes(opts.root, func(n *node) bool {
			s := newShallowNodeStat(n)
			for i, f := range countFns {
				results[i] += f(s)
			}
			return true
		})
	}}
}
