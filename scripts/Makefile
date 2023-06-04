ROOT_DIR:=$(realpath $(dir $(realpath $(lastword $(MAKEFILE_LIST))))../)
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOTEST_ARGS=-coverprofile=coverage.out

BUILD_ARGS=

APP_NAME=dnsupdater

all: windows-amd64 darwin-amd64 linux-amd64 mips-softfloat
.PHONY: all

help:
	@echo Usage: make \<TARGET\>
	@echo TARGET could be windows-amd64, darwin-amd64, linux-amd64, mips-softfloat
	@echo All 4 targets will be built if target is not specified

test:
	$(GOTEST) $(GOTEST_ARGS) $(ROOT_DIR)/internal/...

windows-amd64:
	GOARCH=amd64 GOOS=windows CGO_ENABLED=0 $(GOBUILD) $(BUILD_ARGS) -o $(ROOT_DIR)/bin/$(APP_NAME)-windows-amd64/$(APP_NAME).exe $(ROOT_DIR)/cmd/dnsupdater/main.go
	cp $(ROOT_DIR)/configs/config.yaml.template $(ROOT_DIR)/bin/$(APP_NAME)-windows-amd64/

darwin-amd64:
	GOARCH=amd64 GOOS=darwin CGO_ENABLED=0 $(GOBUILD) $(BUILD_ARGS) -o $(ROOT_DIR)/bin/$(APP_NAME)-darwin-amd64/$(APP_NAME) $(ROOT_DIR)/cmd/dnsupdater/main.go
	cp $(ROOT_DIR)/configs/config.yaml.template $(ROOT_DIR)/bin/$(APP_NAME)-darwin-amd64/

linux-amd64:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 $(GOBUILD) $(BUILD_ARGS) -o $(ROOT_DIR)/bin/$(APP_NAME)-linux-amd64/$(APP_NAME) $(ROOT_DIR)/cmd/dnsupdater/main.go
	cp $(ROOT_DIR)/configs/config.yaml.template $(ROOT_DIR)/bin/$(APP_NAME)-linux-amd64/

mips-softfloat:
	GOARCH=mips GOOS=linux GOMIPS=softfloat CGO_ENABLED=0 $(GOBUILD) $(BUILD_ARGS) -o $(ROOT_DIR)/bin/$(APP_NAME)-linux-mips-softfloat/$(APP_NAME) $(ROOT_DIR)/cmd/dnsupdater/main.go
	cp $(ROOT_DIR)/configs/config.yaml.template $(ROOT_DIR)/bin/$(APP_NAME)-linux-mips-softfloat/

.PHONY: clean
clean:
	rm -r $(ROOT_DIR)/bin
