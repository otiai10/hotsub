package flags

import "github.com/urfave/cli"

// CWLFileFlag specifies CWL file.
var CWLFileFlag = cli.StringFlag{
	Name:  "cwl",
	Usage: "CWL file to run your workflow",
}

// CWLParamFlag represents parameter files of CWL
var CWLParamFlag = cli.StringSliceFlag{
	Name:  "cwl-param",
	Usage: "Parameter files for CWL",
}

// IncludeFlag represents local files to be included onto workflow container.
var IncludeFlag = cli.StringSliceFlag{
	Name:  "include",
	Usage: "Local files to be included onto workflow container",
}
