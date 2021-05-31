package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var flags struct {
	outputFileName    string
	validatePidxFiles bool
	version           bool
}

func printVersionAndLicense() {
	fmt.Printf("vidx2pidx version %v\n", Version)
	fmt.Printf("%v\n", License)
}

// add -o option to print to a file
// add --validate-pdsc to make sure information in pidx are correct
var rootCmd = &cobra.Command{
	Use:   "vidx2pidx vendors.vidx",
	Short: "This utility converts a vendor index file into a vendor independent package index file.",
	Run: func(cmd *cobra.Command, args []string) {
		if flags.version {
			printVersionAndLicense()
			return
		}

		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "E: Empty arguments list. See --help for more information.\n")
			return
		}

		vidxFileName := args[0]

		fmt.Printf("I: Reading '%s'\n", vidxFileName)

		ExitOnError(Vidx.Init(vidxFileName))
		ExitOnError(Pidx.Update())
		ExitOnError(WriteXML(flags.outputFileName, Pidx))
	},
}

func RunCli() {
	rootCmd.PersistentFlags().StringVarP(&flags.outputFileName, "output", "o", "", "Save pidx to this file")
	rootCmd.PersistentFlags().BoolVarP(&flags.validatePidxFiles, "validate-pidx", "", false, "Validate pidx files by checking pdsc files")
	rootCmd.PersistentFlags().BoolVarP(&flags.version, "version", "V", false, "Output the version number of vidx2pidx and exit.")

	ExitOnError(rootCmd.Execute())
}
