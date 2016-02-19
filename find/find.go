package find

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
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
			fmt.Printf("processing: %s\n", path)
			flight, err := igc.NewFlight(path)
			if err != nil {
				return err
			}
			if _, err := os.Stat(fmt.Sprintf("%s.glider", path)); err == nil {
				if g, err := ioutil.ReadFile(fmt.Sprintf("%s.glider", path)); err == nil {
					flight.Glider = strings.Trim(string(g), " \t\r\n")
				}
			}
			flights.Add(flight)
		}
		// manually added flights
		if strings.HasSuffix(f.Name(), "_manual.csv") {
			fmt.Printf("processing manual entries from: %s\n", path)
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
				if len(l) < 7 {
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
					Glider:      l[5],
					Comment:     l[6],
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
