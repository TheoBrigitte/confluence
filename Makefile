GROUP := github.com/TheoBrigitte
NAME := confluence

DOCKER_IMAGE := theo01/${NAME}
PKG := ${GROUP}/${NAME}

OPERATOR_BIN := build/operator
OPERATOR_PKG := ${PKG}/cmd/operator

PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

VERSION := $(shell git describe --always --long --dirty || date)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)

all: server client

# Server

server: build-server publish-server

build-server:
	@docker build -t ${DOCKER_IMAGE} -f Dockerfile .

publish-server:
	@docker push ${DOCKER_IMAGE}

run-server:
	@ go run ./cmd/confluence -addr=0.0.0.0:7803

	#@docker run --rm -p 7803:7803 -p 50007:50007 theo01/confluence:latest -addr=0.0.0.0:7803


# Client

client: build-client

build-client:
	cd client && npm run build

run-client:
	cd client && npm start

.PHONY: server build-server publish-server run-server client build-client run-client
