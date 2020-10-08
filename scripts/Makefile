ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))/..))
GOCMD=go
GOBUILD=$(GOCMD) build

BUILD_ARGS=-i

APP_NAME=dnsupdater

all: windows-amd64 darwin-amd64 linux-amd64 mips-softfloat
.PHONY: all

help:
	@echo Usage: make \<TARGET\>
	@echo TARGET could be windows-amd64, darwin-amd64, linux-amd64, mips-softfloat
	@echo All 4 targets will be built if target is not specified

windows-amd64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) $(BUILD_ARGS) -o $(ROOT_DIR)/bin/$(APP_NAME)-windows-amd64/$(APP_NAME).exe $(ROOT_DIR)/cmd/dnsupdater.go
	cp config.yaml.template $(ROOT_DIR)/bin/$(APP_NAME)-windows-amd64/

darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) $(BUILD_ARGS) -o $(ROOT_DIR)/bin/$(APP_NAME)-darwin-amd64/$(APP_NAME) $(ROOT_DIR)/cmd/dnsupdater.go
	cp config.yaml.template $(ROOT_DIR)/bin/$(APP_NAME)-darwin-amd64/

linux-amd64:
	GOARCH=amd64 GOOS=linux $(GOBUILD) $(BUILD_ARGS) -o $(ROOT_DIR)/bin/$(APP_NAME)-linux-amd64/$(APP_NAME) $(ROOT_DIR)/cmd/dnsupdater.go
	cp config.yaml.template $(ROOT_DIR)/bin/$(APP_NAME)-linux-amd64/

mips-softfloat:
	GOARCH=mips GOOS=linux GOMIPS=softfloat $(GOBUILD) $(BUILD_ARGS) -o $(ROOT_DIR)/bin/$(APP_NAME)-linux-mips-softfloat/$(APP_NAME) $(ROOT_DIR)/cmd/dnsupdater.go
	cp config.yaml.template $(ROOT_DIR)/bin/$(APP_NAME)-linux-mips-softfloat/

.PHONY: clean
clean:
	rm -r $(ROOT_DIR)/bin
