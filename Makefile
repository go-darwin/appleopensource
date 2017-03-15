GO_LDFLAGS = -X main.gitCommit=`git rev-parse --short HEAD`

VENODR_TOOL = dep
VENODR_CMD = $(if $(shell which $(VENODR_TOOL)),,$(error Please install $(VENODR_TOOL))) $(VENODR_TOOL)

default: build

build:  ## Build the appleopensource binary.
	go build -ldflags "$(GO_LDFLAGS)" ./cmd/appleopensource

install:  ## Install the appleopensource binary.
	go install -v -x -ldflags "$(GO_LDFLAGS)" ./cmd/appleopensource

run: build  ## Run the appleopensource with fetch the xnu 4 versions.
	./appleopensource fetch xnu $(shell ./appleopensource versions xnu) .
	@$(RM) xnu-*.tar.gz

test:  ## Test appleopensource package.
	go test -v -race $(shell go list ./... | grep -v vendor)

vendor/list:  ## List vendor packages.
	$(VENODR_CMD) status

vendor/clean:  ## clean vendor package files.
	@find vendor -type f \( -name '*_test.go' -o -name '.gitignore' -o -name '*appveyor.yml' -o -name '.travis.yml' -o -name 'circle.yml' -o -name '*.json' -o -name '*.flake8' -o -name 'generate-flag-types' -o -name 'runtests' \) -print -exec rm -f {} ";"
	@find vendor -type d \( -name 'testdata' -o -name 'examples' -o -name '.gx' -o -name 'autocomplete' \) -print | xargs rm -rf

help:  ## Print this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

clean:
	$(RM) ./appleopensource
