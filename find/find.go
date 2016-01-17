package find

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/marcsauter/igc"
)

func Flights(dir string) *igc.Flights {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatal(err)
	}
	flights := igc.NewFlights()
	//
	evaluate := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}
		if strings.HasSuffix(f.Name(), ".igc") {
			flight, err := igc.NewFlight(path)
			if err != nil {
				return err
			}
			flights.Add(flight)
		}
		if strings.HasSuffix(f.Name(), "_manual.csv") {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			r := csv.NewReader(f)
			d, err := r.ReadAll()
			if err != nil {
				return err
			}
			for _, l := range d {
				if len(l) < 6 {
					log.Printf("%s: wrong number of columns", path)
					continue
				}
				date, err := time.Parse(l[0], "02.01.2006")
				if err != nil {
					log.Printf("%s: unknown date format '%s'", path, l[0])
					continue
				}
				takeoffTime, err := time.Parse(l[1], "15:04")
				if err != nil {
					log.Printf("%s: unknown time format '%s'", path, l[1])
					continue
				}
				landingTime, err := time.Parse(l[3], "15:04")
				if err != nil {
					log.Printf("%s: unknown time format '%s'", path, l[3])
					continue
				}
				duration, err := time.ParseDuration(l[5])
				if err != nil {
					log.Printf("%s: unknown duration format '%s'", path, l[5])
					continue
				}
				flight := &igc.Flight{
					Date:        date,
					TakeOff:     takeoffTime,
					TakeOffSite: l[2],
					Landing:     landingTime,
					LandingSite: l[4],
					Duration:    duration,
				}
				flights.Add(flight)
			}
		}
		return nil
	}
	//
	if err := filepath.Walk(dir, evaluate); err != nil {
		log.Print(err)
	}
	sort.Sort(flights)
	return flights
}
