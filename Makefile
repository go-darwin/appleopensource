# ----------------------------------------------------------------------------
# global

SHELL=/usr/bin/env bash
.DEFAULT_GOAL=static

# hack for replace all whitespace to comma
comma=,
empty=
space=$(empty) $(empty)

# ----------------------------------------------------------------------------
# Go

ifneq ($(shell command -v go),)
GO_PATH  ?= $(shell go env GOPATH)
GO_OS    ?= $(shell go env GOOS)
GO_ARCH  ?= $(shell go env GOARCH)
GO_BIN   ?= $(CURDIR)/bin
GO_FLAGS ?=

PKG_NAME := $(subst $(GO_PATH)/src/,,$(CURDIR))
ifneq ($(shell go list ${GO_MOD_FLAGS} ./... > /dev/null 2>&1),)
GO_PKGS := $(shell go list ${GO_MOD_FLAGS} ./... | grep -v -e '.pb.go')
GO_APP_PKGS := $(shell go list ${GO_MOD_FLAGS} -f '{{if and (or .GoFiles .CgoFiles) (ne .Name "main")}}{{.ImportPath}}{{end}}' ${PKG_NAME}/...)
GO_TEST_PKGS := $(shell go list ${GO_MOD_FLAGS} -f='{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' ./...)
GO_VENDOR_PKGS=
ifneq ($(wildcard ./vendor),)  # exist vender directory
GO_VENDOR_PKGS = $(shell go list ${GO_MOD_FLAGS} -f '{{if and (or .GoFiles .CgoFiles) (ne .Name "main")}}./vendor/{{.ImportPath}}{{end}}' ./vendor/...)
endif
endif
endif


APP=aos
CMD=$(PKG_NAME)/cmd/$(APP)

CGO_ENABLED ?= 0

GO_GCFLAGS=all=-trimpath=${GO_PATH}/src
# https://tip.golang.org/doc/diagnostics.html#debugging
GO_GCFLAGS_DEBUG=all=-N -l -dwarflocationlists=true -smallframes  # https://tip.golang.org/doc/diagnostics.html#debugging
GO_GCFLAGS_CHECKPTR_FLAGS=-d=checkptr=1 -d=checkptr=2

GO_LDFLAGS=all=-s -w
GO_LDFLAGS_STATIC="-extldflags=-static -v"
GO_LDFLAGS_DEBUG=all=-compressdwarf=false

GO_ASMFLAGS=all=-trimpath=${GO_PATH}/src
GO_ASMFLAGS_DEBUG=

GO_BUILDTAGS=
ifeq (${CGO_ENABLED},0)
	GO_BUILDTAGS=osusergo netgo
endif
GO_BUILDTAGS_STATIC=static static_build
GO_INSTALLSUFFIX_STATIC=-installsuffix 'netgo'
GO_FLAGS += -tags='$(subst $(space),$(comma),${GO_BUILDTAGS})'

GO_FLAGS+=-trimpath
GO_FLAGS+=-gcflags='${GO_GCFLAGS}'
GO_FLAGS+=-ldflags='${GO_LDFLAGS}'
GO_FLAGS+=-asmflags='${GO_ASMFLAGS}'

GO_TEST ?= go test
GO_TEST_FUNC ?= .
GO_TEST_FLAGS ?=
GO_BENCH_FUNC ?= .
GO_BENCH_FLAGS ?= -benchmem
GO_TEST_COVERAGE_OUT ?= coverage.out

ifneq ($(wildcard go.mod),)  # exist go.mod
ifneq ($(wildcard ./vendor),)  # exist vender directory
	GO_FLAGS+=-mod=vendor
	GO_TEST_FLAGS+=-mod=vendor
	GO_BENCH_FLAGS+=-mod=vendor
endif  # ifneq ($(wildcard ./vendor),)
endif  # ifneq ($(wildcard go.mod),)

# ----------------------------------------------------------------------------
# defines

GOPHER=""
define target
@printf "$(GOPHER)  \\x1b[1;32m$(patsubst ,$@,$(1))\\x1b[0m\\n"
endef

# $1: package import path, $2 revision
define tools
$(call target,tools/$(@F))
@{ \
	printf "downloadnig $(@F) ...\\n\\n" ;\
	set -e ;\
	CGO_ENABLED=0 GOOS=${GO_OS} GOARCH=${GO_ARCH} GOBIN=${GO_BIN} \
		go install -v -tags='tools,osusergo,netgo,static,static_build' -mod=mod -modfile=tools/go.mod -gcflags="all=-trimpath=${GO_PATH}/src" -ldflags='all=-s -w "-extldflags=-static"' -asmflags="all=-trimpath=${GO_PATH}/src" -installsuffix 'netgo' ${1} ;\
}
endef
GOBIN=$PWD/bin CGO_ENABLED=0 go install -v -x -tags=tools -mod=mod -modfile=tools/go.mod -gcflags="all=-trimpath=$(go env GOPATH)/src" -ldflags='all=-s -w "-extldflags=-static"' -asmflags="all=-trimpath=$(go env GOPATH)/src" github.com/golangci/golangci-lint/cmd/golangci-lint

# ----------------------------------------------------------------------------
# targets

##@ build and install

.PHONY: $(APP)
$(APP):
	$(call target,${TARGET})
	@mkdir -p $(@D)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GO_OS) GOARCH=$(GO_ARCH) go build -v $(strip $(GO_FLAGS)) -o $@ $(CMD)

.PHONY: build
build: TARGET=build
build: $(APP)  ## Builds a dynamic executable or package.

.PHONY: debug
debug: GO_GCFLAGS=${GO_GCFLAGS_DEBUG}
debug: GO_LDFLAGS=${GO_LDFLAGS_DEBUG}
debug: GO_ASMFLAGS=${GO_ASMFLAGS_DEBUG}
debug: $(APP)  ## Builds a dynamic executable with disable inline and optimization for debugger.

.PHONY: static
static: GO_LDFLAGS+=${GO_LDFLAGS_STATIC}
static: GO_BUILDTAGS+=${GO_BUILDTAGS_STATIC}
static: GO_FLAGS+=${GO_INSTALLSUFFIX_STATIC}
static: $(APP)  ## Builds a static executable.

.PHONY: install
install: GO_LDFLAGS+=${GO_LDFLAGS_STATIC}
install: GO_BUILDTAGS+=${GO_BUILDTAGS_STATIC}
install: GO_FLAGS+=${GO_INSTALLSUFFIX_STATIC}
install:  ## Installs the executable or package.
	$(call target)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GO_OS) GOARCH=$(GO_ARCH) go install -v $(strip $(GO_FLAGS)) $(CMD)

.PHONY: pkg/install
pkg/install: GO_FLAGS+=${GO_MOD_FLAGS}
pkg/install: GO_LDFLAGS=
pkg/install: GO_BUILDTAGS=
pkg/install: GO_FLAGS=
pkg/install:
	$(call target)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GO_OS) GOARCH=$(GO_ARCH) go install -v ${GO_APP_PKGS}

##@ test, bench and coverage

.PHONY: test
test: CGO_ENABLED=1  # needs race test
test:  ## Runs package test including race condition.
	$(call target)
	CGO_ENABLED=$(CGO_ENABLED) $(GO_TEST) -v -race $(strip $(GO_TEST_FLAGS)) -run=$(GO_TEST_FUNC) $(GO_TEST_PKGS)

.PHONY: test
trace: CGO_ENABLED=1  # needs race test
trace: GO_FLAGS+=-trace=trace.out
trace:  ## Runs package test including race condition.
	$(call target)
	CGO_ENABLED=$(CGO_ENABLED) $(GO_TEST) -v -race $(strip $(GO_TEST_FLAGS)) -run=$(GO_TEST_FUNC) $(GO_TEST_PKGS)

.PHONY: bench
bench:  ## Take a package benchmark.
	$(call target)
	CGO_ENABLED=$(CGO_ENABLED) $(GO_TEST) -v $(strip $(GO_BENCH_FLAGS)) -run='^$$' -bench=$(GO_BENCH_FUNC) $(GO_TEST_PKGS)

.PHONY: coverage
coverage: CGO_ENABLED=1
coverage:  ## Takes packages test coverage.
	$(call target)
	CGO_ENABLED=$(CGO_ENABLED) $(GO_TEST) -v $(strip $(GO_TEST_FLAGS)) $(strip $(GO_FLAGS)) -covermode=atomic -coverpkg=./... -coverprofile=${GO_TEST_COVERAGE_OUT} $(GO_PKGS)

tools/go-junit-report:  ## Find or download go-junit-report.
tools/go-junit-report: ${GO_BIN}/go-junit-report
${GO_BIN}/go-junit-report:
ifeq (, $(shell test -f $@))
	$(call tools,github.com/jstemmer/go-junit-report)
GO_JUNIT_REPORT=${GO_BIN}/go-junit-report
endif

.PHONY: coverage/ci
coverage/ci: tools/go-junit-report
coverage/ci: CGO_ENABLED=1
coverage/ci: GO_LDFLAGS+=${GO_LDFLAGS_STATIC}
coverage/ci: GO_BUILDTAGS+=${GO_BUILDTAGS_STATIC}
coverage/ci: GO_FLAGS+=${GO_INSTALLSUFFIX_STATIC}
coverage/ci:  ## Takes packages test coverage, and output coverage results to CI artifacts.
	$(call target)
	@mkdir -p /tmp/artifacts /tmp/test-results
	CGO_ENABLED=$(CGO_ENABLED) $(GO_TEST) -v $(strip $(GO_TEST_FLAGS)) $(strip $(GO_FLAGS)) -covermode=atomic -coverpkg=./... -coverprofile=${GO_TEST_COVERAGE_OUT} $(GO_PKGS) 2>&1 | tee /dev/stderr | go-junit-report -set-exit-code > /tmp/test-results/junit.xml
	@if [[ -f '/tmp/artifacts/coverage.out' ]]; then go tool cover -html=/tmp/artifacts/coverage.out -o /tmp/artifacts/coverage.html; fi


##@ lint

.PHONY: lint
lint: lint/golangci-lint  ## Run all linters.

tools/golangci-lint:  # go get 'golangci-lint' binary
tools/golangci-lint: ${GO_BIN}/golangci-lint
${GO_BIN}/golangci-lint:
ifeq (, $(shell test -f $@))
	$(call tools,github.com/golangci/golangci-lint/cmd/golangci-lint,master)
GOLANGCI_LINT=${GO_BIN}/golangci-lint
endif

.PHONY: lint/golangci-lint
lint/golangci-lint: tools/golangci-lint .golangci.yml  ## Run golangci-lint.
	$(call target)
	@${GOLANGCI_LINT} run ./...


##@ mod

.PHONY: mod/get
mod/get:  ## Updates all module packages and go.mod.
	$(call target)
	@go get -u -m -v -x all

.PHONY: mod/tidy
mod/tidy:  ## Makes sure go.mod matches the source code in the module.
	$(call target)
	@go mod tidy -v

.PHONY: mod/vendor
mod/vendor: mod/tidy  ## Resets the module's vendor directory and fetch all modules packages.
	$(call target)
	@go mod vendor -v

.PHONY: mod/graph
mod/graph:  ## Prints the module requirement graph with replacements applied.
	$(call target)
	@go mod graph | modgraphviz | dot -Tpng -o mod-graph.png

.PHONY: mod/install
mod/install: mod/tidy mod/vendor
mod/install:  ## Install the module vendor package as an object file.
	$(call target)
	GO111MODULE=on go install -mod=vendor -v $(GO_VENDOR_PKGS)

.PHONY: mod/update
mod/update: mod/get mod/tidy mod/vendor mod/install  ## Updates all of vendor packages.
	@go mod edit -go 1.14

.PHONY: mod
mod: mod/tidy mod/vendor mod/install
mod:  ## Updates the vendoring directory using go mod.
	@go mod edit -go 1.14


##@ clean

.PHONY: clean
clean:  ## Cleanups binaries and extra files in the package.
	$(call target)
	@$(RM) -r ./bin *.out *.test *.prof trace.log


## boilerplate

.PHONY: boilerplate/go/%
boilerplate/go/%: BOILERPLATE_PKG_DIR=$(shell printf $@ | cut -d'/' -f3- | rev | cut -d'/' -f2- | rev)
boilerplate/go/%: BOILERPLATE_PKG_NAME=$(if $(findstring $@,cmd),main,$(shell printf $@ | rev | cut -d/ -f2 | rev))
boilerplate/go/%: hack/boilerplate/boilerplate.go.txt
boilerplate/go/%:  ## Creates a go file based on boilerplate.go.txt in % location.
	@if [ ! -d ${BOILERPLATE_PKG_DIR} ]; then mkdir -p ${BOILERPLATE_PKG_DIR}; fi
	@cat hack/boilerplate/boilerplate.go.txt <(printf "\\npackage ${BOILERPLATE_PKG_NAME}\\n") > $*
	@sed -i "s|YEAR|$(shell date '+%Y')|g" $*


##@ miscellaneous

.PHONY: AUTHORS
AUTHORS:  ## Creates AUTHORS file.
	@$(file >$@,# This file lists all individuals having contributed content to the repository.)
	@$(file >>$@,# For how it is generated, see `make AUTHORS`.)
	@printf "$(shell git log --format="\n%aN <%aE>" | LC_ALL=C.UTF-8 sort -uf)" >> $@

.PHONY: TODO
TODO:  ## Print the all of (TODO|BUG|XXX|FIXME|NOTE) in packages.
	@rg -t go -t asm -C 3 -e '(TODO|BUG|XXX|FIXME|NOTE)(\(.+\):|:)' --follow --hidden --glob='!vendor'

.PHONY: nolint
nolint:  ## Print the all of //nolint:... pragma in packages.
	@rg -t go -C 3 -e '//nolint.+' --follow --hidden --glob='!vendor'


##@ help

.PHONY: help
help:  ## Show make target help.
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[33m<target>\033[0m\n"} /^[a-zA-Z_0-9\/_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
