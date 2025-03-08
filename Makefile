PROJECT_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

VERSION := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_TIME := $(shell date +"%Y-%m-%d %H:%M:%S")
GIT_COMMIT := $(shell git rev-parse --short HEAD)
GO_VERSION := $(shell go version | awk '{print $$3}')
FEATURES := $(or $(ENV_LUBRICANT_ENABLE_FEATURES),default)
BUILD_HOST_PLATFORM := $(shell uname -s | tr '[:upper:]' '[:lower:]')/$(shell uname -m)
ifeq ($(shell uname -s),Linux)
PLATFORM_VERSION := $(shell grep -E '^(NAME|VERSION)=' /etc/os-release | tr -d '"' | awk -F= '{print $$2}' | paste -sd ' ' -)
else ifeq ($(shell uname -s),Windows)
PLATFORM_VERSION := $(shell systeminfo | findstr /B /C:"OS Name" /C:"OS Version" | awk -F: '{print $$2}' | paste -sd ' ' -)
else
PLATFORM_VERSION := unknown
endif

.PHONY: test test-coverage install mock

make-output-dir:
	rm -rf ./bin
	mkdir -p ./bin

test: mock
	go test -v $(shell go list ./... | grep -v /integration)

test-coverage:
	go test -race -v -coverprofile=coverage.out -covermode=atomic $(shell go list ./... | grep -v /integration)

install:
	@if ! which mockgen > /dev/null 2>&1; then \
		echo "mockgen is not installed, installing..."; \
		go install github.com/golang/mock/mockgen@v1.6.0; \
	else \
		echo "mockgen is already installed"; \
	fi
	@if ! which mockery > /dev/null 2>&1; then \
		echo "mockery is not installed, installing..."; \
		go install github.com/vektra/mockery/v2@latest; \
	else \
		echo "mockery is already installed"; \
	fi

mock: install
	mockgen \
		-source=services/gateway/repo/igateway.go -destination=pkg/mock/db/mockdb-gateway.go \
		-package=db
	mockgen \
		-source=services/lubricant/repo/icore.go -destination=pkg/mock/db/mockdb-lubricant.go \
		-package=db
	mockery --dir=services/gateway/services/async --name=Task --output=pkg/mock/async --outpkg=async
	mockery --dir=pkg/utils/mq --name=Mq --output=pkg/mock/mq --outpkg=mq
	mockery --dir=protobuf/core --name=CoreServiceClient --output=pkg/mock/grpc --outpkg=grpc
	mockery --dir=protobuf --name=BidiStreamingServer --output=pkg/mock/grpc --outpkg=grpc
	# mockgen will add a dependency on the mockgen package, which is needed by ut but this is not present in go.mod.
	go mod tidy -v

build-agent:
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME="$(BUILD_TIME)" \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg FEATURES=$(FEATURES) \
		--build-arg BUILD_HOST_PLATFORM=$(BUILD_HOST_PLATFORM) \
		--build-arg PLATFORM_VERSION="$(PLATFORM_VERSION)" \
		-t hub.iotroom.top/aenjoy/lubricant-agent:nightly \
		-f cmd/agent/Dockerfile .

build-gateway: make-output-dir
	go build -o ./bin/lubricant-gateway \
	-tags=sonic -tags=avx -ldflags "\
	-w -s \
	-X 'main.Version=$(VERSION)' \
	-X 'main.BuildTime=$(BUILD_TIME)' \
	-X 'main.GoVersion=$(GO_VERSION)' \
	-X 'main.GitCommit=$(GIT_COMMIT)' \
	-X 'main.Features=$(FEATURES)' \
	-X 'main.BuildHostPlatform=$(BUILD_HOST_PLATFORM)' \
	-X 'main.PlatformVersion=$(PLATFORM_VERSION)' \
	" \
	./cmd/gateway/main.go ./cmd/gateway/start.go

build-gateway-container:
	docker build -t hub.iotroom.top/aenjoy/lubricant-gateway:nightly -f cmd/agent_proxy/Dockerfile .

build-core:
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME="$(BUILD_TIME)" \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg FEATURES=$(FEATURES) \
		--build-arg BUILD_HOST_PLATFORM=$(BUILD_HOST_PLATFORM) \
		--build-arg PLATFORM_VERSION="$(PLATFORM_VERSION)" \
		-t hub.iotroom.top/aenjoy/lubricant-core:nightly \
		-f cmd/lubricant/Dockerfile .

build-lubricant: build-core

docker-build: build-agent build-gateway build-lubricant

load-to-kind-agent: build-agent
	kind load docker-image hub.iotroom.top/aenjoy/lubricant-agent:nightly
load-to-kind-gateway: build-gateway
	kind load docker-image hub.iotroom.top/aenjoy/lubricant-gateway:nightly
load-to-kind-core: build-lubricant
	kind load docker-image hub.iotroom.top/aenjoy/lubricant-core:nightly

load-to-kind: load-to-kind-agent load-to-kind-core load-to-kind-gateway

test-driver-clock:
	docker build -t hub.iotroom.top/aenjoy/test-driver-clock:nightly \
		-f test/mock_driver/clock/Dockerfile test/mock_driver/clock
load-to-kind-test-driver-clock: test-driver-clock
	kind load docker-image hub.iotroom.top/aenjoy/test-driver-clock:nightly

load-to-kind-test-driver: load-to-kind-test-driver-clock

build-test-driver: test-driver-clock

load-test-driver: load-to-kind-test-driver

