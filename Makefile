SHELL := /bin/sh

.PHONY: dep
dep:
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run
