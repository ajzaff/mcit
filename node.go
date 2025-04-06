package mcit

import (
	"math"
	"slices"
)

const exploreTerm = math.Pi

type Node struct {
	Parent   *Node
	Action   string
	Height   int
	Payload  any
	Minimize bool
	lazyQueue
	// Exhausted marks whether we are done with this node.
	// 	* When true, we will not simulate this node further and will rely on the Bandit policy.
	// 	* When false, we will generate more simulations (and possibly children) in the future.
	Exhausted bool
	Children  map[string]*Node
}

func newRoot() *Node { return &Node{} }

// NewChild creates a new child on the parent Node.
// Pushes a node stat to the list of bandits.
func (parent *Node) NewChild(action string, prior float32) (created bool) {
	if _, found := parent.Children[action]; found {
		return false
	}
	if parent.Children == nil {
		parent.Children = map[string]*Node{}
	}

	// NOTE: We defer child creation until node is actually opened.
	//       This saves allocations for nodes that are never explored.
	parent.Children[action] = nil
	stat := Stat{Action: action, Prior: prior, Priority: float32(math.Inf(+1))}
	// NOTE: We don't use heap.Push here. The majority of actions are never tried so we don't waste time with the O(log N) heap.Push operation.
	//       LazyHeap keeps track of the first index of frontier nodes.
	parent.append(stat)
	return true
}

func (s *Node) next() Stat {
	if s.hasLazyElements() {
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
		stat := s.Bandits[s.lazyIndex]
		child := &Node{
			Parent: s,
			Height: s.Height + 1,
			Action: stat.Action,
		}
		if s.Children == nil {
			s.Children = map[string]*Node{}
		}
		s.Children[stat.Action] = child
		s.upLazy()
	}
	// NOTE: We always take the first action.
	// If we ever implemented a temperature feature, we'd need to keep track of this index.
	return s.top()
}

// Detatched returns a shallow clone of the stat object detatched from patents, children, and the frontier
// without modifying the original stat object. The tree Height is not reset.
func (s *Node) Detatched() *Node {
	copy := *s
	copy.Parent = nil
	copy.Children = nil
	return &copy
}

// Stat attempts to locate the stat entry in the parent node.
// FIXME: Currently, this is not an efficient operation.
func (s *Node) Stat() *Stat {
	if s.Parent == nil {
		return nil
	}
	for _, e := range s.Parent.Bandits {
		if e.Action == s.Action {
			return &e
		}
	}
	return nil
}

func (s *Node) Line() []string { return s.AppendLine(nil) }

func (s *Node) AppendLine(buf []string) []string {
	i := len(buf)
	buf = slices.Grow(buf[i:], 1+s.Height)
	for ; s.Parent != nil; s = s.Parent {
		buf = append(buf, s.Action)
	}
	slices.Reverse(buf[i:])
	return buf
}
