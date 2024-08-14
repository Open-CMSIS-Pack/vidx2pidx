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
- update your `vendor.pidx` file as [documented](https://open-cmsis-pack.github.io/Open-CMSIS-Pack-Spec/main/html/packIndexFile.html#pidxFile)
  ```
  <?xml version="1.0" encoding="UTF-8" ?>
  <index schemaVersion="1.0.0" xs:noNamespaceSchemaLocation="PackIndex.xsd" xmlns:xs="http://www.w3.org/2001/XMLSchema-instance">
    <vendor>MyVendor</vendor>
    <url>https://www.MyVendor.com/pack/</url>
    <timestamp>2024-08-15T15:00:10</timestamp>
    <pindex>
      <pdsc url="https://www.MyVendor.com/pack/mypack/" vendor="MyVendor" name="MyPack" version="1.1.0"/>
      ...
    </pindex>
  </index>
  ```
- create a vendor index file as [documented]( )
  ```
  <?xml version="1.0" encoding="UTF-8" ?>
  <index schemaVersion="1.0" xmlns:xs="http://www.w3.org/2001/XMLSchema-instance" xs:noNamespaceSchemaLocation="PackIndex.xsd">
    <vendor>MyVendor</vendor>
    <url>www.MyVendor.com/pack</url>
    <timestamp>2024-08-T15:30:00</timestamp>
    <vindex>
      <pidx url="https://www.MyVendor.com/pack/" vendor="MyVendor" />
      ...
    </vindex>
  </index>
  ```
- invoke `vidx2pidx vendor.vidx` 

```bash
$ vidx2pidx <index>.vidx

Options:

  -h, --help        show usage and help info
  -V, --version     show version and copyright info
  -v, --verbose     show progress details
  -o, --output      specify index file directory and name
  -c, --cachedir    specify directory where downloaded pidx and pdsc files are stored (default ./.idxcache)
  -f, --force       force update – ignore timestamp information
  ```

Now the generated `index.pidx` can be used with cpackget to validate that all listed packs can be installed:
- invoke `cpackget init ./index.pidx -R ./pack_root_test` to use the generated index.pidx in pack_root_test/.Web/index.pidx
- invoke `cpackget --public -R /pack_root_test` to list all latest public pack versions

