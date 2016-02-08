# igcstat
`igcstat` recursively search through a directory for IGC files and creates a flight statistics.
The default output format is _xlsx_ and the default directory is the current working directory.

## Installation
Download the appropriate binary:
* [Linux (amd64)](dist/igcstat-linux-amd64.tar?raw=true)
* [OSX (amd64)](dist/igcstat-osx-amd64.tar?raw=true)
* [Windows (amd64)](dist/igcstat-windows-amd64.zip?raw=true)

### Windows
* Extract `igcstat-windows-amd64.zip` in the parent directory of your IGC files.
* Execute `igcstat.cmd` (double-click in explorer)
* Open `igcstat.xlsx`

### Linux 
* Extract the archive in the parent directory of your IGC files: `tar -xf igcstat-linux-amd64.tar` 
* Execute `igcstat.sh`
* Open `igcstat.xlsx`

### OSX 
* Extract the archive in the parent directory of your IGC files: `tar -xf igcstat-osx-amd64.tar` (or just double-click in Finder) 
* Execute `igcstat.sh` (first run: right click - Open With - Terminal)
* Open `igcstat.xlsx`

### Takeoff and landing sites
The waypoint files (GPX) with the known takeoff and landing sites are available on flyland.ch. There filenames defaults to:
* Waypoints_Landeplatz.gpx 
* Waypoints_Startplatz.gpx

... and will be searched in the current working directory. Other filenames can be defined on the command line.

The maximal distance to an official takeoff or landing site is 300m and can be adjusted on the command line.

__If no waypoint file are available or no official takeoff or landing site can be found, the Google Maps API will be used.__

###

## Add flights manually
You can add flights manually to the statistics by adding CSV formatted files somewhere below your starting directory.
The name of the CSV file has to match the pattern `<name>_manual.csv` (e.g. `addflights_manual.csv`).

### Fields
Comment lines start with `#` and will be ignored.
* Date:        dd.mm.yyyy
* TakeOffTime: hh.mm
* TakeOffSite: e.g. Amisbuehl
* LandingTime: hh:mm
* LandingSite: e.g. Lehn
* Comment:     e.g. "vario malfunction"

### Example
    # Date,TakeOffTime,TakeOffSite,LandingTime,LandingSite,Comment
    "18.04.2015","19:15","Amisbuehl","19:30","Lehn","no vario"
