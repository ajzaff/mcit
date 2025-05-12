package mcts

import (
	"github.com/ajzaff/lazyq"
)

type Flags int32

const (
	FlagsMinimize = 1 << iota

	// FlagsExhausted marks whether we are done expanding for this node.
	// 	* When set, we will not simulate this node further and will rely on the Bandit policy.
	// 	* When unset, we will generate more simulations (and possibly children) in the future.
	FlagsExhausted
)

func (f Flags) Minimize() bool  { return f&FlagsMinimize != 0 }
func (f Flags) Exhausted() bool { return f&FlagsExhausted != 0 }

type Child struct {
	Action string
	Stat
	*Node
}

type Node struct {
	Parent *Node
	Action string
	Trials float32
	Flags
	Queue lazyq.Queue[Child]
}

func newRoot() *Node { return &Node{} }

func (n *Node) indexChild(action string) (int, bool) {
	m := n.Queue.Len()
	for i := range m {
		if lazyq.At(n.Queue, i).Action == action {
			return i, true
		}
	}
	return -1, false
}

// NewChild creates a new child on the parent Node.
// Pushes a node stat to the list of bandits.
func (parent *Node) NewChild(action string, exploreFactor float32) (created bool) {
	if _, found := parent.indexChild(action); found {
		return false
	}
	// NOTE: We defer child creation until node is actually opened.
	//       This saves allocations for nodes that are never explored.
	// NOTE: We don't use heapify here. The majority of actions are never tried so we don't waste time with the O(log N) heap.Push operation.
	//       lazyq keeps track of the first index of frontier nodes.
	parent.Queue.AppendMax(Child{Action: action, Stat: Stat{ExploreFactor: exploreFactor}})
	return true
}

func (s *Node) next() Child {
	if lazyq.HasMaxElems(s.Queue) {
		// We have at least one node which has never been tried before.
		// Use this time to fix the position in the heap so we can select it.
		// Nodes which have never been tried before always take priority.
		//
		// Waiting until now to fix this position is largely an optimization
		// as we don't expect the majority of nodes of large trees to be tried
		// we don't need to waste time with the O(log N) heap.Push operation.
		//
		// Create the new child now.
		// By defering child creation until the last minute
		// we save tons of allocations for nodes which are never explored.
		stat := lazyq.FirstMaxElem(s.Queue)
		child := &Node{
			Parent: s,
			Action: stat.Action,
			// Copy the setting from the parent. Run will have a chance to override this.
			// See Context.Minimize.
			Flags: s.Flags & FlagsMinimize,
		}
		stat.Node = child
		lazyq.ReplacePayload(s.Queue, lazyq.MaxIndex(s.Queue), stat)
	}
	// NOTE: We always take the first action.
	// If we ever implemented a temperature feature, we'd need to keep track of this index.
	return s.Queue.Next()
}

// Stat attempts to locate the Child entry in the parent node.
// FIXME: Currently, this is not an efficient operation.
func (s *Node) Stat() *Child {
	if s == nil || s.Parent == nil {
		return nil
	}
	for e := range lazyq.Payloads(s.Parent.Queue) {
		if e.Action == s.Action {
			return &e
		}
	}
	return nil
}
