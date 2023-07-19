.PHONY: linux docs local
# ==================== [START] Global Variable Declaration =================== #
SHELL := /bin/bash
# 'shell' removes newlines
ARCH := $(shell go env GOARCH)

BASE_DIR := $(shell pwd)

COMMIT := $(shell git rev-parse --short HEAD)

OS := $(shell go env GOOS)

UNAME_S := $(shell uname -s)

VERSION=1.0.0
# VERSION := $(shell grep "version=" install.sh | cut -d= -f2)

BINARY := "terraform-provider-namecheap_v$(VERSION)"

# exports all variables
export
# ===================== [END] Global Variable Declaration ==================== #

local:
	go build -o $(BINARY) -ldflags='-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)' .
	rm -rf ~/.terraform/plugins/terraform-namecheap
	rm -rf ~/.terraform.d/plugins/registry.terraform.io/tao/namecheap/$(VERSION)/$(OS)_$(ARCH)
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/tao/namecheap/$(VERSION)/$(OS)_$(ARCH)/
	mv $(BINARY) ~/.terraform.d/plugins/registry.terraform.io/tao/namecheap/$(VERSION)/$(OS)_$(ARCH)/
	chmod +x ~/.terraform.d/plugins/registry.terraform.io/tao/namecheap/$(VERSION)/$(OS)_$(ARCH)/$(BINARY)

docs:
	@go generate


