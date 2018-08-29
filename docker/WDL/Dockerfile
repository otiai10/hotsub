FROM ubuntu:16.04

# Install cromwell
RUN apt-get update -qq
RUN apt-get install -y \
  wget \
  default-jre

RUN wget -q https://github.com/broadinstitute/cromwell/releases/download/34/cromwell-34.jar

# For Docker client-binary and daemon
RUN apt-get install -y \
  apt-transport-https \
  software-properties-common \
  lsb-release \
  gnupg \
  curl

RUN curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add -
RUN add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
RUN apt-get update -qq
RUN apt-get install -y docker-ce

RUN apt-get clean

VOLUME /var/lib/docker
ENTRYPOINT ["dockerd"]