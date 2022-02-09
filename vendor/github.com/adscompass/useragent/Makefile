SHELL       =   /bin/sh
PKG_PREFIX  :=  github.com/balerter/balerter
TAG         ?=  latest

.SUFFIXES:
.PHONY: help test fuzz-build fuzz

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test: ## Run tests
	go test ./...

fuzz-build: ## Build go fuzz
	go-fuzz-build

fuzz: ## Fuzz
	go-fuzz
