SHELL := bash
.ONESHELL:

VER=$(shell git describe --tags)
GO=$(shell which go)
GOGET=$(GO) get
GOMOD=$(GO) mod
GOFMT=$(GO) fmt
GOBUILD=$(GO) build -mod=readonly -ldflags "-X main.version=$(VER)"

dir:
	@if [ ! -d bin ]; then mkdir -p bin; fi

mod:
	$(GOMOD) download

format:
	$(GOFMT) ./...

build/assets:
	$(GOGET) github.com/shuLhan/go-bindata/...
	go-bindata -o internal/assets/assets.go -pkg assets assets/
	
	
build/linux/amd64: dir mod build/assets
	export CGO_ENABLED=0
	export GOOS=linux
	export GOARCH=amd64
	$(GOBUILD) -o bin/email-webtools-linux-$(VER:v%=%)-amd64 main.go

build/linux: build/linux/amd64

build: build/linux

clean:
	@rm -rf bin
	@rm -f internal/assets/assets.go

assets: build/assets

all: format build
