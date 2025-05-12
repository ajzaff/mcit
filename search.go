package mcit

import (
	"math/rand/v2"
	"time"

	"github.com/ajzaff/lazyq"
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

	r := rand.New(searchOpts.src)
	maxItersDefined := searchOpts.maxIters > 0
	exploreFactor := searchOpts.exploreFactor

	c := &Context{
		actions: make([]string, 0, 64),
		expand:  make([]string, 0, 64),
		priors:  make([]float32, 0, 64),
	}

	for {
		// 1. Select a frontier node with the maximum bandit at each step. Construct replay actions.
		frontier := root
		c.reset()
		for frontier.Exhausted() && frontier.Queue.Len() > 0 {
			next := frontier.next()
			if next.Node == nil {
				break // Search from here.
			}
			c.actions = append(c.actions, next.Action)
			frontier = next.Node
		}

		// 2. Run simulations at the frontier node.
		runFn(c)
		frontier.Flags &= ^FlagsMinimize
		frontier.Flags |= c.flags & FlagsMinimize

		// 	2b. (optional) Shuffle expanded nodes before inserting them.
		if searchOpts.expandShuffle && len(c.expand) > 1 {
			r.Shuffle(len(c.expand), func(i, j int) { c.expand[i], c.expand[j] = c.expand[j], c.expand[i] })
		}

		// 	2c. (optional) Expand the node, and add children to the state.
		for i, action := range c.expand {
			//	2ca. (optional) Priors, if provided, should match the slice of expanded nodes.
			// FIXME: Implement prior normalization and renormalizaion.
			prior := float32(exploreFactor)
			if len(c.priors) > 0 {
				prior *= c.priors[i]
			}
			frontier.NewChild(action, prior)
		}

		// 	2d. (optional) Keep the frontier node in the frontier set.
		if c.flags.Exhausted() {
			frontier.Flags |= FlagsExhausted
		}

		// 	2e. Backpropagate the results up the tree and fix the bandit heaps along the way.
		for head := frontier.Parent; head != nil; head = head.Parent {
			head.addValueRuns(c.value, c.count)
			// Recompute the PUCT policy value for the frontier.
			bandit := lazyq.First(head.Queue)
			head.Queue.Decrease(bandit.computePriority(head.logTrials()))
		}

		// 	3. State keeping and termination.
		iters++ //	3a. Increment iterations.

		if searchOpts.done || c.done { //	3b. (optional) Stop search if done.
			result.Err = c.err
			break
		}

		//	3c. End the search when maxIters is reached.
		if maxItersDefined && iters >= searchOpts.maxIters {
			break
		}

		// 4. Restart from step 1.
	}

	// 5. Store search results.
	return Result{
		Root:       root,
		Iterations: iters,
		Duration:   time.Since(start),
	}
}
