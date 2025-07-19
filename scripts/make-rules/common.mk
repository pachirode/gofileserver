SHELL := /bin/bash

COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/../../ && pwd -P))
OUTPUT_DIR := $(ROOT_DIR)/_output

ROOT_PACKAGE=github.com/pachirode/gofileserver
VERSION_PACKAGE=github.com/pachirode/gofileserver/pkg/version

ifeq ($(origin TMP_DIR),undefined)
TMP_DIR := $(OUTPUT_DIR)/tmp
$(shell mkdir -p $(TMP_DIR))
endif

ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --tags --always --match='v*')
endif

GIT_TREE_STATE := "dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
GIT_TREE_STATE="clean"
endif

GIT_COMMIT := $(shell git rev-parse HEAD)

GIT_LDFLAGS += \
	-X $(VERSION_PACKAGE).GitVersion=$(VERSION) \
	-X $(VERSION_PACKAGE).GitCommit=$(GIT_COMMIT) \
	-X $(VERSION_PACKAGE).GitTreeState=$(GIT_TREE_STATE) \
	-X $(VERSION_PACKAGE).BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

PLATFORM ?= darwin_amd64 windows_amd64 linux_amd64 linux_arm64

ifeq ($(origin PLATFORM), undefined)
	ifeq ($(origin GOOS), undefined)
		GOOS := $(shell go env GOOS)
	endif
	ifeq ($(origin GOARCH), undefined)
		GOARCH := $(shell go env GOARCH)
	endif
	PLATFORM := $(GOOS)_$(GOARCH)
	IMAGE_PLAT := linux_$(GOARCH)
else
	GOOS := $(word 1, $(subst _, ,$(PLATFORM)))
	GOARCH := $(word 2, $(subst _, ,$(PLATFORM)))
	IMAGE_PLAT := $(PLATFORM)
endif

ifndef v
	MAKEFLAGS += --no-print-directory
endif

XARGS := xargs --no-run-if-empty


