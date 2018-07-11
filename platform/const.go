package platform

// Provider ...
type Provider string

const (
	AWS Provider = "aws"
	GCP Provider = "gcp"
)

// Driver ...
type Driver string

const (
	AmazonEC2 Driver = "amazonec2"
	Google    Driver = "google"
)

const (
	// HotsubSecurityStructureVersion ...
	HotsubSecurityStructureVersion = "2018-06-07"
)
