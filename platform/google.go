package platform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/otiai10/hotsub/params"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
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
func (p GoogleCloudPlatform) Validate(ctx params.Context) error {

	client, err := google.DefaultClient(context.Background())
	if err != nil {
		return err
	}
	service, err := compute.New(client)
	if err != nil {
		return err
	}

	if err := p.createGoogleFirewallIfNotExists(service, ctx); err != nil {
		return err
	}

	return nil
}

func (p GoogleCloudPlatform) errorizeOperationError(ope *compute.OperationError) error {
	if ope == nil {
		return nil
	}
	buf := new(bytes.Buffer)
	json.NewDecoder(buf).Decode(ope)
	return fmt.Errorf(buf.String())
}
