GOCMD=go
GOINSTALL=$(GOCMD) install

GOPATH:=$(GOPATH):$(CURDIR)

PACKAGE_NAME=dnsupdater

all: mips-softfloat

help:
	@echo Usage: make \<TARGET\>
	@echo TARGET could be windows-amd64, darwin-amd64, linux-amd64, mips-softfloat
	@echo Default target is mips-softfloat

windows-amd64:
	GOARCH=amd64 GOOS=windows $(GOINSTALL) $(PACKAGE_NAME)

darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(GOINSTALL) $(PACKAGE_NAME)

linux-amd64:
	GOARCH=amd64 GOOS=linux $(GOINSTALL) $(PACKAGE_NAME)

mips-softfloat:
	GOARCH=mips GOOS=linux GOMIPS=softfloat $(GOINSTALL) $(PACKAGE_NAME)
