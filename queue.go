package mcit

import "iter"

// lazyQueue contains the multi-armed bandit heap queue.
type lazyQueue struct {
	lazyIndex int // when lazyIndex < Len(), we have elements to fix.
	// Bandits is a slice of stats representing different choices of actions (bandit arms).
	Bandits []Stat
}

func (h *lazyQueue) hasLazyElements() bool { return h.lazyIndex < h.Len() }

// StatSeq returns an iterator over the Stats for a Node in the correct priority order.
func (h *lazyQueue) StatSeq() iter.Seq[Stat] {
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

func (h lazyQueue) top() Stat { return h.Bandits[0] }

func (h *lazyQueue) append(x Stat) { h.Bandits = append(h.Bandits, x) }

// upLazy calls up on the lazy element.
// upLazy panics if !h.hasLazyElements().
func (h *lazyQueue) upLazy() {
	h.up(h.lazyIndex)
	h.lazyIndex++
}

func (h lazyQueue) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.less(j, i) {
			break
		}
		h.swap(i, j)
		j = i
	}
}

func (h lazyQueue) down(i0 int) bool {
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

func (h lazyQueue) Len() int      { return len(h.Bandits) }
func (h lazyQueue) swap(i, j int) { h.Bandits[i], h.Bandits[j] = h.Bandits[j], h.Bandits[i] }
func (h lazyQueue) less(i, j int) bool {
	if ui, uj := h.Bandits[i].Priority, h.Bandits[j].Priority; ui != uj {
		// Higher priority nodes first.
		return ui > uj
	}
	// When priorities are equal (often +∞), fall back to prior comparison.
	return h.Bandits[i].Prior > h.Bandits[j].Prior
}
