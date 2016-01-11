package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/marcsauter/flightstat"
	"github.com/marcsauter/igc"
	"github.com/marcsauter/wpt"
)

var (
	takeoffsites, landingsites, outfile string
	distance                            int
	flights                             igc.Flights
)

func init() {
	flag.StringVar(&takeoffsites, "takeoff", "Waypoints_Startplatz.gpx", "takeoff sites")
	flag.StringVar(&landingsites, "landing", "Waypoints_Landeplatz.gpx", "landing sites")
	flag.IntVar(&distance, "distance", 300, "maximal distance to the nearest known site")
	flag.StringVar(&outfile, "outfile", fmt.Sprintf("%s.csv", filepath.Base(os.Args[0])), "output file")
}

func main() {
	flag.Parse()
	takeoff, err := wpt.NewWaypoints(takeoffsites)
	if err == nil {
		igc.RegisterTakeoffSiteSource(takeoff)
	}
	landing, err := wpt.NewWaypoints(landingsites)
	if err == nil {
		igc.RegisterLandingSiteSource(landing)
	}
	igc.MaxDistance = distance
	flights = igc.Flights{}
	err = filepath.Walk(flag.Args()[0], evaluate)
	if err != nil {
		log.Print(err)
	}
	sort.Sort(flights)
	stat := flightstat.NewFlightStat()
	out, err := os.Create(outfile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	w := csv.NewWriter(out)
	for _, f := range flights {
		if err := stat.Add(f); err != nil {
			log.Fatal(err)
		}
		if err := w.Write(f.Record()); err != nil {
			log.Fatal(err)
		}
	}
	for _, s := range stat.Output() {
		if err := w.Write(s); err != nil {
			log.Fatal(err)
		}
	}
	w.Flush()
}

func evaluate(path string, f os.FileInfo, err error) error {
	if f.IsDir() {
		return nil
	}
	if strings.HasSuffix(f.Name(), ".igc") {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		flight, err := igc.NewFlight(f)
		if err != nil {
			return err
		}
		flights = append(flights, flight)
	}
	return nil
}
