package mcit

// RunResults contains results returned from the user search function.
type RunResults struct {
	// Expand is the slice of actions available from this node which should be available as children.
	// Expand may be a subset of all possible actions. See Replace for tips
	// implementing partial expansion.
	Expand []string
	// Priors is a slice of prior values to apply to new expanded nodes.
	// Priors may be empty in which case the default prior value is used.
	Priors []float64
	// Replace can be set to true when the node should be returned to the frontier queue.
	// This is useful when allowing nodes to be partially expanded on each new visit OR
	// when the node is a leaf node of the current search and we want to repeatedly explore it.
	Replace  bool
	Minimize bool
	Count    float64
	Value    float64
	// Payload contains optional user generated payload to store in the resulting tree Node.
	// This can be a reference to the user land state of the node which can be more convenient
	// when direct replay is difficult.
	Payload any
	Done    bool
}

// NodeSelector provides a choice of methods to select the current frontier node.
// Either a slice of actions, or a reference to the Payload returned from RunResults.
type NodeSelector struct {
	Actions []string
	Payload any
}

// Func is a search function containing user code which selects a frontier node and returns the results of experiments on it.
type Func func(selector NodeSelector) (results RunResults)

// Continuation is a structure which contains a root node to pass to continue a previous search from memory.
type Continuation struct {
	root *Node
}
