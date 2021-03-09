PWD:=$(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

GOBIN=$(PWD)/build/bin
export GOBIN
GOPROXY?=https://proxy.golang.org
export GOPROXY
GOSUMDB?=sum.golang.org
export GOSUMDB
COVERAGE=$(PWD)/build/coverage

ifeq ($(OS),Windows_NT)
	GO="$(shell where go)"
	PATH:=$(GOBIN);$(PATH)
	BINEXT:=".exe"
else
	GO=$(shell which go)
	PATH:=$(GOBIN):$(PATH)
endif
export PATH

BIN_SERVER=powwow-server
BIN_CLIENT=powwow-client

VERSION=$(shell git describe --tags --always | sed 's/^v//g')

LDFLAGS=-ldflags "-X main.version=$(VERSION)"

.DEFAULT_GOAL: all
.ONESHELL:
.PHONY: all goenv generate tidy test clean clean-cache compose $(BIN_SERVER) $(BIN_CLIENT)

all: $(BIN_SERVER) $(BIN_CLIENT)

goenv:
	$(GO) env

generate:
	$(GO) generate ./...

tidy:
	$(GO) mod tidy

test:
	mkdir -p "$(COVERAGE)"
	$(GO) test -v ./pkg/... -coverprofile "$(COVERAGE)"/test.out

clean:
	rm -rf "$(GOBIN)"
	rm -rf "$(COVERAGE)"

clean-cache:
	$(GO) clean -cache ./...

compose:
	docker-compose up --build

$(BIN_SERVER):
	$(GO) install -v $(LDFLAGS) ./cmd/$(BIN_SERVER)

$(BIN_CLIENT):
	$(GO) install -v $(LDFLAGS) ./cmd/$(BIN_CLIENT)
