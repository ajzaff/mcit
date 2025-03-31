package mcit

import (
	"iter"
	"math"
	"math/rand/v2"
	"slices"
)

const exploreTerm = math.Pi

type Node struct {
	Parent  *Node
	Action  string
	Height  int
	Payload any
	// Frontier is set to the index of the first Bandit which has never been tried
	// Otherwise it equals len(Bandits).
	Frontier int
	// Bandits nodes which
	Bandits []Stat
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
	//       parent.Frontier keeps track of the first index of frontier nodes.
	parent.Bandits = append(parent.Bandits, stat)

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

func (s *Node) Hist(hist Hist, valueFn func(Stat) float64) {
	for _, e := range s.Bandits {
		x := valueFn(e)
		hist.Insert(x)
	}
	for _, child := range s.Children {
		child.Hist(hist, valueFn)
	}
}

type byUCB1 []Stat

func (a byUCB1) Len() int      { return len(a) }
func (a byUCB1) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byUCB1) Less(i, j int) bool {
	if ui, uj := a[i].Priority, a[j].Priority; ui != uj {
		// Higher priority nodes first.
		return ui > uj
	}
	// When priorities are equal (often +∞), fall back to prior comparison.
	return a[i].Prior > a[j].Prior
}
func (a *byUCB1) Push(x any) {
	e := x.(Stat)
	// Now is a good time to compute the priority.
	// Set the frontier index to the current position.
	e.RecomputePriority()
	*a = append(*a, e)
}
func (a *byUCB1) Pop() any {
	n := len(*a) - 1
	x := (*a)[n]
	*a = (*a)[:n]
	return x
}

func selectChildFunc(r *rand.Rand, cmpFn func(a, b Stat) int) func(*Node) *Node {
	return func(root *Node) *Node {
		// Create an equivalence slice for implementing fair random choice
		// To tie break between equivalent children according to cmpFn.
		equal := []Stat{{}}
		for _, b := range root.Bandits {
			a := equal[0]
			switch c := cmpFn(a, b); {
			case c > 0:
				// Swap out for a better node, and reset the equivalence slice.
				equal[0] = b
				equal = equal[:1]
			case c == 0:
				// When equal, add it to the equivalence slice.
				equal = append(equal, b)
			}
		}
		if len(equal) > 1 {
			// Fair random choice between equivalent children.
			i := r.IntN(len(equal))
			equal[0], equal[i] = equal[i], equal[0]
		}
		// When nil, no action was selected.
		return root.Children[equal[0].Action]
	}
}

func getSelectLine(root *Node, selectFn func(*Node) *Node) *Node {
	for len(root.Children) > 0 {
		next := selectFn(root)
		if next == nil {
			break
		}
		root = next
	}
	return root
}

func nodeIter(root *Node) iter.Seq[*Node] {
	return func(yield func(*Node) bool) { visitNodes(root, yield) }
}

func visitNodes(root *Node, visitFn func(*Node) bool) {
	if !visitFn(root) {
		return
	}
	for _, child := range root.Children {
		visitNodes(child, visitFn)
	}
}

func statIter(root *Node) iter.Seq[Stat] {
	return func(yield func(Stat) bool) {
		for n := range nodeIter(root) {
			for _, e := range n.Bandits {
				if !yield(e) {
					break
				}
			}
		}
	}
}
