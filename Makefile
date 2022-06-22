DOCKER_HOSTNAME=ghcr.io
DOCKER_TAGNAME=latest
DOCKER_NAMESPACE=roytman
DOCKER_NAME=perf-helm

IMG ?= ${DOCKER_HOSTNAME}/${DOCKER_NAMESPACE}/${DOCKER_NAME}:${DOCKER_TAGNAME}

.PHONY: all
all: docker-build docker-push

.PHONY: source-build
source-build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o perf-helm main.go

# Overwrite docker-build from docker.mk
.PHONY: docker-build
docker-build: source-build
	docker build . -t ${IMG} -f Dockerfile

.PHONY: docker-push
docker-push:
ifneq (${DOCKER_PASSWORD},)
	@docker login \
		--username ${DOCKER_USERNAME} \
		--password ${DOCKER_PASSWORD} ${DOCKER_HOSTNAME}
endif
	docker push ${IMG}
