package mcit

import (
	"math/rand/v2"
	"time"
)

type searchOptions struct {
	src           rand.Source
	maxIters      int64
	continuation  *Continuation
	expandShuffle bool
	exhaustable   bool
	root          *NodeStat
	done          <-chan struct{}
	searchStats   *SearchStats
}

func newSearchOptions() *searchOptions {
	return &searchOptions{
		src:           rand.NewPCG(1337, 0xBEEF),
		expandShuffle: true,
		exhaustable:   true,
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
		visitNodes(opts.root, func(n *NodeStat) bool {
			opts.searchStats.NodeCount++
			if len(n.Children) == 0 {
				opts.searchStats.LeafCount++
			}
			if n.Exhausted() {
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
		*stat = *getSelectLine(opts.root, selectChildFunc(maxNodeStrictlyPositive))
	}}
}

func WorstVariation(stat *NodeStat) Option {
	return Option{postFn: func(opts *searchOptions) {
		if stat == nil {
			return
		}
		*stat = *getSelectLine(opts.root, selectChildFunc(minNodeStrictlyPositive))
	}}
}
func MostPopularVariation(stat *NodeStat) Option {
	return Option{postFn: func(opts *searchOptions) {
		if stat == nil {
			return
		}
		*stat = *getSelectLine(opts.root, selectChildFunc(mostPopularNode))
	}}
}
func MinVariation(stat *NodeStat) Option {
	return Option{postFn: func(opts *searchOptions) {
		if stat == nil {
			return
		}
		*stat = *chooseNode(nil, opts.root, minNode)
	}}
}
func MaxVariation(stat *NodeStat) Option {
	return Option{postFn: func(opts *searchOptions) {
		if stat == nil {
			return
		}
		*stat = *chooseNode(nil, opts.root, maxNode)
	}}
}
func SearchTreeShallow(stat *NodeStat) Option {
	return Option{postFn: func(opts *searchOptions) {
		if stat == nil {
			return
		}
		*stat = *opts.root.Detatched()
	}}
}
func SearchTree(stat *NodeStat) Option {
	return Option{postFn: func(opts *searchOptions) {
		if stat == nil {
			return
		}
		*stat = *opts.root
	}}
}
func Histogram(hist Hist, valueFn func(*NodeStat) float64) Option {
	return Option{postFn: func(opts *searchOptions) {
		for n := range nodeIter(opts.root) {
			x := valueFn(n)
			hist.Insert(x)
		}
	}}
}
func Visit(visitFn func(*NodeStat) bool) Option {
	return Option{postFn: func(opts *searchOptions) {
		visitNodes(opts.root, func(n *NodeStat) bool { return visitFn(n) })
	}}
}
func Count(results []int64, countFns ...func(*NodeStat) int64) Option {
	return Option{postFn: func(opts *searchOptions) {
		visitNodes(opts.root, func(n *NodeStat) bool {
			for i, f := range countFns {
				results[i] += f(n)
			}
			return true
		})
	}}
}
