package main

import (
	"fmt"
	"sync"
)

// NestedSTAR simulates the execution of Nested STAR over P3A measurements.
// Note that this struct doesn't actually implement Nested STAR; it merely
// simulates its nesting (represented as a tree of nodes) to produce CSV output
// that allows us to explore the privacy and utility tradeoff.
type NestedSTAR struct {
	sync.WaitGroup
	inbox           chan []Report
	root            *Node
	k               int
	numMeasurements int
}

func (s *NestedSTAR) NumTags() int {
	return s.root.NumTags()
}
func (s *NestedSTAR) NumLeafTags() int {
	return s.root.NumLeafTags()
}
func (s *NestedSTAR) NumNodes() int {
	return s.root.NumNodes()
}

// NewNestedSTAR returns a new NestedSTAR object.
func NewNestedSTAR(k int) *NestedSTAR {
	return &NestedSTAR{
		inbox: make(chan []Report),
		root:  &Node{make(map[string]*NodeInfo)},
		k:     k,
	}
}

// AddReports adds the given reports to Nested STAR.
func (s *NestedSTAR) AddReports(reports []Report) {
	s.numMeasurements += len(reports)
	for _, r := range reports {
		s.root.Add(r.Prepare())
	}
}

func frac(a, b int) float64 {
	if b == 0 {
		return 0
	}
	return float64(a) / float64(b)
}

// Aggregate aggregates Nested STAR's measurements.  The argument 'numAttrs'
// refers to the number of attributes in a record.
func (s *NestedSTAR) Aggregate(numAttrs int) {
	state := s.root.Aggregate(numAttrs, s.k, []string{})
	if !state.AddsUp() {
		elog.Fatal("Number of partial measurements don't add up.")
	}

	// Determine
	for key := 1; key <= numAttrs; key++ {
		num, exists := state.LenPartialMsmts[key]
		if !exists {
			num = 0
		}
		fmt.Printf("LenPartMsmt,%d,0,0,0,%d,%d\n",
			s.k,
			key,
			num)
	}
	fracFull := frac(state.FullMsmts, s.numMeasurements) * 100
	fracPart := frac(state.PartialMsmts, s.numMeasurements) * 100
	elog.Printf("%d (%.1f%%) full, %d (%.1f%%) partial out of %d; %.1f%% lost\n",
		state.FullMsmts,
		fracFull,
		state.PartialMsmts,
		fracPart,
		s.numMeasurements,
		100-fracFull-fracPart)
	fmt.Printf("Partial,%d,%.3f,%d,%d,0,0\n",
		s.k,
		frac(state.PartialMsmts, s.numMeasurements),
		s.root.NumTags(),
		s.root.NumLeafTags())
}

type AggregationState struct {
	FullMsmts       int
	PartialMsmts    int
	AlreadyCounted  int
	LenPartialMsmts map[int]int
}

func NewAggregationState() *AggregationState {
	return &AggregationState{
		LenPartialMsmts: make(map[int]int),
	}
}

func (s *AggregationState) String() string {
	return fmt.Sprintf("%d full, %d partial msmts.", s.FullMsmts, s.PartialMsmts)
}

func (s *AggregationState) AddLenTags(key, value int) {
	num, exists := s.LenPartialMsmts[key]
	if !exists {
		s.LenPartialMsmts[key] = value
	} else {
		s.LenPartialMsmts[key] = num + value
	}
}

func (s *AggregationState) Augment(s2 *AggregationState) {
	s.FullMsmts += s2.FullMsmts
	s.PartialMsmts += s2.PartialMsmts
	s.AlreadyCounted += s2.AlreadyCounted
	for key, value := range s2.LenPartialMsmts {
		s.AddLenTags(key, value)
	}
}

// AddsUp returns true if the number n-length partial measurements adds up to
// the total number of partial measurements.  The purpose of this function is
// to ensure algorithmic correctness.
func (s *AggregationState) AddsUp() bool {
	totalPartial := 0
	for _, num := range s.LenPartialMsmts {
		totalPartial += num
	}
	return s.PartialMsmts == totalPartial
}
