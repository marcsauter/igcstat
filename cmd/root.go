package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/marcsauter/igc"
	"github.com/marcsauter/wpt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "igcstat",
	Short: "summary and statistics over a collection of igc files",
	Long: `collect all igc files in a given directory and generate a list of flights
and statistics over all flights either in xlsx or in csv format`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Run(xlsxCmd, []string{})
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

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/igcstat.yaml)")

	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	RootCmd.Flags().StringP("srcpath", "s", path, "start directory")
	viper.BindPFlag("srcpath", RootCmd.Flags().Lookup("srcpath"))

	RootCmd.Flags().StringP("takeoff", "t", "Waypoints_Startplatz.gpx", "takeoff sites")
	viper.BindPFlag("takeoff", RootCmd.Flags().Lookup("takeoff"))

	RootCmd.Flags().StringP("landing", "l", "Waypoints_Landeplatz.gpx", "landing sites")
	viper.BindPFlag("landing", RootCmd.Flags().Lookup("landing"))

	RootCmd.Flags().IntP("distance", "d", 300, "maximal distance to the nearest known site")
	viper.BindPFlag("distance", RootCmd.Flags().Lookup("distance"))
}

// Read in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	viper.SetConfigType("yaml")
	viper.SetConfigName("igcstat")                  // name of config file (without extension)
	viper.AddConfigPath("$HOME")                    // adding home directory as first search path
	viper.AddConfigPath(path)                       // adding source directory as second search path
	viper.AddConfigPath(viper.GetString("srcpath")) // adding source directory as second search path
	viper.AutomaticEnv()                            // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("using config file:", viper.ConfigFileUsed())
	}

	//
	igc.MaxDistance = viper.GetInt("distance")
	fmt.Printf("searching within a radius of %dm for known takeoff/landing sites\n", viper.GetInt("distance"))

	//
	if to, err := wpt.NewWaypoints(viper.GetString("takeoff")); err == nil {
		fmt.Printf("using %s for takeoff site lookup\n", viper.GetString("takeoff"))
		igc.RegisterTakeoffSiteSource(to)
	}

	//
	if la, err := wpt.NewWaypoints(viper.GetString("landing")); err == nil {
		fmt.Printf("using %s for landing site lookup\n", viper.GetString("landing"))
		igc.RegisterLandingSiteSource(la)
	}

	fmt.Printf("starting search for flights in %s\n", viper.GetString("srcpath"))
}
