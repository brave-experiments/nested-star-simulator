package main

const (
	rootDepth = 1
)

type nodeInfo struct {
	Num  int
	Next *node
}

type node struct {
	// E.g., "US"    -> nodeInfo{}
	//       "Linux" -> nodeInfo{}
	ValueToInfo map[string]*nodeInfo
}

func (n *node) Aggregate(maxDepth, k int, m []string) *aggregationState {
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
			state.AddLenTags(depth, info.Num)
			continue
		}
		if info.Next == nil {
			l.Fatalf("Encountered incomplete measurement: %s\n", m)
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

func (n *node) Add(orderedMsmt []string) {
	info, exists := n.ValueToInfo[orderedMsmt[0]]
	if !exists {
		info = &nodeInfo{Num: 1}
		n.ValueToInfo[orderedMsmt[0]] = info
	} else {
		info.Num++
	}

	if len(orderedMsmt[1:]) > 0 {
		if info.Next == nil {
			newNode := &node{ValueToInfo: make(map[string]*nodeInfo)}
			info.Next = newNode
			newNode.Add(orderedMsmt[1:])
		} else {
			info.Next.Add(orderedMsmt[1:])
		}
	}
}

func (n *node) NumTags() int {
	var num = len(n.ValueToInfo)

	for _, info := range n.ValueToInfo {
		if info.Next != nil {
			num += info.Next.NumTags()
		}
	}
	return num
}

func (n *node) NumNodes() int {
	var num = 1

	for _, info := range n.ValueToInfo {
		if info.Next != nil {
			num += info.Next.NumNodes()
		}
	}
	return num
}

func (n *node) NumLeafTags() int {
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
