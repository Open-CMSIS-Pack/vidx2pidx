[![Release](https://github.com/Open-CMSIS-Pack/vidx2pidx/actions/workflows/release.yml/badge.svg)](https://github.com/Open-CMSIS-Pack/vidx2pidx/actions/workflows/release.yml)
[![Build](https://github.com/open-cmsis-pack/vidx2pidx/actions/workflows/build.yml/badge.svg)](https://github.com/open-cmsis-pack/vidx2pidx/actions/workflows/build.yml/badge.svg)
[![Tests](https://github.com/open-cmsis-pack/vidx2pidx/actions/workflows/test.yml/badge.svg)](https://github.com/open-cmsis-pack/vidx2pidx/actions/workflows/test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/open-cmsis-pack/vidx2pidx)](https://goreportcard.com/report/github.com/open-cmsis-pack/vidx2pidx)
[![GoDoc](https://godoc.org/github.com/open-cmsis-pack/vidx2pidx?status.svg)](https://godoc.org/github.com/open-cmsis-pack/vidx2pidx)

[![Maintainability](https://qlty.sh/badges/bcdfb931-e2dc-4d84-b248-52e0f053f30c/maintainability.svg)](https://qlty.sh/gh/Open-CMSIS-Pack/projects/vidx2pidx)
[![Test Coverage](https://qlty.sh/badges/bcdfb931-e2dc-4d84-b248-52e0f053f30c/test_coverage.svg)](https://qlty.sh/gh/Open-CMSIS-Pack/projects/vidx2pidx)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/Open-CMSIS-Pack/vidx2pidx/badge)](https://securityscorecards.dev/viewer/?uri=github.com/Open-CMSIS-Pack/vidx2pidx)

# vidx2pidx: Open-CMSIS-Pack Package Index Generator Tool

This is the git repository for the `vidx2pidx` tool. It takes in `*.vidx` file
and generate a `pidx`-formatted output listing packages.

## Install

Just head to the release page and download the binary for your system.

## Usage

- update your `vendor.pidx` file as [documented](https://open-cmsis-pack.github.io/Open-CMSIS-Pack-Spec/main/html/packIndexFile.html#pidxFile)

  ```xml
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

- create a vendor index file as documented

  ```xml
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

- invoke `cpackget init ./index.pidx -R ./pack_root_test` to use the generated index.pidx in
 pack_root_test/.Web/index.pidx
- invoke `cpackget --public -R /pack_root_test` to list all latest public pack versions

## Clone and Build

```bash
$ git clone https://github.com/open-cmsis-pack/vidx2pidx
$ cd vidx2pidx
$ make
```

## License

vidx2pidx is licensed under Apache 2.0.

## Contributions and Pull Requests

Contributions are accepted under Apache 2.0. Only submit contributions where you have authored all of the code.

### Issues and Labels

Please feel free to raise an [issue on GitHub](https://github.com/Open-CMSIS-Pack/vidx2pidx/issues)
to report misbehavior (i.e. bugs) or start discussions about enhancements. This
is your best way to interact directly with the maintenance team and the community.
We encourage you to append implementation suggestions as this helps to decrease the
workload of the very limited maintenance team. 

We will be monitoring and responding to issues as best we can.
Please attempt to avoid filing duplicates of open or closed items when possible.
In the spirit of openness we will be tagging issues with the following:

- **bug** – We consider this issue to be a bug that will be investigated.

- **wontfix** - We appreciate this issue but decided not to change the current behavior.
	
- **enhancement** – Denotes something that will be implemented soon. 

- **future** - Denotes something not yet schedule for implementation.

- **out-of-scope** - We consider this issue loosely related to CMSIS. It might by implemented outside of CMSIS. Let us know about your work.
	
- **question** – We have further questions to this issue. Please review and provide feedback.

- **documentation** - This issue is a documentation flaw that will be improved in future.

- **review** - This issue is under review. Please be patient.
	
- **DONE** - We consider this issue as resolved - please review and close it. In case of no further activity this issues will be closed after a week.

- **duplicate** - This issue is already addressed elsewhere, see comment with provided references.

- **Important Information** - We provide essential informations regarding planned or resolved major enhancements.
