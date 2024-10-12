PROJECT_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

.PHONY: test test-coverage

test: mock
	go test -v $(shell go list ./... | grep -v /integration)

test-coverage:
	go test -race -v -coverprofile=coverage.out -covermode=atomic $(shell go list ./... | grep -v /integration)

install:
	@which mockgen > /dev/null || go install github.com/golang/mock/mockgen@v1.6.0

mock: install
	mockgen \
		-source=pkg/model/api.go -destination=pkg/mock/db/mockdb.go \
		-package=db
