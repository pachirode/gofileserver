SHELL := /bin/bash
VERSION = 0.0.0

COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/ && pwd -P))
OUTPUT_DIR := $(ROOT_DIR)/_output

server:
	@go build -v -o $(OUTPUT_DIR)/gofileserver $(ROOT_DIR)/cmd/gofileserver/main.go

clean:
	@-rm -vrf $(OUTPUT_DIR)
