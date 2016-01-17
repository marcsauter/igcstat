// Copyright ©2016 NAME HERE <EMAIL ADDRESS>
//
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/marcsauter/flightstat"
	"github.com/marcsauter/igcstat/find"
	"github.com/spf13/cobra"
)

var csvfile string
var stdout bool

// csvCmd respresents the csv command
var csvCmd = &cobra.Command{
	Use:   "csv",
	Short: "output in csv format",
	Run: func(cmd *cobra.Command, args []string) {
		flights := find.Flights(dir)
		stat, err := flightstat.NewFlightStat(flights)
		if err != nil {
			log.Fatal(err)
		}
		// write file
		f, err := os.Create(csvfile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		out := io.MultiWriter(f)
		if stdout {
			out = io.MultiWriter(f, os.Stdout)
		}
		w := csv.NewWriter(out)
		flights.Csv(w)
		stat.Csv(w)
		w.Flush()
	},
}

func init() {
	RootCmd.AddCommand(csvCmd)
	csvCmd.Flags().StringVarP(&csvfile, "file", "f", fmt.Sprintf("%s.csv", os.Args[0]), "output filename")
	csvCmd.Flags().BoolVar(&stdout, "stdout", false, "write to stdout")
}
