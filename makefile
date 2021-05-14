# Hack to parse extra args from `make run $@` to `$(PROG) $@`
# ref: https://stackoverflow.com/a/14061796/3908350
ifeq (run,$(firstword $(MAKECMDGOALS)))
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(RUN_ARGS):;@:)
endif


# Having these will allow CI scripts to build for many OS's and ARCH's
OS   := $(or ${OS},${OS},linux)
ARCH := $(or ${ARCH},${ARCH},amd64)


# Determine binary file name
PROG := build/cmpack
ifneq (,$(findstring windows,$(OS)))
    PROG=build/cmpack.exe
endif


SOURCES := main.go $(wildcard commands/*.go)


all:
	@echo Pick one of:
	@echo $$ make $(PROG)
	@echo $$ make run
	@echo $$ make clean
	@echo
	@echo Build for different OS's and ARCH's by defining these variables. Ex:
	@echo $$ make OS=windows ARCH=amd64 build/cmpack.exe  \# build for windows 64bits
	@echo $$ make OS=darwin  ARCH=amd64 build/cmpack       \# build for MacOS 64bits
	@echo
	@echo Clean everything
	@echo $$ make clean


$(PROG): $(SOURCES)
	@echo Building project
	GOOS=$(OS) GOARCH=$(ARCH) go build -o $(PROG)


run: $(PROG)
	@./$(PROG) $(RUN_ARGS) || true


clean:
	rm -rf build/*
