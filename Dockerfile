FROM golang:latest

RUN apt-get update -qq
RUN apt-get install -y git build-essential

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN go get github.com/joho/godotenv
RUN go get -u github.com/go-sql-driver/mysql

WORKDIR /go/src/github.com/camdram/email-webtools
CMD [ "make", "all" ]
