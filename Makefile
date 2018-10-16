GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install

GOPATH:=$(GOPATH)

APP_NAME=dnsupdater

PACKAGE_NAME=github.com/boris1993/$(APP_NAME)

all: windows-amd64 darwin-amd64 linux-amd64 mips-softfloat

help:
	@echo Usage: make \<TARGET\>
	@echo TARGET could be windows-amd64, darwin-amd64, linux-amd64, mips-softfloat
	@echo All 4 targets will be built if target is not specified

clean:
	rm -rf $(GOPATH)/bin/$(APP_NAME)/$(APP_NAME)

get-dep:
	go get

windows-amd64:
	get-dep
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(GOPATH)/bin/$(APP_NAME)-windows-amd64/$(APP_NAME) -i -v $(PACKAGE_NAME)
	cp config.yaml.template $(GOPATH)/bin/$(APP_NAME)-windows-amd64/

darwin-amd64:
	get-dep
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(GOPATH)/bin/$(APP_NAME)-darwin-amd64/$(APP_NAME) -i -v $(PACKAGE_NAME)
	cp config.yaml.template $(GOPATH)/bin/$(APP_NAME)-darwin-amd64/

linux-amd64:
	get-dep
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(GOPATH)/bin/$(APP_NAME)-linux-amd64/$(APP_NAME) -i -v $(PACKAGE_NAME)
	cp config.yaml.template $(GOPATH)/bin/$(APP_NAME)-linux-amd64/

mips-softfloat:
	get-dep
	GOARCH=mips GOOS=linux GOMIPS=softfloat $(GOBUILD) -o $(GOPATH)/bin/$(APP_NAME)-linux-mips-softfloat/$(APP_NAME) -i -v $(PACKAGE_NAME)
	cp config.yaml.template $(GOPATH)/bin/$(APP_NAME)-linux-mips-softfloat/

