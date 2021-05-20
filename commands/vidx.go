package commands


import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/chaws/cmpack-idx-gen/xml"
    "os"
)


//
//  Command: vidx
//
//  Should be followed be "add", "rm" or "list"
//
var VidxCmd = &cobra.Command{
    Use:   "vidx",
    Short: "Add, list or remove Vendors",
    Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        unknownCommand(args[0], "vidx")
    },
}


//
//  Command: vidx add <vidx-name> <vidx-path>
//
//  Adds a Vendor pidx file to cmpack-idx-gen vidx
//
var VidxAddCmd = &cobra.Command{
    Use: "add <vendor-name> <vendor-pidx-url>",
    Args: cobra.ExactArgs(2),
    Run: func(cmd *cobra.Command, args []string) {
        var vendorName = args[0]
        var pidxURL = args[1]
        if err := xml.Vidx.AddPidx(vendorName, pidxURL); err != nil {
            fmt.Fprintf(os.Stderr, "E: %s\n", err)
            os.Exit(-1)
        }
    },
}


//
//  Command: vidx list
//
//  List all Vendor pidx
//
var VidxListCmd = &cobra.Command{
    Use: "list",
    Args: cobra.ExactArgs(0),
    Run: func(cmd *cobra.Command, args []string) {
        for i, pidx := range xml.Vidx.ListPidx() {
            fmt.Printf("I: [%d] %s %s\n", i, pidx.Vendor, pidx.URL)
        }
    },
}


//
//  Command: vidx rm <vendor-name>
//
//  Remove a vendor-pidx
//
var VidxRmCmd = &cobra.Command{
    Use: "rm",
    Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        vendorName := args[0]
        if err := xml.Vidx.RemovePidx(vendorName); err != nil {
            fmt.Fprintf(os.Stderr, "E: vendor '%s' does not seem to exist.", vendorName)
        }
    },
}


func init() {
    VidxCmd.AddCommand(
        VidxAddCmd,
        VidxListCmd,
        VidxRmCmd,
    )
}
