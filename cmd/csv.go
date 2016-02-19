// Copyright ©2016 Marc Sauter <marc.sauter@bluewin.ch>
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
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/marcsauter/flightstat"
	"github.com/marcsauter/igcstat/find"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// csvCmd respresents the csv command
var csvCmd = &cobra.Command{
	Use:   "csv",
	Short: "output in csv format",
	Run: func(cmd *cobra.Command, args []string) {
		flights := find.Flights(viper.GetString("srcpath"))
		stat, err := flightstat.NewFlightStat(flights, viper.GetString("glider"))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("writing output to %s ... ", viper.GetString("csvfile"))
		flightstat.Csv(flights, stat, viper.GetString("csvfile"))
		fmt.Println("DONE")
	},
}

func init() {
	RootCmd.AddCommand(csvCmd)
	l := len(os.Args[0]) - len(filepath.Ext(os.Args[0])) // len of filename without extension
	viper.SetDefault("csvfile", fmt.Sprintf("%s.csv", os.Args[0][0:l]))
	csvCmd.Flags().StringP("file", "f", viper.GetString("csvfile"), "csv output filename")
	viper.BindPFlag("csvfile", csvCmd.Flags().Lookup("file"))

}
