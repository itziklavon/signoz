-include .env

VERSION := 1.0.0
BUILD := $(BUILD_ID)
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME = $(shell date +'%Y-%m-%dT%H:%M:%S')
PROJECTNAME := goapm
ARTIFACTNAME := ${PROJECTNAME}-${VERSION}-${BUILD}

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD) -X=main.Commit=$(COMMIT) -X=main.BuildTime=$(BUILD_TIME)"

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## install: Install missing dependencies. Runs `go get` internally. e.g; make install get=github.com/foo/bar
install: go mod tidy

go-artifact-name:
	@echo "${ARTIFACTNAME}"

go-build:
	@echo "  >  build time - ${BUILD_TIME}"
	@echo "  >  Syncing Dependencies..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go mod tidy
	@echo "  >  Building binary..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build $(LDFLAGS) -o ${ARTIFACTNAME}

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo