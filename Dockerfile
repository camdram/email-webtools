FROM golang:latest

RUN apt-get update -qq
RUN apt-get install -y git make

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

WORKDIR /go/src/github.com/camdram/email-webtools
CMD make clean && make all
