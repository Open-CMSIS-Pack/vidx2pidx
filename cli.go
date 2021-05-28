package main


import (
    "fmt"
    "github.com/spf13/cobra"
)


var flags struct {
    outputFileName string
    validatePidxFiles bool
}


// add -o option to print to a file
// add --validate-pdsc to make sure information in pidx are correct
var rootCmd = &cobra.Command{
    Use:   "vidx2pidx vendors.vidx",
    Short: "Generates package index based on CMSIS-Pack vendors",
    Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
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

    ExitOnError(rootCmd.Execute())
}
