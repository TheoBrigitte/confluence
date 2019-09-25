GROUP := github.com/TheoBrigitte
NAME := confluence

DOCKER_IMAGE := theo01/${NAME}:latest

VERSION := $(shell git describe --always --long --dirty --tags || date)

all: build

build:
	@go install -v -ldflags '-s -w' ./cmd/confluence

docker-image:
	@docker build -t ${DOCKER_IMAGE} .

docker-push:
	@docker push ${DOCKER_IMAGE}

run:
	@go run ./cmd/confluence -addr=0.0.0.0:7803
	#@docker run --rm -p 7803:7803 -p 50007:50007 theo01/confluence:latest -addr=0.0.0.0:7803

systemd.unit:
	@docker run \
		--rm \
		-i \
		-v $(PWD)/systemd:/tmp/systemd \
		theo01/template \
		-template /tmp/systemd/confluence.service.tmpl \
		-out /tmp/systemd/confluence.service \
		<<< '{"Version":"$(VERSION)"}'

integration-test:
include /etc/confluence/opensubtitles_credentials
export
	@go test ./... \
	    -tags=integration \
	    -osUser=$(OPENSUBTITLES_USER) \
	    -osPassword=$(OPENSUBTITLES_PASSWORD) \
	    -osUserAgent=$(OPENSUBTITLES_USERAGENT)

.PHONY: build docker-image docker-push run
