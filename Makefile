GROUP := github.com/TheoBrigitte
NAME := confluence

DOCKER_IMAGE := theo01/${NAME}

VERSION := $(shell git describe --always --long --dirty --tags || date)

all: server

# Server

server: package-server publish-server

build-server:
	@go install -v -ldflags '-s -w' ./cmd/confluence

package-server:
	@docker build -t ${DOCKER_IMAGE} -f Dockerfile.server .

publish-server:
	@docker push ${DOCKER_IMAGE}

run-server:
	@go run ./cmd/confluence -addr=0.0.0.0:7803
	#@docker run --rm -p 7803:7803 -p 50007:50007 theo01/confluence:latest -addr=0.0.0.0:7803

client: package-client publish-client

build-client:
	cd client && $(MAKE) $@

package-client:
	cd client && $(MAKE) $@

publish-client:
	cd client && $(MAKE) $@

run-client:
	cd client && $(MAKE) $@

.PHONY: server build-server package-server publish-server run-server
