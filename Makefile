BINARY = webhook-docker
GOARCH = amd64

VERSION?=$(shell git describe --tags --always --dirty --match=* 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo v0)
SHORT_VERSION?=$(shell git describe --tags --always --dirty --match=* | cut -d'-' -f1 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo v0)
MAIN_BUILD_PATH=./cmd/webhook-docker
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
DATE    ?= $(shell date +%FT%T%z)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags '-w -s -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.branch=${BRANCH} -X main.buildDate=${DATE} -X main.shortVersion=${SHORT_VERSION} -X container.buildDate=${DATE} -X container.shortVersion=${SHORT_VERSION}  -X main.build="production"'


.PHONY: build
build:
	@echo "Start: GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-linux-${GOARCH} $(MAIN_BUILD_PATH)"; \
	GO111MODULE=on GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-linux-${GOARCH}  $(MAIN_BUILD_PATH) ; \
	echo "End build"



.PHONY: docker_install
docker_install: build
	@cp -f ${BINARY}-linux-${GOARCH} ${GOPATH}/bin/${BINARY}



.PHONY: version
version:
	@echo "Current dir: $(CURDIR)"
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Date: $(DATE)"

