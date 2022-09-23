package main

import (
	"fmt"
	"sync"
)

// nestedSTAR simulates the execution of Nested STAR over P3A measurements.
// Note that this struct doesn't actually implement Nested STAR; it merely
// simulates its nesting (represented as a tree of nodes) to produce CSV output
// that allows us to explore the privacy and utility tradeoff.
type nestedSTAR struct {
	sync.WaitGroup
	inbox           chan []Record
	root            *node
	numMeasurements int
}

func (s *nestedSTAR) NumTags() int {
	return s.root.NumTags()
}
func (s *nestedSTAR) NumLeafTags() int {
	return s.root.NumLeafTags()
}
func (s *nestedSTAR) NumNodes() int {
	return s.root.NumNodes()
}

// NewNestedSTAR returns a new NestedSTAR object.
func NewNestedSTAR() *nestedSTAR {
	return &nestedSTAR{
		inbox: make(chan []Record),
		root:  &node{make(map[string]*nodeInfo)},
	}
}

// AddRecords adds the given records to Nested STAR.
func (s *nestedSTAR) AddRecords(records []Record) {
	s.numMeasurements += len(records)
	for _, r := range records {
		s.root.Add(r.Prepare())
	}
}

func frac(a, b int) float64 {
	if b == 0 {
		return 0
	}
	return float64(a) / float64(b)
}

func sumValues(m map[int]int) int {
	total := 0
	for _, v := range m {
		total += v
	}
	return total
}

// Aggregate aggregates Nested STAR's measurements.  The argument 'numAttrs'
// refers to the number of attributes in a record.
func (s *nestedSTAR) Aggregate(numAttrs, k int) {
	state := s.root.Aggregate(numAttrs, k, []string{})
	if !state.AddsUp() {
		l.Fatal("Number of partial measurements don't add up.")
	}

	if len(state.LenPartialMsmts) >= numAttrs {
		l.Fatalf("Expected < %d attributes but got %d.",
			numAttrs, len(state.LenPartialMsmts))
	}

	// Start with partial measurements that unlocked at least 1 attribute and
	// iterate to numAttrs-1, which are partial measurements that are only
	// missing a single attribute.
	totalNum := sumValues(state.LenPartialMsmts)
	for key := 1; key < numAttrs; key++ {
		num, exists := state.LenPartialMsmts[key]
		if !exists {
			num = 0
		}
		fmt.Printf("LenPartMsmt,%d,0,0,0,%d,%.2f\n",
			k,
			key,
			frac(num, totalNum))
	}
	fracFull := frac(state.FullMsmts, s.numMeasurements) * 100
	fracPart := frac(state.PartialMsmts, s.numMeasurements) * 100
	l.Printf("%d (%.1f%%) full, %d (%.1f%%) partial out of %d; %.1f%% lost\n",
		state.FullMsmts,
		fracFull,
		state.PartialMsmts,
		fracPart,
		s.numMeasurements,
		100-fracFull-fracPart)
	// We only print partial measurements here because the number of full
	// measurements is simply the number of total measurements subtracted by
	// the number of partial measurements.
	fmt.Printf("Partial,%d,%.3f,%d,%d,0,0\n",
		k,
		frac(state.PartialMsmts, s.numMeasurements),
		s.root.NumTags(),
		s.root.NumLeafTags())
}

type aggregationState struct {
	FullMsmts       int
	PartialMsmts    int
	AlreadyCounted  int
	LenPartialMsmts map[int]int
}

func NewAggregationState() *aggregationState {
	return &aggregationState{
		LenPartialMsmts: make(map[int]int),
	}
}

func (s *aggregationState) String() string {
	return fmt.Sprintf("%d full, %d partial msmts.", s.FullMsmts, s.PartialMsmts)
}

func (s *aggregationState) AddLenTags(key, value int) {
	num, exists := s.LenPartialMsmts[key]
	if !exists {
		s.LenPartialMsmts[key] = value
	} else {
		s.LenPartialMsmts[key] = num + value
	}
}

func (s *aggregationState) Augment(s2 *aggregationState) {
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
func (s *aggregationState) AddsUp() bool {
	totalPartial := 0
	for _, num := range s.LenPartialMsmts {
		totalPartial += num
	}
	return s.PartialMsmts == totalPartial
}
