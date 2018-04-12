package flags

import "github.com/urfave/cli"

//////////////////////////////////
// Flags for Amazon Web Service //
//////////////////////////////////

// AwsVPCFlag ...
// var AwsVPCFlag = cli.StringFlag{
// 	Name:  "aws-vpc",
// 	Usage: `AWS VPC ID in which AmazonEC2 instances would be launched`,
// }

var awsRegion = cli.StringFlag{
	Name:  "aws-region",
	Usage: `AWS region name in which AmazonEC2 instances would be launched`,
	Value: "ap-northeast-1",
}

// awsEC2InstanceType ...
var awsEC2InstanceType = cli.StringFlag{
	Name:  "aws-ec2-instance-type",
	Usage: `AWS EC2 instance type. If specified, all --min-cores and --min-ram would be ignored.`,
	Value: "t2.micro",
}

// awsIAMInstanceProfile ...
var awsIAMInstanceProfile = cli.StringFlag{
	Name:  "aws-iam-instance-profile",
	Usage: `AWS instance profile from your IAM roles.`,
}
