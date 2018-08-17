# This docker image specifies how to get started with `hotsub` CLI.
FROM debian:stretch-slim

# 0) Packages
RUN apt-get update -qq
RUN apt-get install -y curl wget python-pip sudo groff vim

# 1) docker-machine
# See https://docs.docker.com/machine/install-machine/#install-machine-directly
RUN base=https://github.com/docker/machine/releases/download/v0.14.0 \
  && curl -L $base/docker-machine-$(uname -s)-$(uname -m) >/tmp/docker-machine \
  && sudo install /tmp/docker-machine /usr/local/bin/docker-machine

# 2) aws cli
RUN pip install awscli

ARG HOTSUB_VERSION=v0.5.0
# 3) hotsub itself
RUN wget -q https://github.com/otiai10/hotsub/releases/download/${HOTSUB_VERSION}/hotsub.linux_amd64.tar.gz
RUN tar -xzvf hotsub.linux_amd64.tar.gz && mv ./hotsub /usr/local/bin

# Usage of this image:
#
#     docker run -it hotsub/proxy
#
# Then you can use aws, docker-machine and hotsub commands available