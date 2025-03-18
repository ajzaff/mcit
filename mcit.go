package mcit

import (
	"container/heap"
	"math/rand/v2"
)

type RunResults struct {
	Expand  []string
	Priors  []float64
	Replace bool
	Count   float64
	Value   float64
	Done    bool
}

type Func func(actions []string) RunResults

type Continuation struct {
	root *node
}

func Search(runFn Func, opts ...Option) {
	// 0. Initialize state.
	searchOpts := newSearchOptions()
	// 0aa. Execute pre-run hooks and apply options.
	for _, o := range opts {
		if o.preFn != nil {
			o.preFn(searchOpts)
		}
	}

	//	0ba. Initialize root.
	root := searchOpts.root
	if root == nil {
		searchOpts.root = newRoot()
		root = searchOpts.root
	}

	//	0bb. Initialize a frontier heap.
	frontier := []*node{}
	heap.Push((*byUCB1)(&frontier), root)

	//	0c. Schedule post-run hooks.
	iters := int64(0)
	defer func() {
		searchOpts.searchStats.Iterations = iters
		searchOpts.searchStats.MaxFrontierSize = int64(len(frontier))
		// 4. Execute post-run hooks.
		for _, o := range opts {
			if o.postFn != nil {
				o.postFn(searchOpts)
			}
		}
	}()

	replay := make([]string, 0, 64)
	for ; (searchOpts.maxIters == 0 || iters < searchOpts.maxIters) && (!searchOpts.exhaustable || len(frontier) > 0); iters++ {
		// 1. Select a frontier node from the frontier heap and construct replay actions.
		curr := frontier[0]

		replay = curr.appendLine(replay[:0])

		// 2. Run simulations at the frontier node.
		results := runFn(replay)

		// 	2a. (optional) Shuffle expanded nodes before inserting them.
		if searchOpts.expandShuffle {
			s := results.Expand
			rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
		}

		// 	2b. (optional) Expand the node, and add children to the state.
		for i, e := range results.Expand {
			if _, found := curr.children[e]; found {
				continue
			}

			child := curr.newChild(e)

			//	2aa. (optional) Priors, if provided, should match the slice of expanded nodes.
			if len(results.Priors) > 0 {
				child.prior = results.Priors[i]
			}

			//	2ab. Push the child onto the frontier heap.
			heap.Push((*byUCB1)(&frontier), child)
		}

		// 	2b. (optional) Exhaust the frontier node (or keep it around for next time).
		//                 The exhaust logic by default requires that some other conditions hold
		//                 To avoid the simulation running out of frontier values.
		if !results.Replace && (searchOpts.exhaustable || len(results.Expand) > 0) {
			heap.Remove((*byUCB1)(&frontier), curr.frontierIdx)
		}

		// 	2c. Backpropagate the results up the tree and fix the frontier nodes.
		for head := curr; head != nil; head = head.parent {
			head.runs += results.Count
			head.value += results.Value
			head.recomputeUCB1()
			if head.frontierIdx != -1 {
				heap.Fix((*byUCB1)(&frontier), head.frontierIdx)
			}
		}

		select {
		case <-searchOpts.done: //	3a. (optional) Stop search if done.
			return
		default:
		}

		// 3b. Restart from step 1.
	}
}
