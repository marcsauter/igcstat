package find

import (
	"encoding/csv"
	"fmt"
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
		// manually added flights
		if strings.HasSuffix(f.Name(), "_manual.csv") {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			r := csv.NewReader(f)
			r.Comment = '#'
			d, err := r.ReadAll()
			if err != nil {
				return err
			}
			for _, l := range d {
				if len(l) < 6 {
					log.Printf("%s: wrong number of columns", path)
					continue
				}
				date, err := time.Parse("02.01.2006", l[0])
				if err != nil {
					log.Printf("%s: unknown date format '%s'", path, l[0])
					log.Println(err)
					continue
				}
				takeoffTime, err := time.Parse("02.01.2006 15:04", fmt.Sprintf("%s %s", l[0], l[1]))
				if err != nil {
					log.Printf("%s: unknown time format '%s'", path, l[1])
					continue
				}
				landingTime, err := time.Parse("02.01.2006 15:04", fmt.Sprintf("%s %s", l[0], l[3]))
				if err != nil {
					log.Printf("%s: unknown time format '%s'", path, l[3])
					continue
				}
				flight := &igc.Flight{
					Date: date,
					TakeOff: igc.Fix{
						Time: takeoffTime,
					},
					TakeOffSite: l[2],
					Landing: igc.Fix{
						Time: landingTime,
					},
					LandingSite: l[4],
					Duration:    landingTime.Sub(takeoffTime),
					Comment:     l[5],
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
