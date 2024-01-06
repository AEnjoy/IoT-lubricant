PROJECT_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

.PHONY: test test-coverage

test:
	go test -v $(shell go list ./... | grep -v /integration)

test-coverage:
	go test -race -v -coverprofile=coverage.out -covermode=atomic $(shell go list ./... | grep -v /integration)

