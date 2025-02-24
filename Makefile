PROJECT_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

.PHONY: test test-coverage install mock

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
		-source=internal/model/repo/service.go -destination=pkg/mock/db/mockdb.go \
		-package=db
	mockery --dir=pkg/utils/mq --name=Mq --output=pkg/mock/mq --outpkg=mq
	mockery --dir=protobuf/core --name=CoreServiceClient --output=pkg/mock/grpc --outpkg=grpc
	mockery --dir=protobuf --name=BidiStreamingServer --output=pkg/mock/grpc --outpkg=grpc

build-agent:
	docker build -t hub.iotroom.top/aenjoy/lubricant-agent:nightly -f cmd/agent/Dockerfile .

build-gateway:
	echo "Gateway is not running at container"
	# docker build -t hub.iotroom.top/aenjoy/lubricant-gateway:nightly -f cmd/agent_proxy/Dockerfile .

build-core:
	docker build -t hub.iotroom.top/aenjoy/lubricant-core:nightly -f cmd/core/Dockerfile .

docker-build: build-agent build-gateway build-core

load-to-kind-agent: build-agent
	kind load docker-image hub.iotroom.top/aenjoy/lubricant-agent:nightly
load-to-kind-gateway: build-gateway
	kind load docker-image hub.iotroom.top/aenjoy/lubricant-gateway:nightly
load-to-kind-core: build-core
	kind load docker-image hub.iotroom.top/aenjoy/lubricant-core:nightly

load-to-kind: load-to-kind-agent load-to-kind-core
