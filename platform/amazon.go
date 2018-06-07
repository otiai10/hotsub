package platform

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	// DefaultAWSSecurityGroupName default aws security group name
	DefaultAWSSecurityGroupName = "awsub-default"
)

// AmazonWebServices ...
type AmazonWebServices struct {
	Region            string
	SecurityGroupName string
}

// Validate validates the platform itself, setting up the infrastructures if needed.
// For AWS, it executes:
//     1. Create SecurityGroup for awsub.
func (p *AmazonWebServices) Validate() error {

	// Initialize EC2 API Client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: &p.Region},
	}))

	if err := createSecurityGroupIfNotExists(sess); err != nil {
		return fmt.Errorf("failed to setup security group: %v", err)
	}

	return nil
}

func createSecurityGroupIfNotExists(sess *session.Session) error {
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
		Description: aws.String("default sg of awsub"),
	})
	if err != nil {
		return err
	}

	if _, err := client.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{group.GroupId},
		Tags:      []*ec2.Tag{{Key: aws.String("Name"), Value: aws.String("awsub-default")}},
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
