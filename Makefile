VENODR_TOOL = dep
VENODR_CMD = $(if $(shell which $(VENODR_TOOL)),,$(error Please install $(VENODR_TOOL))) $(VENODR_TOOL)

default: build

build: ## Build gaos
	go build -v ./cmd/gaos

test: ## Test appleopensource package
	go test -v -race ./

vendor-list: ## List vendor packages
	$(VENODR_CMD) status

help: ## Print this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[33m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

clean:
	$(RM) ./gaos
