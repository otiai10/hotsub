package platform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/otiai10/hotsub/params"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

const (
	// DefaultGoogleInstanceTag ....
	DefaultGoogleInstanceTag = "hotsub" + "-" + HotsubSecurityStructureVersion

	// DefaultGoogleFirewallName ...
	DefaultGoogleFirewallName = "hotsub-firewall" + "-" + HotsubSecurityStructureVersion
)

// GoogleCloudPlatform ...
type GoogleCloudPlatform struct{}

// Validate validates Google platform itself.
// 1) Create Firewall Rule for specific tag.
func (gcp GoogleCloudPlatform) Validate(ctx params.Context) error {

	project := ctx.String("google-project")

	client, err := google.DefaultClient(context.Background())
	if err != nil {
		return err
	}
	service, err := compute.New(client)
	if err != nil {
		return err
	}

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
		return gcp.errorizeOperationError(op.Error)
	}

	return nil
}

func (gcp GoogleCloudPlatform) errorizeOperationError(ope *compute.OperationError) error {
	if ope == nil {
		return nil
	}
	buf := new(bytes.Buffer)
	json.NewDecoder(buf).Decode(ope)
	return fmt.Errorf(buf.String())
}
