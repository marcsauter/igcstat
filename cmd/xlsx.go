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
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/marcsauter/flightstat"
	"github.com/marcsauter/igcstat/find"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// xlsxCmd respresents the xlsx command
var xlsxCmd = &cobra.Command{
	Use:   "xlsx",
	Short: "output in xslx format",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("writing output to %s ... ", viper.GetString("xlsxfile"))
		flights := find.Flights(viper.GetString("srcpath"))
		stat, err := flightstat.NewFlightStat(flights, viper.GetString("glider"))
		if err != nil {
			log.Fatal(err)
		}
		flightstat.Xlsx(flights, stat, viper.GetString("xlsxfile"))
		fmt.Println("DONE")
	},
}

func init() {
	RootCmd.AddCommand(xlsxCmd)
	l := len(os.Args[0]) - len(filepath.Ext(os.Args[0])) // len of filename without extension
	viper.SetDefault("xlsxfile", fmt.Sprintf("%s.xlsx", os.Args[0][0:l]))
	xlsxCmd.Flags().StringP("file", "f", viper.GetString("xlsxfile"), "xlsx output filename")
	viper.BindPFlag("xlsxfile", xlsxCmd.Flags().Lookup("file"))
}
