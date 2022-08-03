package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

var l = log.New(
	os.Stderr,
	"nstarsim: ",
	log.Ldate|log.Ltime|log.LUTC|log.Lshortfile,
)

// Record defines an interface that represents a data record that can be fed
// into Nested STAR.
type Record interface {
	Prepare() []string
}

func main() {
	filename := flag.String("filename", "", "Path to file containing CSV records.")
	flag.Parse()

	nstar := NewNestedSTAR()
	var wg sync.WaitGroup
	wg.Add(1)
	if err := parseCSVFile(*filename, nstar, &wg); err != nil {
		l.Fatalf("Failed to parse CSV file: %s", err)
	}
	// Wait until we're done parsing the CSV file before starting the
	// aggregation.
	wg.Wait()

	fmt.Println("type,k,frac,num_tags,num_leaf_tags,len_part_msmts,num_part_msmts")
	for _, k := range []int{5, 10, 25, 50, 75, 100, 125, 250, 375, 500, 625, 750, 875, 1000} {
		l.Printf("Aggregating for k=%d.", k)
		nstar.Aggregate(numLocGranularities, k)
	}
}
