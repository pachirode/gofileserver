.DEFAULT_GOAL := all

.PHONY: all
all: go.format go.build

include scripts/make-rules/common.mk
include scripts/make-rules/tools.mk
include scripts/make-rules/golang.mk

.PHONY: build
build: go.tidy
	@$(MAKE) go.build

.PHONY: clean
clean:
	@echo "==========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)

.PHONY: format
format:
	@$(MAKE) go.format

.PHONY: tidy
tidy:
	@$(MAKE) go.tidy
