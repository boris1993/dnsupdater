GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=%(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BINARY_NAME=dnsupdater

all: build build-mips-softfloat

build:
    $(GOBUILD) -o $(BINARY_NAME) -v

build-mips-softfloat:
	GOARCH=mips GOOS=linux GOMIPS=softfloat $(GOBUILD) -o $(BINARY_NAME) -v
