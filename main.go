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
	distance                   int
	flights                    igc.Flights
)

func init() {
	flag.StringVar(&takeoffsites, "takeoff", "Waypoints_Startplatz.gpx", "takeoff sites")
	flag.StringVar(&landingsites, "landing", "Waypoints_Landeplatz.gpx", "landing sites")
	flag.IntVar(&distance, "distance", 300, "maximal distance to the nearest known site")
	flag.BoolVar(&writeCsv, "csv", false, "write csv output")
	flag.StringVar(&csvFile, "csvfile", fmt.Sprintf("%s.csv", filepath.Base(os.Args[0])), "output file")
	flag.BoolVar(&writeXlsx, "xlsx", true, "write xlsx output")
	flag.StringVar(&xlsxFile, "xlsxfile", fmt.Sprintf("%s.xlsx", filepath.Base(os.Args[0])), "output file")
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(0)
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
	flights = igc.Flights{}
	err = filepath.Walk(flag.Args()[0], evaluate)
	if err != nil {
		log.Print(err)
	}
	sort.Sort(flights)
	stat, err := flightstat.NewFlightStat(&flights)
	if err != nil {
		log.Fatal(err)
	}
	data := append(*flights.Output(), *stat.Output()...)
	if writeCsv {
		if err := writeCsvOutput(csvFile, &data); err != nil {
			log.Fatal(err)
		}
	}
	if writeXlsx {
		if err := writeXlsxOutput(xlsxFile, &data); err != nil {
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
	flights = append(flights, flight)
	return nil
}

func writeCsvOutput(outfile string, data *[][]string) error {
	out, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer out.Close()
	w := csv.NewWriter(out)
	//for _, f := range flights {
	//if err := w.Write(f.Record()); err != nil {
	//log.Fatal(err)
	//}
	//}
	for _, s := range *data {
		if err := w.Write(s); err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}

func writeXlsxOutput(outfile string, data *[][]string) error {
	//out, err := os.Create(outfile)
	//if err != nil {
	//return err
	//}
	//defer out.Close()
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Filght Statistics")
	if err != nil {
		return err
	}
	for _, r := range *data {
		row := sheet.AddRow()
		for _, c := range r {
			cell := row.AddCell()
			cell.Value = c
		}
	}
	err = file.Save(outfile)
	if err != nil {
		return err
	}
	return nil
}
