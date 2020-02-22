# maker
<p>
  <a href="https://github.com/njuettner/maker/releases/latest"><img alt="GitHub release" src="https://img.shields.io/github/v/release/njuettner/maker.svg?logo=github&style=flat-square"></a>
  <a href="https://github.com/njuettner/maker/actions?workflow=goreleaser"><img alt="Release workflow" src="https://github.com/njuettner/maker/workflows/goreleaser/badge.svg"></a>
</p>

Creates Makefile for Go projects

## Usage

```console
make help
```

```console
Usage:

  build          builds a local binary
  build-linux    builds binary for linux/amd64
  build-darwin   builds binary for darwin/amd64
  install        install the application
  run            runs go run main.go
  clean          cleans the binary
  lint           runs golangci-lint
  test           runs go test with default values
  setup          setup go modules
  build-docker   builds docker image to registry
  push-docker    pushes docker image to registry
  help           prints this help message
```

## Example Makefile:

```make
APPLICATION   ?= example
REGISTRY      ?= njuettner
VERSION       ?= $(shell git describe --tags --always --dirty)
SOURCES       = $(shell find . -name '*.go')
GOPKGS        = $(shell go list ./...)
TAG           ?= $(VERSION)
BUILD_FLAGS   ?= -v
LDFLAGS       ?= -X main.version=$(VERSION) -w -s

default: build

.PHONY: build
## build: builds a local binary
build: clean $(SOURCES)
	CGO_ENABLED=0 go build -o ${APPLICATION} $(BUILD_FLAGS) .

.PHONY: build-linux
## build-linux: builds binary for linux/amd64
build-linux: clean $(SOURCES)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) .

.PHONY: build-darwin
## build-darwin: builds binary for darwin/amd64
build-darwin: clean $(SOURCES)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) .

.PHONY: install
## install: install the application
install: $(SOURCES)
	go install .

.PHONY: run
## run: runs go run main.go
run: $(SOURCES)
	go run -race .

.PHONY: clean
## clean: cleans the binary
clean:
	go clean

.PHONY: lint
## lint: runs golangci-lint
lint:
	golangci-lint run --timeout=10m ./...

.PHONY: test
## test: runs go test with default values
test:
	go test -v -race -cover $(GOPKGS)

.PHONY: setup
## setup: setup go modules
setup:
	@go mod init \
		&& go mod tidy \
		&& go mod vendor

.PHONY: build-docker
## build-docker: builds docker image to registry
build-docker: build-linux
	docker build -t ${APPLICATION}:${TAG} .

.PHONY: push-docker
## push-docker: pushes docker image to registry
push-docker: build-docker
	docker push ${REGISTRY}/${APPLICATION}:${TAG}

.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
```
