package commands


import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/chaws/cmpack-idx-gen/xml"
    "os"
)


//
//  Command: pidx
//
//  Should be followed be "update" or "list"
//
var PidxCmd = &cobra.Command{
    Use:   "pidx",
    Short: "Update or list all available packages (pdsc)",
    Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        unknownCommand(args[0], "pidx")
    },
}


//
//  Command: pidx update
//
//  Updates the list of all packages (pdsc)
//
var AddCmd = &cobra.Command{
    Use: "update",
    Args: cobra.ExactArgs(2),
    Run: func(cmd *cobra.Command, args []string) {
        if err := xml.Pidx.Update(); err != nil {
            fmt.Fprintf(os.Stderr, "E: %s\n", err)
            os.Exit(-1)
        }
    },
}


//
//  Command: pidx list
//
//  List all available pdsc's
//
var ListCmd = &cobra.Command{
    Use: "list",
    Args: cobra.ExactArgs(0),
    Run: func(cmd *cobra.Command, args []string) {
        for i, pdsc := range xml.Pidx.ListPdsc() {
            fmt.Printf("I: [%d] %s/%s %s\n", i, pdsc.Vendor, pdsc.Name, pdsc.Version)
        }
    },
}


func init() {
    PidxCmd.AddCommand(
        UpdateCmd,
        ListCmd
    )
}
