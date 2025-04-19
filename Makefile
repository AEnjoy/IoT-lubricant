PROJECT_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

VERSION := $(shell git describe --tags --exact-match 2>/dev/null || git rev-parse --abbrev-ref HEAD)
BUILD_TIME := $(shell date +"%Y-%m-%d %H:%M:%S")
GIT_COMMIT := $(shell git rev-parse --short HEAD)
GO_VERSION := $(shell go version | awk '{print $$3}')
FEATURES := $(or $(ENV_LUBRICANT_ENABLE_FEATURES),$(shell git rev-parse --abbrev-ref HEAD))
BUILD_HOST_PLATFORM := $(shell uname -s | tr '[:upper:]' '[:lower:]')/$(shell uname -m)
ifeq ($(shell uname -s),Linux)
PLATFORM_VERSION := $(shell grep -E '^(NAME|VERSION)=' /etc/os-release | tr -d '"' | awk -F= '{print $$2}' | paste -sd ' ' -)
else ifeq ($(shell uname -s),Windows)
PLATFORM_VERSION := $(shell systeminfo | findstr /B /C:"OS Name" /C:"OS Version" | awk -F: '{print $$2}' | paste -sd ' ' -)
else
PLATFORM_VERSION := unknown
endif

CGO_ENABLED ?= 0
CGO_COMPONENTS := datastore

GO_TAGS := -tags=sonic,avx
LD_FLAGS = -w -s \
    -X 'main.Version=$(VERSION)' \
    -X 'main.BuildTime=$(BUILD_TIME)' \
    -X 'main.GoVersion=$(GO_VERSION)' \
    -X 'main.GitCommit=$(GIT_COMMIT)' \
    -X 'main.Features=$(FEATURES)' \
    -X 'main.BuildHostPlatform=$(BUILD_HOST_PLATFORM)' \
    -X 'main.PlatformVersion=$(PLATFORM_VERSION)'

COMPONENTS := gateway apiserver agent logg grpcserver reporter # datastore

.PHONY: all test test-coverage install mock docker-build clean help make-output-dir

all: build-all

make-output-dir: clean
	@mkdir -p ./bin

cgo-init:
ifeq ($(CGO_ENABLED),1)
	@echo "Initializing CGO environment for $(COMPONENT)..."
	@if [ -f /etc/redhat-release ]; then \
		sudo yum install -y gcc glibc-devel; \
	elif [ -f /etc/debian_version ]; then \
		sudo apt-get update && sudo apt-get install -y build-essential; \
	elif [ "$(shell uname -s)" = "Darwin" ]; then \
		xcode-select --install || true; \
	fi
#	@if [ "$(COMPONENT)" = "gateway" ]; then \
#		echo "Installing SQLite dependencies..."; \
#		sudo apt-get install -y libsqlite3-dev || sudo yum install -y sqlite-devel; \
#	elif [ "$(COMPONENT)" = "datastore" ]; then \
#		echo "Installing ZMQ dependencies..."; \
#		sudo apt-get install -y libzmq3-dev || sudo yum install -y zeromq-devel; \
#	fi
endif

define go-build-template
build-$(1): make-output-dir
	$(if $(filter $(1),$(CGO_COMPONENTS)),\
		@$(MAKE) cgo-init COMPONENT=$(1) CGO_ENABLED=$(CGO_ENABLED))
	@echo "Building $(1) with CGO_ENABLED=$(CGO_ENABLED)"
	CGO_ENABLED=$(CGO_ENABLED) go build -v -o ./bin/$(1) \
	$(GO_TAGS) -ldflags "$(LD_FLAGS)" \
	./cmd/$(1)/main.go
endef

$(foreach comp,$(COMPONENTS),$(eval $(call go-build-template,$(comp))))

define docker-build-template
docker-build-$(1):
ifeq ($(FAST_BUILD),1)
	docker build \
		--build-arg CGO_ENABLED=$(CGO_ENABLED) \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME="$(BUILD_TIME)" \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg FEATURES=$(FEATURES) \
		--build-arg BUILD_HOST_PLATFORM=$(BUILD_HOST_PLATFORM) \
		--build-arg PLATFORM_VERSION="$(PLATFORM_VERSION)" \
		-t $(2) \
		-f cmd/$(1)/Dockerfile.FastBuild .
else
	docker build \
		--build-arg CGO_ENABLED=$(CGO_ENABLED) \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME="$(BUILD_TIME)" \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg FEATURES=$(FEATURES) \
		--build-arg BUILD_HOST_PLATFORM=$(BUILD_HOST_PLATFORM) \
		--build-arg PLATFORM_VERSION="$(PLATFORM_VERSION)" \
		-t $(2) \
		-f cmd/$(1)/Dockerfile .
endif
endef

$(eval $(call docker-build-template,agent, hub.iotroom.top/aenjoy/lubricant-agent:nightly))
$(eval $(call docker-build-template,gateway, hub.iotroom.top/aenjoy/lubricant-gateway:nightly))
$(eval $(call docker-build-template,apiserver, hub.iotroom.top/aenjoy/lubricant-apiserver:nightly))
$(eval $(call docker-build-template,logg, hub.iotroom.top/aenjoy/lubricant-logg:nightly))
$(eval $(call docker-build-template,grpcserver, hub.iotroom.top/aenjoy/lubricant-grpcserver:nightly))
$(eval $(call docker-build-template,reporter, hub.iotroom.top/aenjoy/lubricant-reporter:nightly))

# $(eval $(call docker-build-template,datastore, hub.iotroom.top/aenjoy/lubricant-datastore:nightly))

build-all: $(addprefix build-,$(COMPONENTS))

docker-build: $(addprefix docker-build-,$(COMPONENTS))

test: mock
	go test -v $(shell go list ./... | grep -v /integration)

test-coverage:
	go test -race -v -coverprofile=coverage.out -covermode=atomic $(shell go list ./... | grep -v /integration)

install:
	@command -v mockgen >/dev/null 2>&1 || go install github.com/golang/mock/mockgen@v1.6.0
	@command -v mockery >/dev/null 2>&1 || go install github.com/vektra/mockery/v2@latest

mock: install
	$(info Generating mocks...)
	@mockgen -source=services/gateway/repo/igateway.go -destination=pkg/mock/db/mockdb-gateway.go -package=db
	@mockgen -source=services/corepkg/repo/icore.go -destination=pkg/mock/db/mockdb-apiserver.go -package=db
	@mockery --dir=services/gateway/services/async --name=Task --output=pkg/mock/async --outpkg=async
	@mockery --dir=pkg/utils/mq --name=Mq --output=pkg/mock/mq --outpkg=mq
	@mockery --dir=protobuf/core --name=CoreServiceClient --output=pkg/mock/grpc --outpkg=grpc
	@mockery --dir=protobuf --name=BidiStreamingServer --output=pkg/mock/grpc --outpkg=grpc
	go mod tidy -v

copy-files:
	@$(foreach comp,$(COMPONENTS),\
		cp ./bin/$(comp) ./cmd/$(comp)/$(comp);)

load-to-kind: docker-build
	@$(foreach comp,$(COMPONENTS),\
		kind load docker-image hub.iotroom.top/aenjoy/lubricant-$(comp):nightly;)

clean:
	@rm -rf bin

help:
	@echo "Available targets:"
	@echo "  build-all          Build all components (CGO_ENABLED=0 by default)"
	@echo "  build-<component>  Build specific component"
	@echo "  docker-build       Build all Docker images"
	@echo "  test               Run unit tests"
	@echo "  test-coverage      Run tests with coverage"
	@echo "  mock               Generate mock files"
	@echo "  copy-files         Copy binaries to cmd directories"
	@echo "  clean              Clean build artifacts"
	@echo ""
	@echo "Environment variables:"
	@echo "  CGO_ENABLED=1      Enable CGO for supported components"
	@echo "  FAST_BUILD=1       Use fast Docker build method"
	@echo ""
	@echo "Example:"
	@echo "  make build-gateway CGO_ENABLED=1"
	@echo "  make docker-build FAST_BUILD=1"
