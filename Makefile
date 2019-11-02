GO=$(shell which go)
GOGET=$(GO) get
GOFMT=$(GO) fmt
GOBUILD=$(GO) build

export GOARCH=amd64
export GOOS=linux

dir:
	@if [ ! -d $(CURDIR)/bin ] ; then mkdir -p $(CURDIR)/bin ; fi

get:
	@$(GOGET) github.com/joho/godotenv
	@$(GOGET) github.com/go-sql-driver/mysql

format:
	$(GOFMT) main.go
	$(GOFMT) internal/client/client.go
	$(GOFMT) internal/server/controller.go internal/server/driver.go internal/server/server.go

build:
	$(GOBUILD) -o $(CURDIR)/bin/email-webtools main.go

clean:
	@rm -rf $(CURDIR)/bin

all: dir get format build
