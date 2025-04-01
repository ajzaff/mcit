package mcit

import "iter"

// LazyQueue contains the multi-armed bandit heap queue.
//
// It has a few differences from a regular heap queue:
// 	* an append method which adds a lazy element to the heap.
//	* a next method which handles promoting the lazy elements and returning the max element.
type LazyQueue struct {
	// lazyIndex is set to the index of the first lazy element.
	// When equal to Len(), indicates all elements are heapified.
	lazyIndex int
	// Bandits is a slice of stats representing chosing different actions.
	Bandits []Stat
}

func (h *LazyQueue) hasLazyElements() bool { return h.lazyIndex < h.Len() }

// StatSeq returns an iterator over the Stats for a Node in the correct priority order.
func (h *LazyQueue) StatSeq() iter.Seq[Stat] {
	return func(yield func(Stat) bool) {
		for i := h.lazyIndex; i < h.Len(); i++ {
			if !yield(h.Bandits[i]) {
				return
			}
		}
		for i := range h.lazyIndex {
			if !yield(h.Bandits[i]) {
				return
			}
		}
	}
}

func (h *LazyQueue) next() Stat {
	if h.hasLazyElements() {
		// We have at least one node which has never been tried before.
		// Use this time to fix the position in the heap so we can select it.
		// Nodes which have never been tried before always take priority.
		//
		// Waiting until now to fix this position is largely an optimization
		// as we don't expect the majority of nodes of large trees to be tried
		// we don't need to waste time with the O(log N) heap.Push operation.
		h.up(h.lazyIndex)
		h.lazyIndex++
	}
	// NOTE: We always take the first action.
	// If we ever implemented a temperature feature, we'd need to keep track of this index.
	return h.Bandits[0]
}

func (h *LazyQueue) append(x Stat) { h.Bandits = append(h.Bandits, x) }

func (h *LazyQueue) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.less(j, i) {
			break
		}
		h.swap(i, j)
		j = i
	}
}

func (h *LazyQueue) down(i0 int) bool {
	i := i0
	n := h.lazyIndex
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.less(j, i) {
			break
		}
		h.swap(i, j)
		i = j
	}
	return i > i0
}

func (h LazyQueue) Len() int      { return len(h.Bandits) }
func (h LazyQueue) swap(i, j int) { h.Bandits[i], h.Bandits[j] = h.Bandits[j], h.Bandits[i] }
func (h LazyQueue) less(i, j int) bool {
	if ui, uj := h.Bandits[i].Priority, h.Bandits[j].Priority; ui != uj {
		// Higher priority nodes first.
		return ui > uj
	}
	// When priorities are equal (often +∞), fall back to prior comparison.
	return h.Bandits[i].Prior > h.Bandits[j].Prior
}
