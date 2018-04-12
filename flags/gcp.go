package flags

import "github.com/urfave/cli"

/////////////////////////////////////
// Flags for Google Cloud Platform //
/////////////////////////////////////

// googleProject ...
var googleProject = cli.StringFlag{
	Name:  "google-project",
	Usage: "Project ID for GCP",
}

// googleZone ...
var googleZone = cli.StringFlag{
	Name:  "google-zone",
	Usage: "GCP service zone name",
	Value: "asia-northeast1-a",
}
