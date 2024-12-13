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

setup-env:
	# Create ENV file for testing release scripts
	mkdir -p dist
	echo "ARCHIVE_NAME=\"dnsMasqAPI-v0.0.0-test-linux-amd64\"" > dist/envvar.sh
	echo "ARTIFACT_DIR=\"dist/artifacts\"" >> dist/envvar.sh
	echo "ARCHIVE_PATH=\"dist/artifacts/dnsMasqAPI-v0.0.0-test-linux-amd64.tar.gz\"" >> dist/envvar.sh
	echo "GITHUB_REF=\"v0.0.0-test\"" >> dist/envvar.sh
	echo "PLATFORM=\"linux\"" >> dist/envvar.sh
	echo "ARCH=\"amd64\"" >> dist/envvar.sh

lint:
	golangci-lint run

vet:
	go vet ./...

test:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.out ./...

build:
	go build -ldflags "-X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT)' -X 'main.BuildTimeStr=$(BUILD_TIME)'" -o dnsMasqAPI main.go

docker:
	$(DOCKER) build . -f Dockerfile \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg VERSION=$(VERSION) \
		-t dnsmasqapi:$(TAG)

run: docker
	$(DOCKER) run -d -p 8080:8080 -v "dnsmasq.conf:/etc/dnsmasq.d/api.conf" -v "config.yaml:/app/config.yaml" --name dnsapi localhost/dnsmasqapi:latest
	curl localhost:8080/statusz

logs:
	$(DOCKER) logs -f dnsapi

clean:
	$(DOCKER) stop dnsapi
	$(DOCKER) rm dnsapi


.PHONY: setup lint vet test coverage build docker run logs clean
