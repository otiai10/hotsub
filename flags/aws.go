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

// AwsRegion ...
var AwsRegion = cli.StringFlag{
	Name:  "aws-region",
	Usage: `AWS region name in which AmazonEC2 instances would be launched`,
	Value: "ap-northeast-1",
}

// AwsEC2InstanceType ...
var AwsEC2InstanceType = cli.StringFlag{
	Name:  "aws-ec2-instance-type",
	Usage: `AWS EC2 instance type. If specified, all --min-cores and --min-ram would be ignored.`,
	Value: "t2.micro",
}

// AwsIAMInstanceProfile ...
var AwsIAMInstanceProfile = cli.StringFlag{
	Name:  "aws-iam-instance-profile",
	Usage: `AWS instance profile from your IAM roles.`,
}
