package commands


import (
    "fmt"
    "os"
    "strings"
    "github.com/spf13/cobra"
)


var BinName = "cmpack-idx-gen"


func unknownCommand(args ...string) {
    given := args[0]

    args[0] = BinName
    command := strings.Join(args, " ")

    fmt.Printf("E: Unknown command '%s'. See '%s --help'.\n", given, command)
    os.Exit(-1)
}


var rootCmd = &cobra.Command{
    Use:   BinName,
    Short: "Generates package index based on CMSIS-Pack vendors",
    Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        unknownCommand(args[0])
    },
}


func Run() {

    rootCmd.AddCommand(
        VidxCmd,
        UpdateCmd,
    )

    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
