package flags

import "github.com/urfave/cli"

/////////////////////////////////////
// Flags for Google Cloud Platform //
/////////////////////////////////////

// GoogleProject ...
var GoogleProject = cli.StringFlag{
	Name:  "google-project",
	Usage: "Project ID for GCP",
}

// GoogleZone ...
var GoogleZone = cli.StringFlag{
	Name:  "google-zone",
	Usage: "GCP service zone name",
	Value: "asia-northeast1-a",
}
