package commands


import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/chaws/cmpack-idx-gen/packs"
    "os"
)


//
//  Command: vidx
//
//  Should be followed be "add", "rm" or "list"
//
var UpdateCmd = &cobra.Command{
    Use:   "update",
    Short: "Update the list of packs",
    Args: cobra.ExactArgs(0),
    Run: func(cmd *cobra.Command, args []string) {
        if err := packs.UpdateList(); err != nil {
            fmt.Fprintf(os.Stderr, "E: %s\n", err)
            os.Exit(-1)
        }
    },
}
