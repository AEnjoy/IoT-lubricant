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
		-source=pkg/model/api.go -destination=pkg/mock/db/mockdb.go \
		-package=db
	mockery --dir=protobuf/core --name=CoreServiceClient --output=pkg/mock/grpc --outpkg=grpc
	mockery --dir=protobuf --name=BidiStreamingServer --output=pkg/mock/grpc --outpkg=grpc

