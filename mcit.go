package mcit

import "math/rand/v2"

// RunResults contains results returned from the user search function.
type RunResults struct {
	// Expand is the slice of actions available from this node which should be available as children.
	// Expand may be a subset of all possible actions. See Replace for tips
	// implementing partial expansion.
	Expand []string
	// Priors is a slice of prior values to apply to new expanded nodes.
	// Priors may be empty in which case the default prior value is used.
	Priors []float64
	// Replace can be set to true when the node should be returned to the frontier queue.
	// This is useful when allowing nodes to be partially expanded on each new visit OR
	// when the node is a leaf node of the current search and we want to repeatedly explore it.
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
	root *Node
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
		searchOpts.root = newRoot()
		root = searchOpts.root
	}

	//	0c. Schedule post-run hooks.
	iters := int64(0)
	defer func() {
		searchOpts.searchStats.Iterations = iters
		// searchOpts.searchStats.MaxFrontierSize = int64(len(frontier))
		// 4. Execute post-run hooks.
		for _, o := range opts {
			if o.postFn != nil {
				o.postFn(searchOpts)
			}
		}
	}()

	r := rand.New(searchOpts.src)
	replay := make([]string, 0, 64)
	maxItersDefined := searchOpts.maxIters > 0

	for {
		// 1. Select a frontier node with the maximum bandit at each step. Construct replay actions.
		frontier := root
		replay = replay[:0]
		for frontier.Exhausted && len(frontier.Bandits) > 0 {
			action := frontier.next().Action
			next := frontier.Children[action]
			if next == nil {
				break
			}
			replay = append(replay, action)
			frontier = next
		}

		// 2. Run simulations at the frontier node.
		results := runFn(NodeSelector{replay, frontier.Payload})

		// 	2b. (optional) Shuffle expanded nodes before inserting them.
		if searchOpts.expandShuffle {
			r.Shuffle(len(results.Expand), func(i, j int) { results.Expand[i], results.Expand[j] = results.Expand[j], results.Expand[i] })
		}

		// 	2c. (optional) Expand the node, and add children to the state.
		for i, action := range results.Expand {
			//	2ca. (optional) Priors, if provided, should match the slice of expanded nodes.
			// FIXME: Implement prior normalization and renormalizaion.
			prior := 1.
			if len(results.Priors) > 0 {
				prior = results.Priors[i]
			}
			frontier.NewChild(action, prior)
		}

		// 	2d. (optional) Keep the frontier node in the pool.
		if !results.Replace {
			frontier.Exhausted = true
		}

		// 	2e. Backpropagate the results up the tree and fix the bandit heaps along the way.
		for head := frontier.Parent; head != nil; head = head.Parent {
			head.Bandits[0].Runs += results.Count
			head.Bandits[0].Value += results.Value
			head.Bandits[0].RecomputePriority()
			head.down(0)
		}

		// 	3. State keeping and termination.
		iters++ //	3a. Increment iterations.

		if searchOpts.done || results.Done { //	3b. (optional) Stop search if done.
			return
		}

		//	3c. End the search when maxIters is reached.
		if maxItersDefined && iters >= searchOpts.maxIters {
			return
		}

		// 4. Restart from step 1.
	}
}
