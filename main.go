package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"runtime"

	"github.com/jszwec/csvutil"
	geo "github.com/kellydunn/golang-geo"
)

const (
	stopAtDistance = 0.001
)

type StringMap map[string]string

// MarshalCSV generates the string to put additional columns in one
func (sm StringMap) MarshalCSV() ([]byte, error) {
	if len(sm) == 0 {
		return []byte{}, nil
	}
	return []byte(fmt.Sprint(sm)), nil
}

type Row struct {
	ID        int     `csv:"id"`
	Latitude  float64 `csv:"latitude"`
	Longitude float64 `csv:"longitude"`
	Others    map[string]string
	geo       *geo.Point
}

type Result struct {
	ID1      int       `csv:"id searched"`
	ID2      int       `csv:"id data"`
	Distance int       `csv:"meters"`
	Others1  StringMap `csv:"searched others"`
	Others2  StringMap `csv:"data others"`
}

func main() {
	var dataCsv, searchCsv, outFile string
	// Flags
	flag.StringVar(&dataCsv, "d", "", "Data file, for example giata.csv")
	flag.StringVar(&searchCsv, "s", "", "Data to look for example new_properties.csv")
	flag.StringVar(&outFile, "o", "", "Output CSV filename, if not specified, it prints in stdout")
	flag.Parse()

	if dataCsv == "" {
		fmt.Println("Missing data CSV name")
		os.Exit(2)
	}

	if searchCsv == "" {
		fmt.Println("Missing data CSV of date to search for")
		os.Exit(2)
	}

	data := loadCSV(dataCsv)
	toSearch := loadCSV(searchCsv)
	results := processParallel(data, toSearch)
	saveCSV(outFile, results)

}

func processParallel(data, toSearch []Row) (rows []Result) {
	p := runtime.NumCPU()
	s := len(toSearch) / p

	ch := make(chan []Result)
	dispatched := 0
	// Process segments en parallel
	for len(toSearch) > 0 {
		end := s
		if end > len(toSearch) {
			end = len(toSearch)
		}

		// Dispatch a segment
		go func(d1, d2 []Row) {
			r := process(d1, d2)
			ch <- r
		}(data, toSearch[:end])
		dispatched++

		toSearch = toSearch[end:]
	}

	// Collect the results
	for i := 0; i < dispatched; i++ {
		rr := <-ch
		rows = append(rows, rr...)
	}
	return
}

func process(data, toSearch []Row) (rows []Result) {
	for _, s := range toSearch {
		minDistance := math.MaxFloat64
		var picked Row
		for _, d := range data {
			if d.geo == nil {
				continue
			}
			if math.Abs(d.geo.Lat()-s.geo.Lat()) > 1 || math.Abs(d.geo.Lng()-s.geo.Lng()) > 1 {
				continue
			}
			dist := s.geo.GreatCircleDistance(d.geo)
			if dist < minDistance {
				minDistance = dist
				picked = d
				if dist < stopAtDistance {
					break
				}
			}
		}
		if minDistance == math.MaxFloat64 {
			continue
		}
		res := Result{
			ID1:      s.ID,
			ID2:      picked.ID,
			Distance: int(minDistance * 1000),
			Others1:  s.Others,
			Others2:  picked.Others,
		}
		rows = append(rows, res)
	}
	return
}

func saveCSV(filename string, results []Result) {
	var w *csv.Writer
	if filename == "" {
		w = csv.NewWriter(os.Stdout)
	} else {
		csvFile, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer csvFile.Close()
		w = csv.NewWriter(csvFile)
	}

	enc := csvutil.NewEncoder(w)
	for _, r := range results {
		if err := enc.Encode(r); err != nil {
			log.Fatal(err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

func loadCSV(filename string) (rows []Row) {
	csvFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		log.Fatal(err)
	}
	header := dec.Header()
	for {
		r := Row{Others: make(map[string]string)}

		if err := dec.Decode(&r); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		for _, i := range dec.Unused() {
			r.Others[header[i]] = dec.Record()[i]
		}
		// Create the geo point if lat and lgt are not zero
		if r.Latitude != 0 || r.Longitude != 0 {
			r.geo = geo.NewPoint(r.Latitude, r.Longitude)
		}
		rows = append(rows, r)
	}

	return
}
