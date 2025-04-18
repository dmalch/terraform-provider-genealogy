GREEN:=\033[0;32m
YELLOW:=\033[0;33m
WHITE:=\033[0;37m
NC:=\033[0m # No Color

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
