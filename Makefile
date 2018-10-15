GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install

GOPATH:=$(GOPATH)

PACKAGE_NAME=dnsupdater

all: mips-softfloat

help:
	@echo Usage: make \<TARGET\>
	@echo TARGET could be windows-amd64, darwin-amd64, linux-amd64, mips-softfloat
	@echo Default target is mips-softfloat

windows-amd64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(GOPATH)/bin/$(PACKAGE_NAME)_windows_amd64 -i -v $(PACKAGE_NAME)

darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(GOPATH)/bin/$(PACKAGE_NAME)_darwin_amd64 -i -v $(PACKAGE_NAME)

linux-amd64:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(GOPATH)/bin/$(PACKAGE_NAME)_linux_amd64 -i -v $(PACKAGE_NAME)

mips-softfloat:
	GOARCH=mips GOOS=linux GOMIPS=softfloat $(GOBUILD) -o $(GOPATH)/bin/$(PACKAGE_NAME)_linux_mips_softfloat -i -v $(PACKAGE_NAME)
