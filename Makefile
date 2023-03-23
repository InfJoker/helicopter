GOCMD = go
GOTEST := $(GOCMD) test
GOVET := $(GOCMD) vet
WD = $(shell pwd)
BIN_DIR = $(WD)/bin
BIN_NAME = helicopter
EXPORT_RESULT ?= false

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

GENERATED_DIR = $(WD)/generated

PROTO_DIR = $(WD)/proto
PROTO_OUT_DIR = $(GENERATED_DIR)/proto
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)
PROTOS := $(patsubst $(PROTO_DIR)/%.proto,$(PROTO_OUT_DIR)/%.pb.go,$(PROTO_FILES))

SOURCE_DIRS = cmd internal
SOURCE_FILES := $(shell find $(SOURCE_DIRS) -type f -name '*.go') $(PROTOS)

.PHONY: all
all: help

## Build:
.PHONY: build
build: $(BIN_DIR)/$(BIN_NAME) ## Build your project and put the output binary in bin/

$(BIN_DIR)/$(BIN_NAME): $(SOURCE_FILES) compile_protos
	mkdir -p $(BIN_DIR)
	$(GOCMD) mod tidy
	$(GOCMD) build -o $(BIN_DIR)/$(BIN_NAME) ./cmd/helicopter

.PHONY: compile_protos
compile_protos: $(PROTOS)

$(PROTO_OUT_DIR)/%.pb.go: $(PROTO_DIR)/%.proto ## Compile the .proto files and output them to generated/proto
	mkdir -p $(PROTO_OUT_DIR)
	protoc -I="$(shell dirname $<)" \
    --go_out="$(PROTO_OUT_DIR)" --go_opt=paths=source_relative \
    --go-grpc_out="$(PROTO_OUT_DIR)" --go-grpc_opt=paths=source_relative \
    $<

.PHONY: clean
clean: ## Remove build related file
	rm -rf $(BIN_DIR)
	rm -rf $(GENERATED_DIR)

## Test:
.PHONY: test
test: ## Run the tests of the project
	$(GOTEST) -v -race ./...

coverage: ## Run the tests of the project and export the coverage
	$(GOTEST) -cover -covermode=count -coverprofile=profile.cov ./...
	$(GOCMD) tool cover -func profile.cov

## Help:
help: ## Show this help.
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
