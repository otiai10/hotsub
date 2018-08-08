package platform

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/otiai10/hotsub/params"
	"github.com/otiai10/iamutil"
)

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
