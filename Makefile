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

DOCKER_TAG := $(if $(RELEASE),$(RELEASE),nightly)

CGO_ENABLED ?= 0
CGO_COMPONENTS := gateway datastore

GO_TAGS := -tags=sonic,avx
LD_FLAGS = -w -s \
    -X 'github.com/aenjoy/iot-lubricant/pkg/version.Version=$(VERSION)' \
    -X 'github.com/aenjoy/iot-lubricant/pkg/version.BuildTime=$(BUILD_TIME)' \
    -X 'github.com/aenjoy/iot-lubricant/pkg/version.GoVersion=$(GO_VERSION)' \
    -X 'github.com/aenjoy/iot-lubricant/pkg/version.GitCommit=$(GIT_COMMIT)' \
    -X 'github.com/aenjoy/iot-lubricant/pkg/version.Features=$(FEATURES)' \
    -X 'github.com/aenjoy/iot-lubricant/pkg/version.BuildHostPlatform=$(BUILD_HOST_PLATFORM)' \
    -X 'github.com/aenjoy/iot-lubricant/pkg/version.PlatformVersion=$(PLATFORM_VERSION)'

COMPONENTS := gateway apiserver agent logg grpcserver reporter datastore

.PHONY: all test test-coverage install mock docker-build clean help make-output-dir load-test-driver list-components install-tdengine-driver

help:
	@echo "Available targets:"
	@echo "  list-components    List all available build components"
	@echo "  build-all          Build all components (CGO_ENABLED=0 by default)"
	@echo "  build-<component>  Build specific component"
	@echo "  docker-build       Build all Docker images"
	@echo "  test               Run unit tests"
	@echo "  test-coverage      Run tests with coverage"
	@echo "  mock               Generate mock files"
	@echo "  copy-files         Copy binaries to cmd directories"
	@echo "  load-to-kind       Load images to kind cluster"
	@echo "  install-driver     Install TDengine driver"
	@echo "  push-image         Push Docker images to registry"
	@echo "  clean              Clean build artifacts"
	@echo "  help               Show this help"
	@echo ""
	@echo "Environment variables:"
	@echo "  CGO_ENABLED=1      Enable CGO for supported components"
	@echo "  FAST_BUILD=1       Use fast Docker build method"
	@echo "  RELEASE=<tag>      Specify a release tag for Docker images(default: nightly)"
	@echo ""
	@echo "Example:"
	@echo "   CGO_ENABLED=1 make build-gateway"
	@echo "   FAST_BUILD=1 make docker-build"
	@echo "   RELEASE=v0.0.1 FAST_BUILD=1 make docker-build"
	@echo "   make load-to-kind"

list-components:
	@echo "Available components:"
	@for comp in $(COMPONENTS); do \
		echo "  $$comp"; \
	done
	@echo "Must CGO_ENABLED components:"
	@for comp in $(CGO_COMPONENTS); do \
    		echo "  $$comp"; \
    done

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
		@$(MAKE) cgo-init COMPONENT=$(1))
	@echo "Building $(1) with CGO_ENABLED=$(if $(filter $(1),$(CGO_COMPONENTS)),1,$(CGO_ENABLED))"
	CGO_ENABLED=$(if $(filter $(1),$(CGO_COMPONENTS)),1,$(CGO_ENABLED)) go build -v -o ./bin/$(1) \
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
		-t $(2):$(DOCKER_TAG) \
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
		-t $(2):$(DOCKER_TAG) \
		-f cmd/$(1)/Dockerfile .
endif
endef

$(eval $(call docker-build-template,agent, hub.iotroom.top/aenjoy/lubricant-agent))
$(eval $(call docker-build-template,gateway, hub.iotroom.top/aenjoy/lubricant-gateway))
$(eval $(call docker-build-template,apiserver, hub.iotroom.top/aenjoy/lubricant-apiserver))
$(eval $(call docker-build-template,logg, hub.iotroom.top/aenjoy/lubricant-logg))
$(eval $(call docker-build-template,grpcserver, hub.iotroom.top/aenjoy/lubricant-grpcserver))
$(eval $(call docker-build-template,reporter, hub.iotroom.top/aenjoy/lubricant-reporter))
$(eval $(call docker-build-template,datastore, hub.iotroom.top/aenjoy/lubricant-datastore))

load-test-driver:
	docker build -t hub.iotroom.top/aenjoy/test-driver-clock:nightly \
		-f scripts/test/mock_driver/clock/Dockerfile scripts/test/mock_driver/clock
	kind load docker-image hub.iotroom.top/aenjoy/test-driver-clock:nightly

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
	# mockgen will add a dependency on the mockgen package, which is needed by ut but this is not present in go.mod.
	go mod tidy -v

copy-files:
	@$(foreach comp,$(COMPONENTS),\
		cp ./bin/$(comp) ./cmd/$(comp)/$(comp);)

load-to-kind: docker-build
	@$(foreach comp,$(COMPONENTS),\
		kind load docker-image hub.iotroom.top/aenjoy/lubricant-$(comp):nightly;)

install-driver:
	wget https://static.iotroom.top/TDengine-client-3.3.5.2-Linux-x64.tar.gz -O /tmp/tdengine.tar.gz && \
    	tar -xzf /tmp/tdengine.tar.gz -C /tmp/ && \
    	cd /tmp/TDengine-client-3.3.5.2 && \
    	bash install_client.sh
push-image:
	@$(foreach comp,$(COMPONENTS),\
		docker push hub.iotroom.top/aenjoy/lubricant-$(comp):$(DOCKER_TAG);)

clean:
	@rm -rf bin
