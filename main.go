package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/marcsauter/igc"
	"github.com/marcsauter/wpt"
)

var (
	takeoffsites, landingsites string
	distance                   int
	flights                    []*igc.Flight
)

func init() {
	flag.StringVar(&takeoffsites, "takeoff", "Waypoints_Startplatz.gpx", "takeoff sites")
	flag.StringVar(&landingsites, "landing", "Waypoints_Landeplatz.gpx", "landing sites")
	flag.IntVar(&distance, "distance", 300, "maximal distance to the nearest known site")
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
	flights = []*igc.Flight{}
	err = filepath.Walk(flag.Args()[0], evaluate)
	if err != nil {
		log.Print(err)
	}
	w := csv.NewWriter(os.Stdout)
	for _, f := range flights {
		if err := w.Write(f.Stat()); err != nil {
			log.Fatalln("error writing record to csv:", err)
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
