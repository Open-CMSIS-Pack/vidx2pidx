# Having these will allow CI scripts to build for many OS's and ARCH's
OS   := $(or ${OS},${OS},linux)
ARCH := $(or ${ARCH},${ARCH},amd64)

# Path to lint tool
GOLINTER ?= golangci-lint
GOFORMATTER ?= gofmt

# Determine binary file name
BIN_NAME := vidx2pidx
PROG := build/$(BIN_NAME)
ifneq (,$(findstring windows,$(OS)))
    PROG=build/$(BIN_NAME).exe
endif

SOURCES := $(wildcard *.go)

all:
	@echo Pick one of:
	@echo $$ make $(PROG)
	@echo $$ make run
	@echo $$ make clean
	@echo
	@echo Build for different OS's and ARCH's by defining these variables. Ex:
	@echo $$ make OS=windows ARCH=amd64 build/$(BIN_NAME).exe  \# build for windows 64bits
	@echo $$ make OS=darwin  ARCH=amd64 build/$(BIN_NAME)       \# build for MacOS 64bits
	@echo
	@echo Clean everything
	@echo $$ make clean

$(PROG): $(SOURCES)
	@echo Building project
	GOOS=$(OS) GOARCH=$(ARCH) go build -o $(PROG)

run: $(PROG)
	@./$(PROG) $(ARGS) || true

lint:
	$(GOLINTER) run

format:
	$(GOFORMATTER) -s -w .

format-check:
	$(GOFORMATTER) -d . | tee format-check.out
	test ! -s format-check.out

.PHONY: test
test:
	TESTING=1 go test $(ARGS)

test-all: format-check coverage-check lint test

coverage-report: 
	TESTING=1 go test -coverprofile cover.out
	go tool cover -html=cover.out

coverage-check:
	TESTING=1 go test -coverprofile cover.out
	tail -n +2 cover.out | grep -v -e " 1$$" | grep -v main.go | tee coverage-check.out
	test ! -s coverage-check.out

clean:
	rm -rf build/*
