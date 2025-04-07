package mcit

import (
	"math"
	"math/rand/v2"
	"time"
)

type searchOptions struct {
	src           rand.Source
	maxIters      int
	continuation  *Node
	expandShuffle bool
	exploreFactor float32
	done          bool
}

func newSearchOptions() *searchOptions {
	return &searchOptions{
		src:           rand.NewPCG(1337, 0xBEEF),
		expandShuffle: true,
		exploreFactor: 2 * math.Pi,
	}
}

type Option func(opts *searchOptions)

func Done(done <-chan struct{}) Option {
	return Option(func(opts *searchOptions) {
		go func() {
			<-done
			opts.done = true
		}()
	})
}
func DoneAfter(d time.Duration) Option {
	return Option(func(opts *searchOptions) {
		go func() {
			<-time.After(d)
			opts.done = true
		}()
	})
}
func MaxIters(n int) Option { return Option(func(opts *searchOptions) { opts.maxIters = n }) }

// UseContinuation specifies a root node to continue a previous search from memory.
func UseContinuation(n *Node) Option {
	return Option(func(opts *searchOptions) { opts.continuation = n })
}
func RandSource(src rand.Source) Option {
	return Option(func(opts *searchOptions) { opts.src = src })
}
func ExpandShuffle(expandShuffle bool) Option {
	return Option(func(opts *searchOptions) { opts.expandShuffle = expandShuffle })
}
func ExploreFactor(exploreFactor float32) Option {
	return Option(func(opts *searchOptions) { opts.exploreFactor = exploreFactor })
}
