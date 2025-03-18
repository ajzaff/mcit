package mcit

import "math/rand/v2"

type searchOptions struct {
	runFn         Func
	src           *rand.Source
	expandShuffle bool
	root          *node
	searchStats   *SearchStats
}

func newSearchOptions() *searchOptions {
	return &searchOptions{
		runFn:         nil,
		src:           nil,
		expandShuffle: true,
	}
}

type Option struct {
	preFn  func(opts *searchOptions)
	postFn func(opts *searchOptions)
}

func RunFunc(runFn Func) Option {
	return Option{preFn: func(so *searchOptions) { so.runFn = runFn }}
}
func RandSource(src *rand.Source) Option {
	return Option{preFn: func(opts *searchOptions) { opts.src = src }}
}
func ExpandShuffle(expandShuffle bool) Option {
	return Option{preFn: func(opts *searchOptions) { opts.expandShuffle = expandShuffle }}
}
func CollectSearchStats(stats *SearchStats) Option {
	return Option{preFn: func(opts *searchOptions) { opts.searchStats = stats }}
}
func MostPopularVariation(stat *NodeStat) Option {
	return Option{postFn: func(opts *searchOptions) {
		if stat == nil {
			return
		}
		n := getSelectLine(opts.root, getMostRunChild)
		*stat = *newVariationStat(n)
	}}
}
func MinVariation(stat *NodeStat) Option {
	return Option{postFn: func(opts *searchOptions) {
		if stat == nil {
			return
		}
		n := getSelectLine(opts.root, getMinChild)
		*stat = *newVariationStat(n)
	}}
}
func MaxVariation(stat *NodeStat) Option {
	return Option{postFn: func(opts *searchOptions) {
		if stat == nil {
			return
		}
		n := getSelectLine(opts.root, getMaxChild)
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
		*stat = *newNodeStat(opts.root)
	}}
}
