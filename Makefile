################################################################################
# Variables                                                                    #
################################################################################

export GO111MODULE ?= on
export GOPROXY ?= https://goproxy.cn
export GOSUMDB ?= sum.golang.google.cn

GIT_COMMIT  = $(shell git rev-list -1 HEAD)
GIT_VERSION = $(shell git describe --always --abbrev=7 --dirty)
# By default, disable CGO_ENABLED. See the details on https://golang.org/cmd/cgo
CGO         ?= 0

LOCAL_ARCH := $(shell uname -m)
LOCAL_OS := $(shell uname)
ifeq ($(LOCAL_OS),Linux)
   GOLANGCI_LINT:=golangci-lint
else ifeq ($(LOCAL_OS),Darwin)
   GOLANGCI_LINT:=golangci-lint
else
   GOLANGCI_LINT:=golangci-lint.exe
endif

# Use the variable H to add a header (equivalent to =>) to informational output
H = $(shell printf "\033[34;1m=>\033[0m")

################################################################################
.PHONY: help
help: ## display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: all
all: help

.PHONY: tidy
tidy: ## up the go modules
	go mod tidy

.PHONY: lint
lint: ## fmt and lint the whole project
	go fmt .
	$(GOLANGCI_LINT) run --timeout=20m --sort-results

.PHONY: check-diff
check-diff: ## check no changes
	git diff --exit-code ./go.mod