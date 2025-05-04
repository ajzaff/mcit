package mcit

import (
	"math/rand/v2"
	"time"
)

// Result of a search containing the root search node and total number of iterations of MCTS performed.
type Result struct {
	Root       *Node
	Iterations int
	Duration   time.Duration
	Err        error
}

// Search is the main function from this package which implements Monte-carlo tree search.
//
// It accepts runFn containing user search code and calls it on each frontier node in accordance with the
// multi-armed bandit policy. Using a regret-optimal combination of exploration and exploitation.
//
// It takes options to configure aspects of the search.
func Search(runFn Func, opts ...Option) (result Result) {
	// 0. Initialize state.
	searchOpts := newSearchOptions()
	// 0aa. Execute pre-run hooks and apply options.
	for _, optFn := range opts {
		optFn(searchOpts)
	}

	//	0ba. Initialize root.
	var root *Node
	if root = searchOpts.continuation; root == nil {
		root = newRoot()
	}

	var iters int
	start := time.Now()

	defer func() {
		// 4. Store search results.
		result = Result{
			Root:       root,
			Iterations: iters,
			Duration:   time.Since(start),
		}
	}()

	r := rand.New(searchOpts.src)
	replay := make([]string, 0, 64)
	expand := make([]string, 0, 64)
	priors := make([]float32, 0, 64)
	maxItersDefined := searchOpts.maxIters > 0
	exploreFactor := searchOpts.exploreFactor

	for {
		// 1. Select a frontier node with the maximum bandit at each step. Construct replay actions.
		frontier := root
		replay = replay[:0]
		for frontier.Exhausted && frontier.Queue.Len() > 0 {
			action := frontier.next().Action
			next := frontier.Children[action]
			if next == nil {
				break
			}
			replay = append(replay, action)
			frontier = next
		}

		// 2. Run simulations at the frontier node.
		c := Context{actions: replay, expand: expand[:0], priors: priors[:0]}
		runFn(&c)
		frontier.Minimize = c.minimize
		expand = c.expand
		priors = c.priors

		// 	2b. (optional) Shuffle expanded nodes before inserting them.
		if searchOpts.expandShuffle && len(expand) > 1 {
			r.Shuffle(len(expand), func(i, j int) { expand[i], expand[j] = expand[j], expand[i] })
		}

		// 	2c. (optional) Expand the node, and add children to the state.
		for i, action := range expand {
			//	2ca. (optional) Priors, if provided, should match the slice of expanded nodes.
			// FIXME: Implement prior normalization and renormalizaion.
			prior := float32(1)
			if len(c.priors) > 0 {
				prior = c.priors[i]
			}
			frontier.NewChild(action, prior)
		}

		// 	2d. (optional) Keep the frontier node in the frontier set.
		if !c.preserve {
			frontier.Exhausted = true
		}

		// 	2e. Backpropagate the results up the tree and fix the bandit heaps along the way.
		for head := frontier.Parent; head != nil; head = head.Parent {
			head.addValueRuns(c.value, c.count)
			head.recomputePriority(exploreFactor)
		}

		// 	3. State keeping and termination.
		iters++ //	3a. Increment iterations.

		if searchOpts.done || c.done { //	3b. (optional) Stop search if done.
			result.Err = c.err
			return
		}

		//	3c. End the search when maxIters is reached.
		if maxItersDefined && iters >= searchOpts.maxIters {
			return
		}

		// 4. Restart from step 1.
	}
}
