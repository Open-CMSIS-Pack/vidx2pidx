[![Release](https://github.com/Open-CMSIS-Pack/vidx2pidx/actions/workflows/release.yml/badge.svg)](https://github.com/Open-CMSIS-Pack/vidx2pidx/actions/workflows/release.yml)
[![Build](https://github.com/open-cmsis-pack/vidx2pidx/actions/workflows/build.yml/badge.svg)](https://github.com/open-cmsis-pack/vidx2pidx/actions/workflows/build.yml/badge.svg)
[![Tests](https://github.com/open-cmsis-pack/vidx2pidx/actions/workflows/test.yml/badge.svg)](https://github.com/open-cmsis-pack/vidx2pidx/actions/workflows/test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/open-cmsis-pack/vidx2pidx)](https://goreportcard.com/report/github.com/open-cmsis-pack/vidx2pidx)
[![GoDoc](https://godoc.org/github.com/open-cmsis-pack/vidx2pidx?status.svg)](https://godoc.org/github.com/open-cmsis-pack/vidx2pidx)

# vidx2pidx: Open-CMSIS-Pack Package Index Generator Tool

This is the git repository for the `vidx2pidx` tool. It takes in `*.vidx` file
and generate a `pidx`-formatted output listing packages.

## Install

Just head to the release page and download the binary for your system.


## Usage
```bash
$ vidx2pidx &lt;index&gt;.vidx

Options:

  -h, --help        show usage and help info
  -V, --version     show version and copyright info
  -v, --verbose     show progress details
  -o, --output      specify index file directory and name
  -c, --cachedir    specify directory where downloaded pidx and pdsc files are stored (default ./.idxcache)
  -f, --force       force update – ignore timestamp information
  ```

## Developing

Make sure to have Go [installed](https://golang.org/doc/install) in your environment.

```bash
$ git clone https://github.com/open-cmsis-pack/vidx2pidx
$ cd vidx2pidx
$ make
```
