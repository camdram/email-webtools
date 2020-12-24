FROM golang:latest

RUN apt-get update -qq
RUN apt-get install -y git make

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

WORKDIR /app
CMD make clean && make all
