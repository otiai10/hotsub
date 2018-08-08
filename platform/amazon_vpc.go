package platform

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/otiai10/hotsub/params"
)

func createWorkspaceVpcIfNotExists(sess *session.Session, ctx params.Context) error {

	client := ec2.New(sess)

	vpcsout, err := client.DescribeVpcs(&ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{{Name: aws.String("tag:Name"), Values: aws.StringSlice([]string{DefaultAWSWorkspaceVpcName})}},
	})
	if err != nil {
		return err
	}

	// Good, existing.
	if len(vpcsout.Vpcs) != 0 {
		err = ctx.Set("aws-vpc-id", *vpcsout.Vpcs[0].VpcId)
		if err != nil {
			return err
		}
		subnetout, err := client.DescribeSubnets(&ec2.DescribeSubnetsInput{
			Filters: []*ec2.Filter{{Name: aws.String("tag:Name"), Values: aws.StringSlice([]string{DefaultAWSSubnetName})}},
		})
		if err != nil {
			return err
		}
		if len(subnetout.Subnets) == 0 {
			return fmt.Errorf("no subnets found, clean up Workspace VPC beforehand")
		}
		err = ctx.Set("aws-subnet-id", *subnetout.Subnets[0].SubnetId)
		if err != nil {
			return err
		}
		return nil
	}

	// Create because it doesn't exist
	createout, err := client.CreateVpc(&ec2.CreateVpcInput{
		CidrBlock: aws.String(privateIPAddressRange + "/16"),
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
		CidrBlock:        aws.String(privateIPAddressRange + "/18"),
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

	if err := ctx.Set("aws-subnet-id", *subnet.SubnetId); err != nil {
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
