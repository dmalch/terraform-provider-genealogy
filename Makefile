GREEN:=\033[0;32m
YELLOW:=\033[0;33m
WHITE:=\033[0;37m
NC:=\033[0m # No Color

GOLANGCI_LINT_VERSION := v2.11.4
GOLANGCI_LINT := bin/golangci-lint

# VERSION and PLATFORM feed build-local's Terraform filesystem-mirror path, so it
# always matches the latest release tag and the host platform instead of a stale
# hardcoded value. VERSION falls back to 0.0.0-dev when no tag is reachable.
VERSION := $(or $(patsubst v%,%,$(shell git describe --tags --abbrev=0 2>/dev/null)),0.0.0-dev)
PLATFORM := $(shell go env GOOS)_$(shell go env GOARCH)

.PHONY: build-local
build-local:
	@echo "${WHITE}=====================${NC}"
	@echo "${YELLOW}Building...${NC}"
	go build -o bin/registry.terraform.io/dmalch/geni/$(VERSION)/$(PLATFORM)/terraform-provider-genealogy
	@echo "${YELLOW}Building...${NC} ${GREEN}Done${NC}"

.PHONY: build
build:
	@echo "${WHITE}=====================${NC}"
	@echo "${YELLOW}Building...${NC}"
	go build -o bin/terraform-provider-genealogy
	@echo "${YELLOW}Building...${NC} ${GREEN}Done${NC}"

.PHONY: clean
clean:
	@echo "${WHITE}=====================${NC}"
	@echo "${YELLOW}Cleaning...${NC}"
	rm -rf bin/
	@echo -e "${YELLOW}Cleaning...${NC} ${GREEN}Done${NC}"

.PHONY: test
test:
	@echo "${WHITE}=====================${NC}"
	@echo "${YELLOW}Testing...${NC}"
	go test -v ./...
	@echo "${YELLOW}Testing...${NC} ${GREEN}Done${NC}"

.PHONY: docs
docs:
	@echo "${WHITE}=====================${NC}"
	@echo "${YELLOW}Generating docs...${NC}"
	tfplugindocs generate --provider-name geni
	@echo "${YELLOW}Generating docs...${NC} ${GREEN}Done${NC}"

$(GOLANGCI_LINT):
	@echo "${WHITE}=====================${NC}"
	@echo "${YELLOW}Installing golangci-lint ${GOLANGCI_LINT_VERSION}...${NC}"
	@GOBIN=$(CURDIR)/bin go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	@echo "${YELLOW}Installing golangci-lint...${NC} ${GREEN}Done${NC}"

.PHONY: lint
lint: $(GOLANGCI_LINT)
	@echo "${WHITE}=====================${NC}"
	@echo "${YELLOW}Linting...${NC}"
	$(GOLANGCI_LINT) run
	@echo "${YELLOW}Linting...${NC} ${GREEN}Done${NC}"

.PHONY: lint-fix
lint-fix: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run --fix
