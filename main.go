package main

import (
	"fmt"
	"os"
	"text/template"

	log "github.com/sirupsen/logrus"
)

var makefile = `APPLICATION   ?= {{.Application}}
REGISTRY      ?= {{.Registry}}
VERSION       ?= $(shell git describe --tags --always --dirty)
TAG           ?= $(VERSION)
BUILD_FLAGS   ?= -v
LDFLAGS       ?= -X main.version=$(VERSION) -w -s

.PHONY: build
## build: build the application
build: clean
	@echo "Building..."
	@go build -o ${APPLICATION} main.go

.PHONY: build
## install: install the application
install:
	@echo "Installing..."
	@go install ./...

.PHONY: run
## run: runs go run main.go
run:
	go run -race main.go

.PHONY: clean
## clean: cleans the binary
clean:
	@echo "Cleaning"
	@go clean

.PHONY: lint
## lint: runs golangci-lint
lint:
	golangci-lint run --timeout=5m ./...

.PHONY: test
## test: runs go test with default values
test:
	go test -v -count=1 -race ./...

.PHONY: setup
## setup: setup go modules
setup:
	@go mod init \
		&& go mod tidy \
		&& go mod vendor

.PHONY: docker-build
## docker-build: builds docker image to registry
docker-build: build
	docker build -t ${APPLICATION}:${TAG} .

.PHONY: docker-push
## docker-push: pushes docker image to registry
docker-push: docker-build
	docker push ${REGISTRY}/${APPLICATION}:${TAG}

.PHONY: help
## help: Prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
`

type Config struct {
	Application string
	Registry    string
}

func main() {
	file, err := os.Create("Makefile")
	if err != nil {
		log.Fatalf("Unable to create 'Makefile', %s", err)
	}
	config := &Config{}
	makefileTemplate := template.New("Makefile")
	t := template.Must(makefileTemplate.Parse(makefile))
	fmt.Println("Please enter the name of your application:")
	fmt.Fscanln(os.Stdin, &config.Application)
	fmt.Println("Please enter the url of registry, e.g. docker user:")
	fmt.Fscanln(os.Stdin, &config.Registry)
	t.Execute(file, config)
}
