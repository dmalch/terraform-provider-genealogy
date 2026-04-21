GREEN:=\033[0;32m
YELLOW:=\033[0;33m
WHITE:=\033[0;37m
NC:=\033[0m # No Color

GOLANGCI_LINT_VERSION := v2.11.4
GOLANGCI_LINT := bin/golangci-lint

.PHONY: build-local
build-local:
	@echo "${WHITE}=====================${NC}"
	@echo "${YELLOW}Building...${NC}"
	go build -o bin/registry.terraform.io/dmalch/geni/0.6.4/darwin_arm64/terraform-provider-genealogy
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
