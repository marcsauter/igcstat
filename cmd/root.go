package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile                    string
	dir                        string
	takeoffSites, landingSites string
	distance                   int
)

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "igcstat",
	Short: "summary and statistics over a collection of igc files",
	Long: `collect all igc files in a given directory and generate a list of flights
and statistics over all flights either in xlsx or in csv format`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

//Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.igcstat.yaml)")
	d, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	RootCmd.PersistentFlags().StringVarP(&dir, "dir", "d", d, "start directory")
	RootCmd.PersistentFlags().StringVar(&takeoffSites, "takeoff", "Waypoints_Startplatz.gpx", "takeoff sites")
	RootCmd.PersistentFlags().StringVar(&landingSites, "landing", "Waypoints_Landeplatz.gpx", "landing sites")
	RootCmd.PersistentFlags().IntVar(&distance, "distance", 300, "maximal distance to the nearest known site")
}

// Read in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".igcstat") // name of config file (without extension)
	viper.AddConfigPath("$HOME")    // adding home directory as first search path
	viper.AutomaticEnv()            // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
