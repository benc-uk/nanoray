# Set ENV to dev, prod, etc. to load .env.$(ENV) file
ENV ?= 
-include .env
export
-include .env.$(ENV)
export

# Internal variables you don't want to change
REPO_ROOT := $(shell git rev-parse --show-toplevel)
SHELL := /bin/bash

GOLINT_PATH := $(REPO_ROOT)/.tools/golangci-lint              # Remove if not using Go
AIR_PATH := $(REPO_ROOT)/.tools/air                           # Remove if not using Go

.EXPORT_ALL_VARIABLES:
.PHONY: help image push build run lint lint-fix
.DEFAULT_GOAL := help

help: ## 💬 This help message :)
	@figlet $@ || true
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install-tools: ## 🔮 Install dev tools into project .tools directory
	@figlet $@ || true
	@$(GOLINT_PATH) > /dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./.tools/
	@$(AIR_PATH) -v > /dev/null 2>&1 || ( wget https://github.com/cosmtrek/air/releases/download/v1.51.0/air_1.51.0_linux_amd64 -q -O .tools/air && chmod +x .tools/air )

lint: ## 🔍 Lint & format check only, sets exit code on error for CI
	@figlet $@ || true
	$(GOLINT_PATH) run ./worker/

lint-fix: ## 📝 Lint & format, attempts to fix errors & modify code
	@figlet $@ || true
	$(GOLINT_PATH) run --fix ./worker/

image: check-vars ## 📦 Build container image from Dockerfile
	@figlet $@ || true
	docker build --file ./build/Dockerfile \
	--build-arg BUILD_INFO="$(BUILD_INFO)" \
	--build-arg VERSION="$(VERSION)" \
	--tag $(IMAGE_PREFIX):$(IMAGE_TAG) . 

push: check-vars ## 📤 Push container image to registry
	@figlet $@ || true
	docker push $(IMAGE_PREFIX):$(IMAGE_TAG)

build: ## 🔨 Run a local build without a container
	@figlet $@ || true
	#go build -o __CHANGE_ME__ $(SRC_DIR)/...
	#cd $(SRC_DIR); npm run build

run: ## 🏃 Run application, used for local development
	@figlet $@ || true
	#$(AIR_PATH) -c .air.toml

clean: ## 🧹 Clean up, remove dev data and files
	@figlet $@ || true
	@rm -rf pkg/proto/*.pb.go

gen-proto: ## 🚀 Generate protobuf files
	@figlet $@ || true
	@protoc > /dev/null 2>&1 || (echo "💥 Error! protoc is not installed!"; exit 1)
	@protoc-gen-go --help > /dev/null 2>&1 || (echo "💥 Error! protoc-gen-go is not installed!"; exit 1)
	@protoc --go_out=pkg/proto --go-grpc_out=pkg/proto \
	  --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --proto_path=pkg/proto \
	  pkg/proto/*.proto

check-vars:
	@if [[ -z "${IMAGE_REG}" ]]; then echo "💥 Error! Required variable IMAGE_REG is not set!"; exit 1; fi
	@if [[ -z "${IMAGE_NAME}" ]]; then echo "💥 Error! Required variable IMAGE_NAME is not set!"; exit 1; fi
	@if [[ -z "${IMAGE_TAG}" ]]; then echo "💥 Error! Required variable IMAGE_TAG is not set!"; exit 1; fi
	@if [[ -z "${VERSION}" ]]; then echo "💥 Error! Required variable VERSION is not set!"; exit 1; fi
