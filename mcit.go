package mcit

import (
	"container/heap"
	"math"
	"slices"
)

type Func func(actions []string) (expand []string, priors []float64, exhaust bool, results RunResults, done bool)

func Search(runFn Func) SearchResults {
	// 0. Initialize state.
	//	0a. Initialize the search tree.
	root := &node{
		results:     RunResults{},
		ucb1:        math.Inf(+1),
		frontierIdx: -1,
		children:    make(map[string]*node),
	}
	//	0b. Initialize a frontier heap.
	frontier := []*node{}
	heap.Push((*byUCB1)(&frontier), root)

	// 0c. Initialize search stats.
	iters := int64(0)
	nodeCount := int64(0)
	maxFrontierSize := int64(0)
	exhaustedNodes := int64(0)
	minSampledPriority := math.Inf(+1)

	replay := make([]string, 0, 64)
	for ; ; iters++ {
		// 1. Select a frontier node from the frontier heap and construct replay actions.
		curr := frontier[0]

		if size := int64(len(frontier)); maxFrontierSize < size {
			maxFrontierSize = size
		}

		if p := curr.ucb1; p < minSampledPriority {
			minSampledPriority = p
		}

		replay = curr.getLine(replay)

		// 2. Run simulations at the frontier node.
		expand, priors, exhaust, results, done := runFn(replay)
		if done { //	2a. (optional) Set stop flag to stop search.
			break
		}

		// 	2a. (optional) Expand the node, and add children to the state.
		for i, e := range expand {
			if _, found := curr.children[e]; found {
				continue
			}
			nodeCount++
			child := &node{
				height:      curr.height + 1,
				action:      e,
				parent:      curr,
				results:     RunResults{},
				prior:       1,
				ucb1:        math.Inf(+1),
				frontierIdx: -1,
				children:    make(map[string]*node),
			}

			//	2aa. (optional) Priors, if provided, should match the slice of expanded nodes.
			if len(priors) > 0 {
				child.prior = priors[i]
			}

			curr.children[e] = child
			heap.Push((*byUCB1)(&frontier), child)
		}

		// 	2b. (optional) Exhaust the frontier node (or keep it around for next time).
		//                 The exhaust logic by default requires that some other conditions hold
		//                 To avoid the simulation running out of frontier values.
		if exhaust && (len(expand) > 0 || results.Count == 0) {
			exhaustedNodes++
			heap.Remove((*byUCB1)(&frontier), curr.frontierIdx)
		}

		// 	2c. Backpropagate the results up the tree and fix the frontier nodes.
		for head := curr; head != nil; head = head.parent {
			head.results.Count += results.Count
			head.results.Value += results.Value
			head.recomputeUCB1()
			if head.frontierIdx != -1 {
				heap.Fix((*byUCB1)(&frontier), head.frontierIdx)
			}
		}

		// 3. Restart from step 1.
	}

	// 4. Summarize the results.
	rootChildResults := make(map[string]RunResults)
	for _, child := range root.children {
		rootChildResults[child.action] = child.results
	}

	// Run histogram buckets (14 bins).
	// [0 1 2 4 8 16 32 64 128 256 512 1024 2048 4096+]
	runHist := makeHist(runBins)
	fillHist(runHist, root, func(n *node) float64 { return n.results.Count })

	// Fill score hist.
	scoreHist := makeHist(scoreBins)
	fillHist(scoreHist, root, func(n *node) float64 { return n.score() })

	// Fill priority hist.
	priorityHist := makeHist(priorityBins)
	fillHist(priorityHist, root, func(n *node) float64 { return n.ucb1 })

	bestChild := getMaxChild(root)

	mostRunChild := getMostRunChild(root)

	bestLine := getBestLine(root)

	mostRunLine := getMostRunLine(root)

	bestNode := new(node)
	getBestNode(bestNode, root)

	leafCount := countLeaves(root)

	return SearchResults{
		RootResults:        root.results,
		RootChildResults:   rootChildResults,
		Iterations:         iters,
		NodeCount:          nodeCount,
		LeafCount:          leafCount,
		MaxFrontierSize:    maxFrontierSize,
		ExhaustedNodes:     exhaustedNodes,
		BestChild:          bestChild.action,
		MostRunChild:       mostRunChild.action,
		BestLine:           bestLine.getLine(nil),
		BestNode:           bestNode.getLine(nil),
		MostRunLine:        mostRunLine.getLine(nil),
		MinSampledPriority: minSampledPriority,
		BestNodeScore:      bestNode.score(),
		MostRunScore:       mostRunLine.score(),
		BestChildScore:     bestChild.score(),
		BestLineScore:      bestLine.score(),
		RunHist:            runHist,
		ScoreHist:          scoreHist,
		PriorityHist:       priorityHist,
	}
}

func fillHist(hist Hist, root *node, valueFn func(*node) float64) {
	hist.Insert(valueFn(root))
	for _, child := range root.children {
		fillHist(hist, child, valueFn)
	}
}

func countLeaves(root *node) int64 {
	if len(root.children) == 0 {
		return 1
	}
	leaves := int64(0)
	for _, child := range root.children {
		leaves += countLeaves(child)
	}
	return leaves
}

func getBestNode(bestNode *node, root *node) {
	if bestNode.score() < root.score() {
		*bestNode = *root
	}
	for _, child := range root.children {
		getBestNode(bestNode, child)
	}
}

func getBestLine(root *node) *node { return getSelectLine(root, getMaxChild) }

func getMaxChild(root *node) *node {
	var maxChild *node
	for _, child := range root.children {
		if maxChild.score() < child.score() {
			maxChild = child
		}
	}
	return maxChild
}

func getMinChild(root *node) *node {
	var minChild *node
	for _, child := range root.children {
		minScore := minChild.score()
		if childScore := child.score(); childScore != math.Inf(-1) && minScore == math.Inf(-1) || childScore < minChild.score() {
			minChild = child
		}
	}
	return minChild
}

func getMostRunLine(root *node) *node { return getSelectLine(root, getMostRunChild) }

func getMostRunChild(root *node) *node {
	var mostRunChild *node
	for _, child := range root.children {
		if child.results.Count > 0 && (mostRunChild == nil || mostRunChild.results.Count < child.results.Count) {
			mostRunChild = child
		}
	}
	return mostRunChild
}

func getSelectLine(root *node, selectFn func(*node) *node) *node {
	for len(root.children) > 0 {
		next := selectFn(root)
		if next == nil {
			break
		}
		root = next
	}
	return root
}

type SearchResults struct {
	RootResults        RunResults
	RootChildResults   map[string]RunResults
	BestChild          string
	MostRunChild       string
	BestLine           []string
	BestNode           []string
	MostRunLine        []string
	MinSampledPriority float64
	BestChildScore     float64
	BestLineScore      float64
	BestNodeScore      float64
	MostRunScore       float64
	Iterations         int64
	NodeCount          int64
	LeafCount          int64
	MaxFrontierSize    int64
	ExhaustedNodes     int64
	RunHist            Hist
	ScoreHist          Hist
	PriorityHist       Hist
}

type RunResults struct {
	Count float64
	Value float64
}

type node struct {
	height      int
	action      string
	parent      *node
	results     RunResults
	ucb1        float64
	prior       float64
	frontierIdx int
	children    map[string]*node
}

const exploreTerm = math.Pi

func (n *node) recomputeUCB1() {
	n.ucb1 = n.computeUCB1()
}

func (n *node) computeUCB1() float64 {
	if n.results.Count == 0 {
		return math.Inf(+1)
	}
	return (n.results.Value + n.prior*exploreTerm) / n.results.Count
}

func (n *node) score() float64 {
	if n == nil || n.results.Count == 0 {
		return math.Inf(-1)
	}
	return n.results.Value / n.results.Count
}

func (n *node) getLine(line []string) []string {
	if n == nil {
		return line[:0]
	}
	line = slices.Grow(line[:0], n.height)[:n.height]
	for i := len(line) - 1; n.parent != nil; i, n = i-1, n.parent {
		line[i] = n.action
	}
	return line
}

type byUCB1 []*node

func (a byUCB1) Len() int { return len(a) }
func (a byUCB1) Swap(i, j int) {
	a[i].frontierIdx, a[j].frontierIdx = a[j].frontierIdx, a[i].frontierIdx
	a[i], a[j] = a[j], a[i]
}
func (a byUCB1) Less(i, j int) bool { return a[i].ucb1 > a[j].ucb1 }
func (a *byUCB1) Push(x any)        { n := len(*a); *a = append(*a, x.(*node)); x.(*node).frontierIdx = n }
func (a *byUCB1) Pop() any {
	n := len(*a) - 1
	x := (*a)[n]
	*a = (*a)[:n]
	x.frontierIdx = -1
	return x
}
