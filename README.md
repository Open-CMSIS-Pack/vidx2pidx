# vidx2pidx: Open-CMSIS-Pack Package Index Generator Tool

This is the git repository for the `vidx2pidx` tool. It takes in `*.vidx` file
and generate a `pidx`-formatted output listing packages.

## Install

Just head to the release page and download the binary for your system.


## Usage
```bash
$ vidx2pidx vendor.vidx

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
