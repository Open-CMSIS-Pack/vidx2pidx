package commands


import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/chaws/cmpack/config"
    "os"
)


//
//  Command: vidx
//
//  Should be followed be "add", "rm" or "list"
//
var VidxCmd = &cobra.Command{
    Use:   "vidx",
    Short: "Add, list or remove vendor index sources",
    Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        unknownCommand(args[0], "vidx")
    },
}


//
//  Command: vidx add <vidx-name> <vidx-path>
//
//  Adds a vidx source to cmpack config file
//
var vidxAddCmd = &cobra.Command{
    Use: "add <name> <path>",
    Args: cobra.ExactArgs(2),
    Run: func(cmd *cobra.Command, args []string) {
        var name = args[0]
        var path = args[1]
        if err := config.AddVidx(name, path); err != nil {
            fmt.Fprintf(os.Stderr, "E: %s\n", err)
            os.Exit(-1)
        }
    },
}


func init() {
    VidxCmd.AddCommand(vidxAddCmd)
}
