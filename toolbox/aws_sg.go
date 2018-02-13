/*Package toolbox is an utility functions for each providers.*/
package toolbox

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// CreateSecurityGroupIfNotExists ...
func CreateSecurityGroupIfNotExists(name, region string) error {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: &region},
	}))
	client := ec2.New(sess)

	res, err := client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		GroupNames: []*string{&name},
	})

	// If exists, pass
	if err == nil && len(res.SecurityGroups) != 0 {
		return nil
	}

	// If any error other than "NotFound", mark it error
	if err != nil {
		ae, ok := err.(awserr.Error)
		if !ok {
			return err
		}
		if ae.Code() != "InvalidGroup.NotFound" {
			return err
		}
	}

	group, err := client.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		GroupName:   &name,
		Description: aws.String("default sg of awsub"),
	})
	if err != nil {
		return err
	}

	anywhere := []*ec2.IpRange{&ec2.IpRange{CidrIp: aws.String("0.0.0.0/0")}}
	tcp := "tcp"

	_, err = client.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: group.GroupId,
		IpPermissions: []*ec2.IpPermission{
			&ec2.IpPermission{IpProtocol: &tcp, IpRanges: anywhere, FromPort: aws.Int64(22), ToPort: aws.Int64(22)},
			&ec2.IpPermission{IpProtocol: &tcp, IpRanges: anywhere, FromPort: aws.Int64(2376), ToPort: aws.Int64(2376)},
		},
	})

	// If failed to authorize, we shouldn't keep it alive.
	if err != nil {
		_, derr := client.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{GroupId: group.GroupId})
		return fmt.Errorf("SecurityGroup failed to be authorized: %v\nAlso failed to be deleted: %v\nPlease delete this security group `%s` manually", err, derr, name)
	}

	return nil
}
