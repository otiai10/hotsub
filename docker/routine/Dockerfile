# This docker image is supposed to be used as hotsub/routine, in which
# input files are donwloaded onto the host machine and
# output files are uploaded to cloud storage services.
FROM debian:stretch-slim

LABEL maintainer="Hiromu Ochiai<otiai10@gmail.com>"

RUN apt-get update -qq
RUN apt-get install -y -qq \
    python \
    python-pip \
    python-dev \
    build-essential \
    groff \
    less \
    man \
    ssh \
    curl

# For AWS cli
RUN pip install pip --upgrade --quiet
RUN pip install awscli --upgrade --quiet

# For gcloud, gsutil
RUN curl https://sdk.cloud.google.com | bash
ENV PATH=${PATH}:/root/google-cloud-sdk/bin

# Downloader/Uploader scripts
COPY ./scripts /scripts
