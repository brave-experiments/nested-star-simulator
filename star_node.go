package main

const (
	rootDepth = 1
)

type NodeInfo struct {
	Num  int
	Next *Node
}

type Node struct {
	// E.g., "US"    -> NodeInfo{}
	//       "Linux" -> NodeInfo{}
	ValueToInfo map[string]*NodeInfo
}

func (n *Node) Aggregate(maxDepth, k int, m []string) *AggregationState {
	state := NewAggregationState()
	depth := len(m) + 1

	// Iterate over all values where we are in the tree, e.g., "US", "FR", ...
	for value, info := range n.ValueToInfo {
		// We don't meet our k-anonymity threshold for the given value.
		if info.Num < k {
			continue
		}

		// We've reached the last tag, i.e., we fully unlocked a measurement.
		if depth == maxDepth {
			state.FullMsmts += info.Num
			continue
		}
		if info.Next == nil {
			elog.Printf("ERROR: Encountered incomplete measurement: %s\n", m)
			continue
		}

		// Go deeper down the tree, and try to unlock our next tag.
		subState := info.Next.Aggregate(maxDepth, k, append(m, value))
		state.Augment(subState)

		numNewlyUnlocked := info.Num - subState.FullMsmts - subState.AlreadyCounted
		state.AddLenTags(depth, numNewlyUnlocked)
		state.AlreadyCounted += numNewlyUnlocked

		// Once we're back at our root node, determine the total number of
		// partial measurements.
		if depth == rootDepth {
			state.PartialMsmts += info.Num - subState.FullMsmts
		}
	}

	return state
}

func (n *Node) Add(orderedMsmt []string) {
	info, exists := n.ValueToInfo[orderedMsmt[0]]
	if !exists {
		info = &NodeInfo{Num: 1}
		n.ValueToInfo[orderedMsmt[0]] = info
	} else {
		info.Num++
	}

	if len(orderedMsmt[1:]) > 0 {
		if info.Next == nil {
			newNode := &Node{ValueToInfo: make(map[string]*NodeInfo)}
			info.Next = newNode
			newNode.Add(orderedMsmt[1:])
		} else {
			info.Next.Add(orderedMsmt[1:])
		}
	}
}

func (n *Node) NumTags() int {
	var num = len(n.ValueToInfo)

	for _, info := range n.ValueToInfo {
		if info.Next != nil {
			num += info.Next.NumTags()
		}
	}
	return num
}

func (n *Node) NumNodes() int {
	var num = 1

	for _, info := range n.ValueToInfo {
		if info.Next != nil {
			num += info.Next.NumNodes()
		}
	}
	return num
}

func (n *Node) NumLeafTags() int {
	var num int

	for _, info := range n.ValueToInfo {
		if info.Next == nil {
			num++
		} else {
			num += info.Next.NumLeafTags()
		}
	}
	return num
}
