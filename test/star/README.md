```sh
awsub \
  --tasks ./test/star/star-alignment-tasks.single.csv \
  --script ./test/star/star-alignment.sh \
  --image friend1ws/star-alignment \
  --aws-iam-instance-profile awsubtest \
  --aws-ec2-instance-type m4.2xlarge \
  --disk-size 128 \
  --verbose
```
