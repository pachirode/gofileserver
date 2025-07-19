TOOLS := golangci-lint goimports gotests mockgen

.PHONY: tools.verify
tools.verify: $(addprefix tools.verify., $(TOOLS))

.PHONY: tools.install
tools.install: $(addprefix tools.install., $(TOOLS))

.PHONY: tools.install.%
tools.install.%:
	@echo "==========> Installing $*"
	@$(MAKE) install.$*

.PHONY: tools.verify.%
tools.verify.%:
	@if ! which $* &>/dev/null; then $(MAKE) tools.install.$*; fi

.PHONY: install.golangci-lint
install.golangci-lint:
	@$(GO) install install github.com/golang/golangci-lint/cmd/golangci-lint@v1.50.1
	@golangci-lint completion bash > $(HOME)/.golangci-lint.bash
	@if ! grep -q .golangci-lint.bash $(HOME)/.bashrc; then echo "source \$$HOME/.golangci-lint.bash" >> $(HOME)/.bashrc; fi

.PHONY: install.goimports
install.goimports:
	@$(GO) install golang.org/x/tools/cmd/goimports@v0.4.0

.PHONY: install.gotests
install.gotests:
	@$(GO) install github.com/cweill/gotests/gotests@latest
