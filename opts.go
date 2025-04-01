package mcit

import (
	"math/rand/v2"
	"time"
)

type searchOptions struct {
	src           rand.Source
	maxIters      int
	continuation  *Node
	expandShuffle bool
	done          bool
}

func newSearchOptions() *searchOptions {
	return &searchOptions{
		src:           rand.NewPCG(1337, 0xBEEF),
		expandShuffle: true,
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
func MaxIters(n int) Option { return Option{preFn: func(opts *searchOptions) { opts.maxIters = n }} }

// UseContinuation specifies a root node to continue a previous search from memory.
func UseContinuation(n *Node) Option {
	return Option{preFn: func(opts *searchOptions) { opts.continuation = n }}
}
func RandSource(src rand.Source) Option {
	return Option{preFn: func(opts *searchOptions) { opts.src = src }}
}
func ExpandShuffle(expandShuffle bool) Option {
	return Option{preFn: func(opts *searchOptions) { opts.expandShuffle = expandShuffle }}
}
