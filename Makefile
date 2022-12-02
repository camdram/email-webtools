SHELL := bash
.ONESHELL:

VER=$(shell git describe --tags --always --dirty)
GO=$(shell which go)
GOGET=$(GO) get
GOMOD=$(GO) mod
GOFMT=$(GO) fmt
GOBUILD=$(GO) build -trimpath -mod=readonly -ldflags "-X main.version=$(VER:v%=%) -s -w -buildid="

dir:
	@if [ ! -d bin ]; then mkdir -p bin; fi

mod:
	$(GOMOD) download

format:
	$(GOFMT) ./...

build/linux/amd64: dir mod
	export CGO_ENABLED=0
	export GOOS=linux
	export GOARCH=amd64
	$(GOBUILD) -o bin/email-webtools-linux-$(VER:v%=%)-amd64 main.go

build/linux: build/linux/amd64

build: build/linux

clean:
	@rm -rf bin

all: format build
