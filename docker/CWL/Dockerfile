FROM ubuntu:16.04

# For cwltool
RUN apt-get update -qq
RUN apt-get install -y \
  python-dev \
  python-pip \
  libxml2-dev \
  libxslt1-dev \
  zlib1g-dev

RUN pip install --upgrade \
  mypy-extensions==0.3.0 \
  requests==2.19.1 \
  six \
  cwltool==1.0.20180820141117

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