# basic Go commands
GOCMD=go
GOBUILD=$(GOCMD) build

# get GOPATH according to OS
ifeq ($(OS),Windows_NT) # is Windows_NT on XP, 2000, 7, Vista, 10...
    GOPATH=$(go env GOPATH)
else
    GOPATH=$(shell go env GOPATH)
endif

export PATH := $(GOPATH)/bin:$(PATH)

# targets

.PHONY: install-linter
install-linter:
  define get_latest_lint_release
    curl -s "https://api.github.com/repos/golangci/golangci-lint/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
  endef
  LATEST_LINT_VERSION=$(shell $(call get_latest_lint_release))
  INSTALLED_LINT_VERSION=$(shell golangci-lint --version 2>/dev/null | awk '{print "v"$$4}')
  ifneq "$(INSTALLED_LINT_VERSION)" "$(LATEST_LINT_VERSION)"
    @echo "new golangci-lint version found:" $(LATEST_LINT_VERSION)
    @curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin latest
  endif

# run static analysis tools, configuration in ./.golangci.yml file
.PHONY: lint
lint: install-linter
	golangci-lint run ./...

.PHONY: build-server
build-server:
	@BINARY_NAME=http make build-linux-amd64

.PHONY: build-ingestor
build-ingestor:
	@BINARY_NAME=ingestor make build-linux-amd64

.PHONY: build-linux-amd64
build-linux-amd64:
	@GOOS=linux GOARCH=amd64 make build-binary

.PHONY: build-binary
build-binary:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) -o ./build/$(BINARY_NAME) ./cmd/$(BINARY_NAME)/main.go

.PHONY: gen.go
gen.go:
	@go install github.com/golang/mock/mockgen@latest
	go generate ./...

.PHONY: test
test:
	@echo "Running tests..."
	go test -short -failfast -coverprofile=coverage.out ./...

.PHONY: vulncheck
vulncheck:
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...