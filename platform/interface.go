package platform

import (
	"github.com/otiai10/hotsub/params"
)

// Platform ...
type Platform interface {
	Validate(ctx params.Context) error
}

// Get ...
func Get(ctx params.Context) Platform {
	switch Provider(ctx.String("provider")) {
	case AWS:
		return &AmazonWebServices{Region: ctx.String("aws-region")}
	case GCP:
		return &GoogleCloudPlatform{}
	case VBOX:
		return &Virtualbox{}
	case HYPERV:
		return &HyperV{}
	}
	return &AmazonWebServices{}
}
