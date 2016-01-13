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

	flights = igc.Flights{}
	err = filepath.Walk(dir, evaluate)
	if err != nil {
		log.Print(err)
	}
	sort.Sort(flights)

	stat, err := flightstat.NewFlightStat(&flights)
	if err != nil {
		log.Fatal(err)
	}

	fdata := flights.Output()
	sdata := stat.Output()
	if writeCsv {
		if err := writeCsvOutput(csvFile, fdata, sdata); err != nil {
			log.Fatal(err)
		}
	}
	if writeXlsx {
		if err := writeXlsxOutput(xlsxFile, fdata, sdata); err != nil {
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

func writeCsvOutput(outfile string, fdata, sdata *[][]string) error {
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
	for _, s := range *fdata {
		if err := w.Write(s); err != nil {
			return err
		}
	}
	for _, s := range *sdata {
		if err := w.Write(s); err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}

func writeXlsxOutput(outfile string, fdata, sdata *[][]string) error {
	//out, err := os.Create(outfile)
	//if err != nil {
	//return err
	//}
	//defer out.Close()
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Flight Statistics")
	if err != nil {
		return err
	}
	// flights - header
	// 1st line
	h1 := sheet.AddRow()
	h1.AddCell().SetString("Date")
	to := h1.AddCell()
	to.Merge(1, 0)
	to.SetString("Takeoff")
	h1.AddCell()
	la := h1.AddCell()
	la.Merge(1, 0)
	la.SetString("Landing")
	h1.AddCell()
	h1.AddCell().SetString("Duration")
	h1.AddCell().SetString("Filename")
	// 2nd line
	h2 := sheet.AddRow()
	h2.AddCell()
	h2.AddCell().SetString("Time")
	h2.AddCell().SetString("Site")
	h2.AddCell().SetString("Time")
	h2.AddCell().SetString("Site")
	// flights - data
	for _, r := range *fdata {
		row := sheet.AddRow()
		for _, c := range r {
			cell := row.AddCell()
			cell.Value = c
		}
	}
	// statistics - header
	h := sheet.AddRow()
	h.AddCell().SetString("Period")
	h.AddCell().SetString("Flights")
	h.AddCell().SetString("Duration")
	// statistics - data
	for _, r := range *sdata {
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
