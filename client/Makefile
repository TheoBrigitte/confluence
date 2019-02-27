NAME := confluence
DOCKER_IMAGE := theo01/${NAME}-client

VERSION := $(shell git describe --always --long --dirty --tags || date)

all: client

client: package-client publish-client

build-client:
	npm run build

package-client:
	@docker build -t ${DOCKER_IMAGE} .

publish-client:
	@docker push ${DOCKER_IMAGE}

run-client:
	npm start

.PHONY: client build-client package-client publish-client run-client
