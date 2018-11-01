package platform

// Provider ...
type Provider string

const (
	AWS    Provider = "aws"
	GCP    Provider = "gcp"
	VBOX   Provider = "vbox"
	HYPERV Provider = "hyperv"
)

// Driver ...
type Driver string

const (
	AmazonEC2 Driver = "amazonec2"
	Google    Driver = "google"
	Vbox      Driver = "virtualbox"
	Hyperv    Driver = "hyperv"
)

const (
	// HotsubSecurityStructureVersion ...
	HotsubSecurityStructureVersion = "2018-08-06"
)
