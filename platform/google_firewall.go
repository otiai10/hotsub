package platform

import (
	"github.com/otiai10/hotsub/params"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func (p GoogleCloudPlatform) createGoogleFirewallIfNotExists(service *compute.Service, ctx params.Context) error {

	project := ctx.String("google-project")

	firewall, err := service.Firewalls.Get(project, DefaultGoogleFirewallName).Do()
	if err != nil {
		apierror, ok := err.(*googleapi.Error)
		if !ok {
			return err
		}
		if apierror.Code != 404 {
			return err
		}
	}
	if firewall != nil {
		return nil
	}

	// It was 404, create new one

	firewall = &compute.Firewall{
		Name:        DefaultGoogleFirewallName,
		Description: "Allow NFS connection for ExTL Shared Data Instance model",
		Allowed:     []*compute.FirewallAllowed{{IPProtocol: "tcp", Ports: []string{"2049"}}},
		Direction:   "INGRESS",
		TargetTags:  []string{DefaultGoogleInstanceTag},
		SourceTags:  []string{DefaultGoogleInstanceTag},
	}

	op, err := service.Firewalls.Insert(project, firewall).Do()
	if err != nil {
		return err
	}
	if op.Error != nil {
		return p.errorizeOperationError(op.Error)
	}

	return nil
}
