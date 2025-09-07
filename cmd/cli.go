/* SPDX-License-Identifier: Apache-2.0 */
/* Copyright Contributors to the vidx2pidx project. */

package main

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// flags gathers cmdline flags in a single set.
// They are:
// - outputFileName: file name to be created with vidx2pidx's output
// - force: when enabled, will ignore the info from PDSC tag and read the PDSC file instead
// - cacheDir: will be the home for downloaded files used to generate the index file
// - version: when enabled, force vidx2pidx to print out version and license, then exit
var flags struct {
	outputFileName string
	force          bool
	cacheDir       string
	version        bool
}

// printVersion prints out vidx2pidx current version
func printVersion(file io.Writer) {
	fmt.Fprintf(file, "vidx2pidx version %v %v\n", version, CopyrightNotice)
}

// NewCli creates a new instance of vidx2pidx cli.
func NewCli() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vidx2pidx <index>.vidx",
		Short: "This utility converts a vendor index file into a vendor independent package index file.",
		Run: func(cmd *cobra.Command, args []string) {
			if flags.version {
				printVersion(cmd.OutOrStdout())
				return
			}

			if len(args) == 0 {
				Logger.Error("Empty arguments list. See --help for more information.")
				return
			}

			vidxFileName := args[0]
			Logger.Info("Reading '%s'\n", vidxFileName)

			Vidx := NewVidx()
			Pidx := NewPidx()
			Pidx.SetForce(flags.force)

			CacheDir = flags.cacheDir

			ExitOnError(EnsureDir(CacheDir))
			ExitOnError(Vidx.Init(vidxFileName))

			err := Pidx.Update(Vidx, vidxFileName, flags.outputFileName)
			if err != nil {
				Logger.Error(err.Error())
			}

			ExitOnError(WriteXML(flags.outputFileName, Pidx))
		},
	}

	cmd.Flags().StringVarP(&flags.outputFileName, "output", "o", "index.pidx", "Save pidx to this file")
	cmd.Flags().BoolVarP(&flags.version, "version", "V", false, "Output the version number of vidx2pidx and exit")
	cmd.Flags().BoolVarP(&flags.force, "force", "f", false, "Force the update, ignoring timestamps and digging package descriptor files")
	cmd.Flags().StringVarP(&flags.cacheDir, "cachedir", "c", ".idxcache", "Directory where to download and store pdsc files when using -f/--force flag")

	return cmd
}
