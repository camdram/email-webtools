GO=$(shell which go)
GOGET=$(GO) get
GOMOD=$(GO) mod
GOFMT=$(GO) fmt
GOBUILD=$(GO) build

export GOARCH=amd64
export GOOS=linux

dir:
	@if [ ! -d bin ] ; then mkdir -p bin ; fi

get:
	@$(GOGET) github.com/shuLhan/go-bindata/...

mod:
	@$(GOMOD) download

format:
	$(GOFMT) ./...

build/assets: get
	go-bindata -o internal/assets/assets.go -pkg assets assets/

build: build/assets dir mod
	$(GOBUILD) -o bin/email-webtools main.go

clean:
	@rm -rf bin
	@rm -f internal/assets/assets.go

assets: build/assets

all: format build
