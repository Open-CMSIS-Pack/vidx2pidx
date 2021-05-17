// This module will be responsible for handling cli commands
package commands

import (
    "fmt"
    "github.com/spf13/cobra"
    "os"
)

var rootCmd = &cobra.Command{
    Use:   "cmpack",
    Short: "cmpack is a tool to manage local CMSIS-Pack packages",
    Long: `A useful tool that helps downloading CMSIS-Pack packages from various vendors.
Complete documentation is available at https://cmpack.com`,
    Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Running cmpack")
        fmt.Println(args)
    },
}

/**
 *  Parse cmdline arguments and invoke user's commands
 */
func Run() {
    rootCmd.AddCommand(VidxCmd)
    fmt.Println("Generating cmdline args")
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
