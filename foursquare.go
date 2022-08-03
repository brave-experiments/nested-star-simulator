package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	numLocGranularities = 7
	numCSVFields        = numLocGranularities + 1 + 1
)

type latLon struct {
	lat float32
	lon float32
}

type foursquareRec struct {
	countryCode string
	locations   []*latLon
}

func (f foursquareRec) String() string {
	var locs string
	for _, loc := range f.locations {
		locs += fmt.Sprintf(" %.7f;%.7f", loc.lat, loc.lon)
	}

	return fmt.Sprintf("%s: %s", f.countryCode, locs)
}

func (f *foursquareRec) Prepare() []string {
	var s []string
	for _, loc := range f.locations {
		s = append(s, fmt.Sprintf("%7f;%7f", loc.lat, loc.lon))
	}
	return s
}

func parseLatLon(s string) (*latLon, error) {
	pair := strings.Split(s, ";")
	if len(pair) != 2 {
		return nil, errors.New("expected exactly one ';' separator in lat/lon")
	}

	// The second column of our CSV data contains number like: 41.;28.
	// ParseFloat is not happy with that, so we are removing the trailing
	// period.
	if strings.HasSuffix(pair[0], ".") {
		pair[0] = strings.ReplaceAll(pair[0], ".", "")
		pair[1] = strings.ReplaceAll(pair[1], ".", "")
	}

	lat, err := strconv.ParseFloat(pair[0], 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse lat %q: %s", pair[0], err)
	}
	lon, err := strconv.ParseFloat(pair[1], 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse lon %q: %s", pair[1], err)
	}

	return &latLon{
		lat: float32(lat),
		lon: float32(lon),
	}, nil
}

func parseRawRecord(record []string) (*foursquareRec, error) {
	if len(record) != numCSVFields {
		return nil, fmt.Errorf("expected %d but got %d CSV fields", numCSVFields, len(record))
	}

	f := &foursquareRec{}
	f.countryCode = strings.TrimSpace(record[0])
	for i := 1; i < numLocGranularities+1; i++ {

		latLon, err := parseLatLon(record[i])
		if err != nil {
			return nil, err
		}
		f.locations = append(f.locations, latLon)
	}
	return f, nil
}

func parseCSVFile(filename string, nstar *nestedSTAR, wg *sync.WaitGroup) error {
	defer func() {
		wg.Done()
	}()
	l.Printf("Opening %q for processing.", filename)
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}

	reader := csv.NewReader(bufio.NewReader(fd))
	numRecs := 0
	var strRec []string
	for {
		strRec, err = reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			l.Printf("Error reading from CSV reader: %s", err)
		}
		fsRec, err := parseRawRecord(strRec)
		if err != nil {
			l.Printf("Error parsing Foursquare record: %s", err)
		}
		numRecs++
		if numRecs%1000000 == 0 {
			l.Printf("Parsed %dM records.", numRecs/1000000)
		}

		nstar.AddRecords([]Record{fsRec})
	}
	l.Printf("Parsed %d valid records.", numRecs)

	return nil
}
