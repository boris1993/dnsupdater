GOCMD=go
GOINSTALL=$(GOCMD) install

GOPATH:=$(GOPATH):$(CURDIR)

PACKAGE_NAME=dnsupdater

all: build-mips-softfloat

build-windows-amd64:
	GOARCH=amd64 GOOS=windows $(GOINSTALL) $(PACKAGE_NAME)

build-darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(GOINSTALL) $(PACKAGE_NAME)

build-linux-amd64:
	GOARCH=amd64 GOOS=linux $(GOINSTALL) $(PACKAGE_NAME)

build-mips-softfloat:
	GOARCH=mips GOOS=linux GOMIPS=softfloat $(GOINSTALL) $(PACKAGE_NAME)
