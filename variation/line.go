package variation

import (
	"slices"

	"github.com/ajzaff/mcts"
)

// Line computes the actions leading up to n from the root.
//
// Line is equivalent to AppendLine(n, nil).
func Line(n *mcts.Node) []string { return AppendLine(n, nil) }

// AppendLine appends the line leading up to n from the root to buf and returns the modified slice.
func AppendLine(n *mcts.Node, buf []string) []string {
	i := len(buf)
	buf = slices.Grow(buf[i:], 1+Depth(n))
	for ; n.Parent != nil; n = n.Parent {
		buf = append(buf, n.Action)
	}
	slices.Reverse(buf[i:])
	return buf
}
