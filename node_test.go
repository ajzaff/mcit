package mcit

import "testing"

func TestJustInTimeNodeAllocation(t *testing.T) {
	var n Node

	const action = "foo"

	n.NewChild(action, 1)

	child, ok := n.Children[action]
	if !ok {
		t.Errorf("TestJustInTimeNodeAllocation(): expected child to exist in Node but it was not found")
	}
	if child != nil {
		t.Errorf("TestJustInTimeNodeAllocation(): expected child to be nil but it was non-nil")
	}

	n.next()

	if child := n.Children[action]; child == nil {
		t.Errorf("TestJustInTimeNodeAllocation(): expected child to be allocated after next() but it was nil")
	}
}
