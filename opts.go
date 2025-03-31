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
	root          *Node
	done          bool
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
	return Option{preFn: func(opts *searchOptions) {
		go func() {
			<-done
			opts.done = true
		}()
	}}
}
func DoneAfter(d time.Duration) Option {
	return Option{preFn: func(opts *searchOptions) {
		go func() {
			<-time.After(d)
			opts.done = true
		}()
	}}
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
		visitNodes(opts.root, func(n *Node) bool {
			opts.searchStats.NodeCount++
			if len(n.Children) == 0 {
				opts.searchStats.LeafCount++
			}
			if n.Exhausted {
				opts.searchStats.ExhaustedNodes++
			}
			return true
		})
	}}
}
func MaxVariation(r *rand.Rand, stat *Node) Option {
	return Option{postFn: func(opts *searchOptions) {
		if stat == nil {
			return
		}
		*stat = *getSelectLine(opts.root, selectChildFunc(r, compareMaxStat))
	}}
}

func MinVariation(r *rand.Rand, stat *Node) Option {
	return Option{postFn: func(opts *searchOptions) {
		if stat == nil {
			return
		}
		*stat = *getSelectLine(opts.root, selectChildFunc(r, compareMinStat))
	}}
}
func MostPopularVariation(r *rand.Rand, stat *Node) Option {
	return Option{postFn: func(opts *searchOptions) {
		if stat == nil {
			return
		}
		*stat = *getSelectLine(opts.root, selectChildFunc(r, compareStatPopularity))
	}}
}
func SearchTreeShallow(stat *Node) Option {
	return Option{postFn: func(opts *searchOptions) {
		if stat == nil {
			return
		}
		*stat = *opts.root.Detatched()
	}}
}
func SearchTree(stat *Node) Option {
	return Option{postFn: func(opts *searchOptions) {
		if stat == nil {
			return
		}
		*stat = *opts.root
	}}
}
func Histogram(hist Hist, valueFn func(Stat) float64) Option {
	return Option{postFn: func(opts *searchOptions) {
		for e := range statIter(opts.root) {
			x := valueFn(e)
			hist.Insert(x)
		}
	}}
}
func Visit(visitFn func(*Node) bool) Option {
	return Option{postFn: func(opts *searchOptions) {
		visitNodes(opts.root, func(n *Node) bool { return visitFn(n) })
	}}
}
func Count(results []int64, countFns ...func(*Node) int64) Option {
	return Option{postFn: func(opts *searchOptions) {
		visitNodes(opts.root, func(n *Node) bool {
			for i, f := range countFns {
				results[i] += f(n)
			}
			return true
		})
	}}
}
