GO=$(shell which go)
GOGET=$(GO) get
GOFMT=$(GO) fmt
GOBUILD=$(GO) build

export GOARCH=amd64
export GOOS=linux

dir:
	@if [ ! -d $(CURDIR)/pkg ] ; then mkdir -p $(CURDIR)/pkg ; fi

get:
	@$(GOGET) github.com/joho/godotenv
	@$(GOGET) github.com/go-sql-driver/mysql

format:
	$(GOFMT) main.go controller.go driver.go
	$(GOFMT) client.go

build:
	$(GOBUILD) -o $(CURDIR)/pkg/email-webtools-server main.go controller.go driver.go
	$(GOBUILD) -o $(CURDIR)/pkg/email-webtools-client client.go

clean:
	@rm -rf $(CURDIR)/pkg

all: dir get format build
