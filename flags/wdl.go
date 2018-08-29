package flags

import "github.com/urfave/cli"

// WDLFileFlag specifies WDL file.
var WDLFileFlag = cli.StringFlag{
	Name:  "wdl",
	Usage: "WDL file to run your workflow",
}

// WDLJobFileFlag represents parameter files of WDL
var WDLJobFileFlag = cli.StringSliceFlag{
	Name:  "wdl-job",
	Usage: "Parameter files for WDL",
}
