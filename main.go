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
	"github.com/tealeg/xlsx"
)

var (
	takeoffsites, landingsites string
	csvFile, xlsxFile          string
	writeCsv, writeXlsx        bool
	stdout                     bool
	distance                   int
	flights                    *igc.Flights
)

func init() {
	flag.StringVar(&takeoffsites, "takeoff", "Waypoints_Startplatz.gpx", "takeoff sites")
	flag.StringVar(&landingsites, "landing", "Waypoints_Landeplatz.gpx", "landing sites")
	flag.IntVar(&distance, "distance", 300, "maximal distance to the nearest known site")
	flag.BoolVar(&writeCsv, "csv", false, "write csv output")
	flag.StringVar(&csvFile, "csvfile", fmt.Sprintf("%s.csv", filepath.Base(os.Args[0])), "output file  (stdout writes to stdout)")
	flag.BoolVar(&writeXlsx, "xlsx", true, "write xlsx output")
	flag.StringVar(&xlsxFile, "xlsxfile", fmt.Sprintf("%s.xlsx", filepath.Base(os.Args[0])), "output file")
}

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()
	if len(flag.Args()) >= 1 {
		dir = flag.Args()[0]
	}
	takeoff, err := wpt.NewWaypoints(takeoffsites)
	if err == nil {
		igc.RegisterTakeoffSiteSource(takeoff)
	}
	landing, err := wpt.NewWaypoints(landingsites)
	if err == nil {
		igc.RegisterLandingSiteSource(landing)
	}
	igc.MaxDistance = distance

	flights = igc.NewFlights()
	err = filepath.Walk(dir, evaluate)
	if err != nil {
		log.Print(err)
	}
	sort.Sort(flights)

	stat, err := flightstat.NewFlightStat(flights)
	if err != nil {
		log.Fatal(err)
	}

	if writeCsv {
		if err := writeCsvOutput(csvFile, flights, stat); err != nil {
			log.Fatal(err)
		}
	}
	if writeXlsx {
		if err := writeXlsxOutput(xlsxFile, flights, stat); err != nil {
			log.Fatal(err)
		}
	}
}

func evaluate(path string, f os.FileInfo, err error) error {
	if f.IsDir() || !strings.HasSuffix(f.Name(), ".igc") {
		return nil
	}
	flight, err := igc.NewFlight(path)
	if err != nil {
		return err
	}
	flights.Add(flight)
	return nil
}

func writeCsvOutput(outfile string, f *igc.Flights, s *flightstat.FlightStat) error {
	out := os.Stdout
	if outfile != "stdout" {
		out, err := os.Create(outfile)
		if err != nil {
			return err
		}
		defer out.Close()
	}
	w := csv.NewWriter(out)
	f.Csv(w)
	s.Csv(w)
	w.Flush()
	return nil
}

func writeXlsxOutput(outfile string, flights *igc.Flights, stat *flightstat.FlightStat) error {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Flight Statistics")
	if err != nil {
		return err
	}

	// flights - header
	flights.Xlsx(sheet)
	sheet.AddRow()
	// statistics - data
	stat.Xlsx(sheet)

	err = file.Save(outfile)
	if err != nil {
		return err
	}
	return nil
}
