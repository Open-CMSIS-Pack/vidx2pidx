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

.PHONY: test
test:
	TESTING=1 go test

clean:
	rm -rf build/*
