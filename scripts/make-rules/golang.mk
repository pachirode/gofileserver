GO := go

GO_BUILD_FLAGS := -ldflags "$(GO_LDFLAGS)"

ifeq ($(GOOS),windows)
GO_OUT_EXIT := .exe
endif

ifeq ($(ROOT_PACKAGE),)
$(error the variable ROOT_PACKAGE must be set prior to including golang.mk)
endif

GOPATH := $(shell go env GOPATH)
ifeq ($(origin GOBIN), undefined)
GOBIN := $(GOPATH)/bin
endif

COMMANDS ?= $(filter-out %.md, $(wildcard $(ROOT_DIR)/cmd/gofileserver/*))
BINS ?= $(foreach cmd,${COMMANDS},$(notdir $(cmd)))

ifeq ($(COMMANDS),)
$(error Could not determine COMMANDS, set ROOT_DIR)
endif
ifeq ($(BINS),)
$(error Could not determine BINS, set ROOT_DIR)
endif

.PHONY: go.build.verify
go.build.verify:
	@if ! which go &>/dev/null; then echo "Cannot found go compile tool. Please install go first."; exit 1; fi

.PHONY: go.build.%
go.build.%:
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "==================> Build binary $(COMMAND) $(VERSION) for $(OS) $(ARCH)"
	@mkdir -p $(OUTPUT_DIR)/platforms/$(OS)/$(ARCH)
	@CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) $(GO) build $(GO_BUILD_FLAGS) -o $(OUTPUT_DIR)/platforms/$(OS)/$(ARCH)/$(COMMAND)$(GO_OUT_EXIT) $(ROOT_PACKAGE)/cmd/gofileserver/$(COMMAND)

.PHONY: go.build
go.build: go.build.verify $(addprefix go.build., $(PLATFORM))

.PHONY: print-BINS
print-BINS:
	@echo "BINS=$(BINS)"

.PHONY: go.format
go.format: tools.verify.goimports
	@find . -type f -name '*.go' | $(XARGS) gofmt -s -w
	@find . -type f -name '*.go' | $(XARGS) goimports -w -local $(ROOT_PACKAGE)
	@$(GO) mod edit -fmt

.PHONY: go.tidy
go.tidy:
	@$(GO) mod tidy

.PHONY: go.lint
go.lint: tools.verify.golangci-lint
	@echo "========> Run golangci to lint source codes"
	@golangci-lint run -c $(ROOT_DIR)/.golangci.yaml $(ROOT_DIR)/..
