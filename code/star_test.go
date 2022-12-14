package main

import "testing"

func initSTAR() (*nestedSTAR, int, int) {
	star := NewNestedSTAR()

	maxTags, threshold := 3, 5
	// Six full measurements.
	star.root.Add([]string{"US", "release", "windows"})
	star.root.Add([]string{"US", "release", "windows"})
	star.root.Add([]string{"US", "release", "windows"})
	star.root.Add([]string{"US", "release", "windows"})
	star.root.Add([]string{"US", "release", "windows"})
	star.root.Add([]string{"US", "release", "windows"})
	// Three partial measurements of length two, containing ["US", "release"].
	star.root.Add([]string{"US", "release", "linux"})
	star.root.Add([]string{"US", "release", "linux"})
	star.root.Add([]string{"US", "release", "macos"})
	// Two partial measurements of length one, containing ["US"].
	star.root.Add([]string{"US", "nightly", "windows"})
	star.root.Add([]string{"US", "beta", "windows"})
	// Five partial measurements of length one, containing ["CA"].
	star.root.Add([]string{"CA", "release", "windows"})
	star.root.Add([]string{"CA", "release", "windows"})
	star.root.Add([]string{"CA", "release", "windows"})
	star.root.Add([]string{"CA", "release", "windows"})
	star.root.Add([]string{"CA", "nightly", "windows"})
	// One discarded measurement.
	star.root.Add([]string{"MX", "release", "windows"})

	return star, maxTags, threshold
}

func TestRealSTAR(t *testing.T) {
	star, maxTags, threshold := initSTAR()
	state := star.root.Aggregate(maxTags, threshold, []string{})
	if !state.AddsUp() {
		t.Fatal("number of partial measurements don't add up")
	}

	expectedFull, expectedPartial := 6, 10
	if state.FullMsmts != expectedFull {
		t.Fatalf("expected %d but got %d full measurements.", expectedFull, state.FullMsmts)
	}
	if state.PartialMsmts != expectedPartial {
		t.Fatalf("expected %d but got %d partial measurements.", expectedPartial, state.PartialMsmts)
	}

	expectedLens := map[int]int{
		3: 6, // Six measurements of length three.
		2: 3, // Three measurements of length two.
		1: 7, // Seven measurements of length one.
	}
	if !isMapEqual(state.LenPartialMsmts, expectedLens) {
		t.Fatalf("expected %v but got %v.", expectedLens, state.LenPartialMsmts)
	}
}

func isMapEqual(m1, m2 map[int]int) bool {
	if len(m1) != len(m2) {
		return false
	}

	for key, v1 := range m1 {
		v2, exists := m2[key]
		if !exists {
			return false
		}
		if v1 != v2 {
			return false
		}
	}
	return true
}

func TestNumLeafs(t *testing.T) {
	var n, expected int
	star, _, _ := initSTAR()

	n = star.root.NumNodes()
	expected = 10
	if n != expected {
		t.Fatalf("expected %d but got %d nodes in tree.", expected, n)
	}

	n = star.root.NumTags()
	expected = 17
	if n != expected {
		t.Fatalf("expected %d but got %d tags in tree.", expected, n)
	}

	n = star.root.NumLeafTags()
	expected = 8
	if n != expected {
		t.Fatalf("expected %d but got %d leaf tags in tree.", expected, n)
	}
}

func TestAggregationState(t *testing.T) {
	s1 := NewAggregationState()
	s2 := NewAggregationState()

	s1.AddLenTags(0, 10)
	s1.AddLenTags(1, 5)
	s2.AddLenTags(1, 15)
	s1.Augment(s2)
	if s1.LenPartialMsmts[0] != 10 {
		t.Fatalf("expected 10 but got %d", s1.LenPartialMsmts[0])
	}
	if s1.LenPartialMsmts[1] != 20 {
		t.Fatalf("expected 20 but got %d", s1.LenPartialMsmts[0])
	}
}
