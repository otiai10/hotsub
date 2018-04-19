#!/bin/bash

PROJROOT=$(dirname $(dirname $(cd $(dirname $0) && pwd)))

set -e -v
awsub \
    --tasks ${PROJROOT}/test/wordcount/wordcount.csv \
    --script ${PROJROOT}/test/wordcount/main.sh \
    --env FOO="This is foo, specified as a common env" \
    --aws-iam-instance-profile awsubtest \
    --log-dir /tmp \
    --verbose

# TODO: The command above should be specified with following JSON/YAML
#       in the future ;)
#
# {
#     "task": [
#         ["--env NAME", "--input INPUT01", "--input INPUT02"],
#         ["otiai10", "s3://foo/bar/baz.txt", "s3://foo/bar/qak.txt"],
#     ],
#     "script": "./test/wordcount/main.sh",
#     "image": "debian:latest",
#     "env": {
#         "FOO": "This is foo, specified as a common env"
#     },
#     "platform": "aws",
#     "aws": {
#         "InstanceProfile": "awsubtest",
#         "InstanceType": "t2.micro"
#     }
# }