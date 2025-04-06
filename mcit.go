package mcit

import (
	"errors"
	"iter"
	"slices"
)

// ErrStop is returned from the search when the search is stopped prematurely with Stop.
var ErrStop = errors.New("stop")

// Context encapulates the control methods for the MCIT search.
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
	preserve bool
	minimize bool

	// Payload contains optional user generated payload to store in the resulting tree Node.
	// This can be a reference to the user land state of the node which can be more convenient
	// when direct replay is difficult.
	payload any

	count float32
	value float32

	done bool
	err  error
}

func (c *Context) Len() int { return len(c.actions) }

// Stop stops the search immediately with ErrStop.
func (c *Context) Stop() { c.StopErr(ErrStop) }

// StopErr stops the search immediately with the given error.
func (c *Context) StopErr(err error) { c.done = true; c.err = err }

// Payload returns the current payload.
func (c *Context) Payload() any { return c.payload }

// ReplacePayload replaces the payload on the current node.
func (c *Context) ReplacePayload(v any) { c.payload = v }

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

// Expand adds the given actions to the expand set from the current node.
func (c *Context) Expand(actions ...string) { c.expand = append(c.expand, actions...) }

// Priors optionally adds the given unnormalized priors to the expand set from the current node.
//
// If Priors are used, they must be called one-to-one with Expand.
func (c *Context) Priors(priors ...float32) { c.priors = append(c.priors, priors...) }

// Preserve may be called to keep the node from being marked "exhausted".
//
// This will allow further simulations to be continued from this node as well as new nodes to be expanded.
// The default behavior is to automatically exhaust nodes which have been called once.
func (c *Context) Preserve() { c.preserve = true }

// Minimize sets the objective function of the current node and its subtree.
//
// By default the objective function is "maximize". The objective function is
// inherited from parents to children, but Maximize may be used to set it
// explicitly, for instance, in two player game contexts.
func (c *Context) Minimize(minimize bool) { c.minimize = minimize }

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
