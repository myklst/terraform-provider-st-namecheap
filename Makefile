NAME=namecheap

VERSION=2.2.0

COMMIT := $(shell git rev-parse --short HEAD)

ARCH := $(shell go env GOARCH)

OS := $(shell go env GOOS)

OS_ARCH=darwin_arm64

BINARY := "terraform-provider-namecheap_v$(VERSION)"

format:
	go fmt ./...

check:
	go vet ./...

test:
	go test -v ./namecheap/... -count=1 -cover

# Please set the following ENV variables for this test:
# NAMECHEAP_USER_NAME, NAMECHEAP_API_USER, NAMECHEAP_API_KEY, NAMECHEAP_TEST_DOMAIN, NAMECHEAP_USE_SANDBOX (optional, default is false)
testacc:
	TF_ACC=1 go test ./namecheap -v -run=TestAcc -count=1

build:
	go build -o ${BINARY}


local:
	go build -o $(BINARY) -ldflags='-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)' .
	rm -rf ~/.terraform/plugins/terraform-namecheap
	rm -rf ~/.terraform.d/plugins/registry.terraform.io/st/namecheap/$(VERSION)/$(OS)_$(ARCH)
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/st/namecheap/$(VERSION)/$(OS)_$(ARCH)/
	mv $(BINARY) ~/.terraform.d/plugins/registry.terraform.io/st/namecheap/$(VERSION)/$(OS)_$(ARCH)/
	chmod +x ~/.terraform.d/plugins/registry.terraform.io/st/namecheap/$(VERSION)/$(OS)_$(ARCH)/$(BINARY)


release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=darwin GOARCH=arm64 go build -o ./bin/${BINARY}_${VERSION}_darwin_arm64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

install_darwin_amd64: build
	mkdir -p ~/.terraform.d/plugins/localhost/namecheap/${NAME}/${VERSION}/darwin_amd64
	mv ${BINARY} ~/.terraform.d/plugins/localhost/namecheap/${NAME}/${VERSION}/darwin_amd64

install_darwin_arm64: build
	mkdir -p ~/.terraform.d/plugins/localhost/namecheap/${NAME}/${VERSION}/darwin_arm64
	mv ${BINARY} ~/.terraform.d/plugins/localhost/namecheap/${NAME}/${VERSION}/darwin_arm64

install_linux_amd64: build
	mkdir -p ~/.terraform.d/plugins/localhost/namecheap/${NAME}/${VERSION}/linux_amd64
	mv ${BINARY} ~/.terraform.d/plugins/localhost/namecheap/${NAME}/${VERSION}/linux_amd64

# Make sure you have installed golangci-lint CLI
# https://golangci-lint.run/usage/install/#local-installation
lint:
	golangci-lint run

# Make sure you have installed https://github.com/hashicorp/terraform-plugin-docs
docs:
	tfplugindocs

vendor:
	go mod vendor

.PHONY: format check test build release install_darwin_amd64 install_linux_amd64 lint docs vendor
