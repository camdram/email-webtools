GO=$(shell which go)
GOGET=$(GO) get
GOBUILD=$(GO) build

export GOARCH=amd64
export GOOS=linux

dir:
	@if [ ! -d $(CURDIR)/pkg ] ; then mkdir -p $(CURDIR)/pkg ; fi

get:
	@$(GOGET) github.com/joho/godotenv
	@$(GOGET) github.com/go-sql-driver/mysql

build:
	$(GOBUILD) -o $(CURDIR)/pkg/email-webtools main.go controller.go driver.go

clean:
	@rm -rf $(CURDIR)/pkg

all: dir get build
