package platform

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/otiai10/hotsub/params"
)

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
