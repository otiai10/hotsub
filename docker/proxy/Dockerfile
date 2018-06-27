# This docker image specifies how to get started with `awsub` CLI.
FROM debian:stretch-slim

# 0) Packages
RUN apt-get update -qq
RUN apt-get install -y curl wget python-pip sudo

# 1) docker-machine
# See https://docs.docker.com/machine/install-machine/#install-machine-directly
RUN base=https://github.com/docker/machine/releases/download/v0.14.0 \
  && curl -L $base/docker-machine-$(uname -s)-$(uname -m) >/tmp/docker-machine \
  && sudo install /tmp/docker-machine /usr/local/bin/docker-machine

# 2) aws cli
RUN pip install awscli

# 3) awsub itself
RUN wget -q https://github.com/otiai10/awsub/releases/download/v0.3.0/awsub.linux_amd64.tar.gz
RUN tar -xzvf awsub.linux_amd64.tar.gz && mv ./awsub /usr/local/bin

# Usage of this image:
#
#     docker run -it -v ~/.aws:/root/.aws awsub/proxy
#