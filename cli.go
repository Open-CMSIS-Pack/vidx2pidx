package main

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var flags struct {
	outputFileName string
	force          bool
	cacheDir       string
	version        bool
}

func printVersionAndLicense(file io.Writer) {
	fmt.Fprintf(file, "vidx2pidx version %v\n", Version)
	fmt.Fprintf(file, "%v\n", License)
}

func NewCli() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vidx2pidx vendors.vidx",
		Short: "This utility converts a vendor index file into a vendor independent package index file.",
		Run: func(cmd *cobra.Command, args []string) {
			if flags.version {
				printVersionAndLicense(cmd.OutOrStdout())
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
			ExitOnError(Pidx.Update(Vidx))
			ExitOnError(WriteXML(flags.outputFileName, Pidx))
		},
	}

	cmd.Flags().StringVarP(&flags.outputFileName, "output", "o", "index.pidx", "Save pidx to this file")
	cmd.Flags().BoolVarP(&flags.version, "version", "V", false, "Output the version number of vidx2pidx and exit")
	cmd.Flags().BoolVarP(&flags.force, "force", "f", false, "Force the update, ignoring timestamps and digging package descriptor files")
	cmd.Flags().StringVarP(&flags.cacheDir, "cachedir", "c", ".idxcache", "Directory where to download and store pdsc files when using -f/--force flag")

	return cmd
}
