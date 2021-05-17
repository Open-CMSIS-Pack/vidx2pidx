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


//
//  Command: vidx list
//
//  List available vidx sources
//
var vidxListCmd = &cobra.Command{
    Use: "list",
    Args: cobra.ExactArgs(0),
    Run: func(cmd *cobra.Command, args []string) {
        vidxSources := config.ListVidxs()
        for i := 0; i < len(vidxSources); i++ {
            vidx := vidxSources[i]
            fmt.Printf("%s: %s\n", vidx.Name, vidx.Path)
        }
    },
}


//
//  Command: vidx rm
//
//  Remove a vidx source
//
var vidxRmCmd = &cobra.Command{
    Use: "rm",
    Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        if err := config.RemoveVidx(args[0]); err != nil {
            fmt.Fprintf(os.Stderr, "E: vidx '%s' does not seem to exist.", args[0])
        }
    },
}


func init() {
    VidxCmd.AddCommand(vidxAddCmd)
    VidxCmd.AddCommand(vidxListCmd)
    VidxCmd.AddCommand(vidxRmCmd)
}
