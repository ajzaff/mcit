package mcts

import (
	"errors"
	"iter"
	"slices"
)

// ErrStop is returned from the search when the search is stopped prematurely with Stop.
var ErrStop = errors.New("stop")

// Context encapulates the control methods for the MCTS search.
type Context struct {
	actions []string

	// expand is the slice of actions available from this node which should be available as children.
	// expand may be a subset of all possible actions. See Replace for tips
	// implementing partial expansion.
	expand []string
	// priors is a slice of prior values to apply to new expanded nodes.
	// priors may be empty in which case the default prior value is used.
	priors []float32
	// preserve can be set to true when the node should be returned to the frontier queue.
	// This is useful when allowing nodes to be partially expanded on each new visit OR
	// when the node is a leaf node of the current search and we want to repeatedly explore it.
	flags Flags

	count float32
	value float32

	done bool
	err  error
}

func (c *Context) reset() {
	c.actions = c.actions[:0]
	c.expand = c.expand[:0]
	c.priors = c.priors[:0]
}

func (c *Context) Len() int { return len(c.actions) }

// Stop stops the search immediately with ErrStop.
func (c *Context) Stop() { c.StopErr(ErrStop) }

// StopErr stops the search immediately with the given error.
func (c *Context) StopErr(err error) { c.done = true; c.err = err }

// Actions returns an iterator of actions from root up to the current node.
//
// provides a choice of methods to select the current frontier node.
// Either a slice of actions, or a reference to the Payload returned from RunResults.
func (c *Context) Actions() iter.Seq[string] { return slices.Values(c.actions) }

// Actions returns an iterator of actions from root up to the current node.
func (c *Context) Actions2() iter.Seq2[int, string] { return slices.All(c.actions) }

// ActionAt returns the action at the given index or the invalid action, empty string.
func (c *Context) ActionAt(i int) string {
	if i < 0 || len(c.actions) <= i {
		return ""
	}
	return c.actions[i]
}

// Append is like [Expand] but does not exhaust the node.
func (c *Context) Append(actions ...string) { c.expand = append(c.expand, actions...) }

func (c *Context) exhaust() { c.flags |= FlagsExhausted }

// Expand adds the given actions to the expand set from the current node.
//
// Exhausts the node. An exhausted node cannot expand again in the future,
// but still receives priority updates via backpropoagaion.
func (c *Context) Expand(actions ...string) { c.Append(actions...); c.exhaust() }

// Priors optionally adds the given unnormalized priors to the expand set from the current node.
//
// If Priors are used, they must be called one-to-one with Expand.
func (c *Context) Priors(priors ...float32) { c.priors = append(c.priors, priors...) }

// Minimize sets the objective function of the current node and its subtree to minimize.
//
// The default is maximize.
func (c *Context) Minimize() { c.flags |= FlagsMinimize }

// Maximize sets the objective function of the current node and its subtree to maximize.
//
// This is the default.
func (c *Context) Maximize() { c.flags &= ^FlagsMinimize }

// SetResult sets the result of the experiment to the explicit value and number of experiments.
func (c *Context) SetResult(value, count float32) { c.value = value; c.count = count }

// SetResultValue sets the result of the experiment to the explicit value and one experiment run.
func (c *Context) SetResultValue(value float32) { c.value = value; c.count = 1 }

// AddResultValue adds the result of the experiment and increments the number of runs.
func (c *Context) AddResultValue(value float32) { c.AddResult(value, 1) }

// AddResultValue adds the result of the experiment and increments the number of runs.
func (c *Context) AddResult(value, count float32) { c.value += value; c.count += count }

// AddValue adds the value to the experiment results.
func (c *Context) AddValue(value float32) { c.value += value }

// AddCount adds the count to the number of experiment runs.
func (c *Context) AddCount(count float32) { c.count += count }

// Func is a search function containing user code which selects a frontier node and returns the results of experiments on it.
//
// The search stops if an error is returned.
type Func func(c *Context)
