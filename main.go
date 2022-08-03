package main

import (
	"flag"
	"log"
	"os"
	"sync"
)

const (
	defaultK = 10
)

var elog = log.New(
	os.Stderr,
	"nstarsim: ",
	log.Ldate|log.Ltime|log.LUTC|log.Lshortfile,
)

// Report defines an interface that represents a report in our briefcase.  A
// report must be able to return its crowd ID and payload; and it must be
// marshal-able.
type Report interface {
	Prepare() []string
}

func main() {
	filename := flag.String("filename", "", "Filename containing CSV records.")
	k := flag.Int("k", defaultK, "The k-anonymity threshold.")
	flag.Parse()

	nstar := NewNestedSTAR(*k)

	var wg sync.WaitGroup
	wg.Add(1)
	if err := parseCSVFile(*filename, nstar, &wg); err != nil {
		elog.Fatalf("Failed to parse CSV file: %s", err)
	}
	// Wait until we're done parsing the CSV file before starting the
	// aggregation.
	wg.Wait()
	elog.Printf("# nodes: %d", nstar.NumNodes())
	elog.Printf("# tags: %d", nstar.NumTags())
	elog.Printf("# leaf tags: %d", nstar.NumLeafTags())

	nstar.Aggregate(numLocGranularities)
}
