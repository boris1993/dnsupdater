GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install

GOPATH:=$(GOPATH)

APP_NAME=dnsupdater

PACKAGE_NAME=github.com/boris1993/$(APP_NAME)

all: get-dep mips-softfloat

help:
	@echo Usage: make \<TARGET\>
	@echo TARGET could be windows-amd64, darwin-amd64, linux-amd64, mips-softfloat
	@echo Default target is mips-softfloat

clean:
	rm -rf $(GOPATH)/bin/$(APP_NAME)/$(APP_NAME)

get-dep:
	go get

windows-amd64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(GOPATH)/bin/$(APP_NAME)/$(APP_NAME) -i -v $(PACKAGE_NAME)
	cp config.yaml.template $(GOPATH)/bin/$(APP_NAME)/

darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(GOPATH)/bin/$(APP_NAME)/$(APP_NAME) -i -v $(PACKAGE_NAME)
	cp config.yaml.template $(GOPATH)/bin/$(APP_NAME)/

linux-amd64:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(GOPATH)/bin/$(APP_NAME)/$(APP_NAME) -i -v $(PACKAGE_NAME)
	cp config.yaml.template $(GOPATH)/bin/$(APP_NAME)/

mips-softfloat:
	GOARCH=mips GOOS=linux GOMIPS=softfloat $(GOBUILD) -o $(GOPATH)/bin/$(APP_NAME)/$(APP_NAME) -i -v $(PACKAGE_NAME)
	cp config.yaml.template $(GOPATH)/bin/$(APP_NAME)/
