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
	// AwsubSecurityStructureVersion ...
	AwsubSecurityStructureVersion = "2018-06-07"
)
