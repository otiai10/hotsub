# STAR RNA-seq alignment example

This example specifies how to run `STAR` on AWS by using `awsub`.

```sh
awsub \
  --tasks ./star-alignment.single.csv \
  --script ./main.sh \
  --image friend1ws/star-alignment \
  --aws-iam-instance-profile awsubtest \
  --aws-ec2-instance-type m4.2xlarge \
  --disk-size 128 \
  --verbose
```

# What you need beforehand are

1. CSV file (tasks file) which indicates URLs for your samples
2. Reference files' location, as well
3. IAM instance profile that you have (TODO: It's gonna be NOT required in future version)

# Working example

See `./runner.sh` for more details.