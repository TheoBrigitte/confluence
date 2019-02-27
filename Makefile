GROUP := github.com/TheoBrigitte
NAME := confluence

DOCKER_IMAGE := theo01/${NAME}

VERSION := $(shell git describe --always --long --dirty --tags || date)

all: server

# Server

server: package publish

build:
	@go install -v -ldflags '-s -w' ./cmd/confluence

package:
	@docker build -t ${DOCKER_IMAGE} .

publish:
	@docker push ${DOCKER_IMAGE}

run:
	@go run ./cmd/confluence -addr=0.0.0.0:7803
	#@docker run --rm -p 7803:7803 -p 50007:50007 theo01/confluence:latest -addr=0.0.0.0:7803

.PHONY: server build package publish run
