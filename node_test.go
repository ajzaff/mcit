package mcit

import (
	"testing"

	"github.com/ajzaff/lazyq"
)

func TestJustInTimeNodeAllocation(t *testing.T) {
	var n Node

	const action = "foo"

	n.NewChild(action, 1)

	if n.Queue.Len() < 1 {
		t.Fatalf("TestJustInTimeNodeAllocation(): expected child to exist in Node but it was not found")
	}
	s := lazyq.FirstMaxElem(n.Queue)
	if s.Node != nil {
		t.Errorf("TestJustInTimeNodeAllocation(): expected child to be nil but it was non-nil")
	}

	n.next()

	s = lazyq.First(n.Queue)
	if s.Node == nil {
		t.Errorf("TestJustInTimeNodeAllocation(): expected child to be allocated after next() but it was nil")
	}
}
