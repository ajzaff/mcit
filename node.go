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
func (parent *Node) NewChild(action string, prior float64) (child *Node, created bool) {
	if child, found := parent.Children[action]; found {
		return child, false
	}
	child = &Node{
		Parent: parent,
		Height: parent.Height + 1,
		Action: action,
	}
	if parent.Children == nil {
		parent.Children = map[string]*Node{}
	}
	parent.Children[action] = child

	stat := Stat{Action: action, Prior: prior, Priority: math.Inf(+1)}
	// NOTE: We don't use heap.Push here. The majority of actions are never tried so we don't waste time with the O(log N) heap.Push operation.
	//       LazyHeap keeps track of the first index of frontier nodes.
	parent.append(stat)

	return child, true
}

// Detatched returns a shallow clone of the stat object detatched from patents, children, and the frontier
// without modifying the original stat object. The tree Height is not reset.
func (s *Node) Detatched() *Node {
	var copy Node
	copy = *s
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
