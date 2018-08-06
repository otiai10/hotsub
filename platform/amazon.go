package platform

import (
	"fmt"

	"github.com/otiai10/hotsub/params"

	"github.com/otiai10/iamutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
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

func createWorkspaceVpcIfNotExists(sess *session.Session, ctx params.Context) error {

	client := ec2.New(sess)

	vpcsout, err := client.DescribeVpcs(&ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{Name: aws.String("tag:Name"), Values: aws.StringSlice([]string{DefaultAWSWorkspaceVpcName})},
		},
	})
	if err != nil {
		return err
	}

	// Good, existing.
	if len(vpcsout.Vpcs) != 0 {
		return ctx.Set("aws-vpc-id", *vpcsout.Vpcs[0].VpcId)
	}

	// Create because it doesn't exist
	createout, err := client.CreateVpc(&ec2.CreateVpcInput{
		CidrBlock: aws.String("10.0.0.0/16"),
	})
	if err != nil {
		return err
	}
	vpc := createout.Vpc

	// Name this VPC
	_, err = client.CreateTags(&ec2.CreateTagsInput{
		Tags:      []*ec2.Tag{{Key: aws.String("Name"), Value: aws.String(DefaultAWSWorkspaceVpcName)}},
		Resources: []*string{vpc.VpcId},
	})
	if err != nil {
		return err
	}

	// Create subnet
	subnetout, err := client.CreateSubnet(&ec2.CreateSubnetInput{
		AvailabilityZone: aws.String(ctx.String("aws-region") + "a"), // FIXME: hard coded
		CidrBlock:        aws.String("10.0.0.0/18"),
		VpcId:            vpc.VpcId,
	})
	if err != nil {
		return err
	}
	subnet := subnetout.Subnet

	// Name this subnet
	_, err = client.CreateTags(&ec2.CreateTagsInput{
		Tags:      []*ec2.Tag{{Key: aws.String("Name"), Value: aws.String(DefaultAWSSubnetName)}},
		Resources: []*string{subnet.SubnetId},
	})
	if err != nil {
		return err
	}

	// Create InternetGateway
	gatewayout, err := client.CreateInternetGateway(&ec2.CreateInternetGatewayInput{})
	if err != nil {
		return err
	}
	gateway := gatewayout.InternetGateway

	// Name this internet gateway
	_, err = client.CreateTags(&ec2.CreateTagsInput{
		Tags:      []*ec2.Tag{{Key: aws.String("Name"), Value: aws.String(DefaultAWSInternetGatewayName)}},
		Resources: []*string{gateway.InternetGatewayId},
	})
	if err != nil {
		return err
	}

	// Attach this InternetGateway to the VPC
	_, err = client.AttachInternetGateway(&ec2.AttachInternetGatewayInput{
		InternetGatewayId: gateway.InternetGatewayId,
		VpcId:             vpc.VpcId,
	})
	if err != nil {
		return err
	}

	// Create RouteTable
	routetablesout, err := client.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{Name: aws.String("vpc-id"), Values: []*string{vpc.VpcId}},
		},
	})
	if err != nil {
		return err
	}
	if len(routetablesout.RouteTables) == 0 {
		return fmt.Errorf("no default routetable found on this VPC")
	}
	routetable := routetablesout.RouteTables[0]

	// Name this default RouteTable
	_, err = client.CreateTags(&ec2.CreateTagsInput{
		Tags:      []*ec2.Tag{{Key: aws.String("Name"), Value: aws.String(DefaultAWSRouteTableName)}},
		Resources: []*string{routetable.RouteTableId},
	})
	if err != nil {
		return err
	}

	_, err = client.CreateRoute(&ec2.CreateRouteInput{
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		RouteTableId:         routetable.RouteTableId,
		GatewayId:            gateway.InternetGatewayId,
	})
	if err != nil {
		return err
	}

	return ctx.Set("aws-vpc-id", *vpc.VpcId)
}

func createSecurityGroupIfNotExists(sess *session.Session, ctx params.Context) error {
	client := ec2.New(sess)

	// Check existing SecurityGroup
	name := DefaultAWSSecurityGroupName
	res, err := client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		GroupNames: []*string{&name},
	})

	// If exists, well done!
	if err == nil && len(res.SecurityGroups) != 0 {
		return nil
	}

	// If any error other than "NotFound", mark it error
	if ae, ok := err.(awserr.Error); ok {
		if ae.Code() != "InvalidGroup.NotFound" {
			return err
		}
	} else {
		return err
	}

	// It seems not existing. Let's create new one.
	group, err := client.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		GroupName:   &name,
		VpcId:       aws.String(ctx.String("aws-vpc-id")),
		Description: aws.String("default sg of hotsub"),
	})
	if err != nil {
		return err
	}

	if _, err := client.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{group.GroupId},
		Tags:      []*ec2.Tag{{Key: aws.String("Name"), Value: aws.String("hotsub-default")}},
	}); err != nil {
		return err
	}

	anywhere := []*ec2.IpRange{&ec2.IpRange{CidrIp: aws.String("0.0.0.0/0")}}
	tcp := "tcp"

	_, err = client.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: group.GroupId,
		IpPermissions: []*ec2.IpPermission{
			&ec2.IpPermission{IpProtocol: &tcp, IpRanges: anywhere, FromPort: aws.Int64(22), ToPort: aws.Int64(22)},
			&ec2.IpPermission{IpProtocol: &tcp, IpRanges: anywhere, FromPort: aws.Int64(2376), ToPort: aws.Int64(2376)},
			&ec2.IpPermission{
				IpProtocol: &tcp, FromPort: aws.Int64(2049), ToPort: aws.Int64(2049),
				UserIdGroupPairs: []*ec2.UserIdGroupPair{{GroupId: group.GroupId}},
			},
		},
	})

	// If failed to authorize, we shouldn't keep it alive.
	if err != nil {
		_, derr := client.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{GroupId: group.GroupId})
		return fmt.Errorf("SecurityGroup failed to be authorized: %v\nAlso failed to be deleted: %v\nPlease delete this security group `%s` manually", err, derr, name)
	}

	return nil
}

func createInstanceProfileIfNotExists(sess *session.Session, ctx params.Context) error {

	_, err := iamutil.FindInstanceProfile(sess, DefaultAWSInstanceProfileNameForCompute)
	if err == nil {
		// Found! Well done! Do nothing!
		return nil
	}

	ae, ok := err.(awserr.Error)
	if !ok || ae.Code() != iam.ErrCodeNoSuchEntityException {
		return err
	}

	// Well, it's not found. Create new one.
	profile := &iamutil.InstanceProfile{
		Name: DefaultAWSInstanceProfileNameForCompute,
		Role: &iamutil.Role{
			Description: "hotsub Instance Profile for computing nodes",
			PolicyArns: []string{
				"arn:aws:iam::aws:policy/AmazonS3FullAccess",
			},
		},
	}

	return profile.Create(sess)
}
