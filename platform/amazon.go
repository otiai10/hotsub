package platform

import (
	"fmt"

	"github.com/otiai10/hotsub/params"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (

	// DefaultAWSWorkspaceVpcName ...
	DefaultAWSWorkspaceVpcName = "hotsub-workspace" + "-" + HotsubSecurityStructureVersion

	// DefaultAWSSubnetName ...
	DefaultAWSSubnetName = "hotsub-subnet" + "-" + HotsubSecurityStructureVersion

	// DefaultAWSInternetGatewayName ...
	DefaultAWSInternetGatewayName = "hotsub-gateway" + "-" + HotsubSecurityStructureVersion

	// DefaultAWSRouteTableName ...
	DefaultAWSRouteTableName = "hotsub-default-rt" + "-" + HotsubSecurityStructureVersion

	// DefaultAWSInstanceProfileNameForCompute default aws instance profile name
	DefaultAWSInstanceProfileNameForCompute = "hotsub-compute" + "-" + HotsubSecurityStructureVersion
	// TODO: Separate instance profile for shared data instance

	// DefaultAWSSecurityGroupName default aws security group name
	DefaultAWSSecurityGroupName = "hotsub-default" + "-" + HotsubSecurityStructureVersion

	privateIPAddressRange = "172.31.0.0"
)

// AmazonWebServices ...
type AmazonWebServices struct {
	Region            string
	SecurityGroupName string
}

// Validate validates the platform itself, setting up the infrastructures if needed.
// For AWS, it executes:
//     1. Create SecurityGroup for hotsub.
func (p *AmazonWebServices) Validate(ctx params.Context) error {

	// Initialize EC2 API Client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: &p.Region},
	}))

	if err := createWorkspaceVpcIfNotExists(sess, ctx); err != nil {
		return fmt.Errorf("failed to setup VPC: %v", err)
	}

	if err := createSecurityGroupIfNotExists(sess, ctx); err != nil {
		return fmt.Errorf("failed to setup security group: %v", err)
	}

	if err := createInstanceProfileIfNotExists(sess, ctx); err != nil {
		return fmt.Errorf("failed to setup instance profile: %v", err)
	}

	return nil
}
