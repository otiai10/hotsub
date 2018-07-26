package platform

// Provider ...
type Provider string

const (
	AWS   Provider = "aws"
	GCP   Provider = "gcp"
	Local Provider = "local"
)

// Driver ...
type Driver string

const (
	AmazonEC2 Driver = "amazonec2"
	Google    Driver = "google"
	Vbox      Driver = "virtualbox"
)

const (
	// HotsubSecurityStructureVersion ...
	HotsubSecurityStructureVersion = "2018-06-07"
)
