package mcit

import (
	"container/heap"
	"math/rand/v2"
)

// RunResults contains results returned from the user search function.
type RunResults struct {
	Expand   []string
	Priors   []float64
	Replace  bool
	Minimize bool
	Count    float64
	Value    float64
	// Payload contains optional user generated payload to store in the resulting tree Node.
	// This can be a reference to the user land state of the node which can be more convenient
	// when direct replay is difficult.
	Payload any
	Done    bool
}

// NodeSelector provides a choice of methods to select the current frontier node.
// Either a slice of actions, or a reference to the Payload returned from RunResults.
type NodeSelector struct {
	Actions []string
	Payload any
}

// Func is a search function containing user code which selects a frontier node and returns the results of experiments on it.
type Func func(selector NodeSelector) (results RunResults)

// Continuation is a structure which contains a root node to pass to continue a previous search from memory.
type Continuation struct {
	root *NodeStat
}

// Search is the main function from this package which implements Monte-carlo tree search.
//
// It accepts runFn containing user search code and calls it on each frontier node in accordance with the
// multi-armed bandit policy. Using a regret-optimal combination of exploration and exploitation.
//
// It takes options to configure aspects of the search.
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
		searchOpts.root = newRootStat()
		root = searchOpts.root
	}

	//	0bb. Initialize a frontier heap.
	frontier := []*NodeStat{}
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

		replay = curr.AppendLine(replay[:0])

		// 2. Run simulations at the frontier node.
		results := runFn(NodeSelector{replay, curr.Payload})

		//	2a. Copy minimize setting to current node.
		curr.Minimize = results.Minimize

		// 	2b. (optional) Shuffle expanded nodes before inserting them.
		if searchOpts.expandShuffle {
			s := results.Expand
			rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
		}

		// 	2c. (optional) Expand the node, and add children to the state.
		for i, e := range results.Expand {
			child, created := curr.NewChild(e)
			if !created {
				continue
			}

			//	2ca. (optional) Priors, if provided, should match the slice of expanded nodes.
			if len(results.Priors) > 0 {
				child.Prior = results.Priors[i]
			}

			//	2cb. Push the child onto the frontier heap.
			heap.Push((*byUCB1)(&frontier), child)
		}

		// 	2d. (optional) Exhaust the frontier node (or keep it around for next time).
		//                 The exhaust logic by default requires that some other conditions hold
		//                 To avoid the simulation running out of frontier values.
		if !results.Replace && (searchOpts.exhaustable || len(results.Expand) > 0) {
			heap.Remove((*byUCB1)(&frontier), curr.frontierIdx)
		}

		// 	2e. Backpropagate the results up the tree and fix the frontier nodes.
		for head := curr; head != nil; head = head.Parent {
			head.Runs += results.Count
			head.Value += results.Value
			head.RecomputePriority()
			if head.frontierIdx != -1 {
				heap.Fix((*byUCB1)(&frontier), head.frontierIdx)
			}
		}

		if searchOpts.done { //	3a. (optional) Stop search if done.
			return
		}

		// 3b. Restart from step 1.
	}
}
