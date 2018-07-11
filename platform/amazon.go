package platform

import (
	"fmt"

	"github.com/otiai10/iamutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
)

const (
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
func (p *AmazonWebServices) Validate() error {

	// Initialize EC2 API Client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: &p.Region},
	}))

	if err := createSecurityGroupIfNotExists(sess); err != nil {
		return fmt.Errorf("failed to setup security group: %v", err)
	}

	if err := createInstanceProfileIfNotExists(sess); err != nil {
		return fmt.Errorf("failed to setup instance profile: %v", err)
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

func createInstanceProfileIfNotExists(sess *session.Session) error {

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
