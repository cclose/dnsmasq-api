# Detect the Docker executable
DOCKER_EXEC ?= $(shell \
	if command -v podman >/dev/null 2>&1; then echo podman; \
	elif command -v nerdctl >/dev/null 2>&1; then echo nerdctl; \
	elif command -v docker >/dev/null 2>&1; then echo docker; \
	el echo "error: No suitable container tool found" >&2; exit 1; \
	fi)
DOCKER := $(DOCKER_EXEC)
ifeq ($(VERSION),)
VERSION := $(shell git describe --tags --always)
endif
ifeq ($(COMMIT),)
COMMIT := $(shell git rev-parse --short HEAD)
endif
ifeq ($(BUILD_TIME),)
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
endif
ifeq ($(TAG),)
TAG := "latest"
endif

setup:
	cp config.yaml.dist config.yaml
	cp dnsmasq.conf.dist dnsmasq.conf

lint:
	golangci-lint run

vet:
	go vet ./...

build:
	go build -ldflags "-X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT)' -X 'main.BuildTimeStr=$(BUILD_TIME)'" -o dnsMasqAPI main.go

docker:
	$(DOCKER) build . -f Dockerfile \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg VERSION=$(VERSION) \
		-t dnsmasqapi:$(TAG)

run: docker
	$(DOCKER) run -d -p 8080:8080 -v "dnsmasq.conf:/etc/dnsmasq.conf" -v "config.yaml:/app/config.yaml" --name dnsapi localhost/dnsmasqapi:latest
	curl localhost:8080/statusz

logs:
	$(DOCKER) logs -f dnsapi

clean:
	$(DOCKER) stop dnsapi
	$(DOCKER) rm dnsapi


.PHONY: setup lint vet build docker run logs clean
