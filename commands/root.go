package commands


import (
    "fmt"
    "github.com/spf13/cobra"
    "os"
    "strings"
)


var BinName = "cmpack"


func unknownCommand(args ...string) {
    given := args[0]

    args[0] = BinName
    command := strings.Join(args, " ")

    fmt.Printf("E: Unknown command '%s'. See '%s --help'.\n", given, command)
    os.Exit(-1)
}


var rootCmd = &cobra.Command{
    Use:   BinName,
    Short: "This a tool to manage local CMSIS-Pack packages",
    Long: `A useful tool that helps downloading CMSIS-Pack packages from various vendors.
Complete documentation is available at https://cmpack.com`,
    Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        unknownCommand(args[0])
    },
}


func Run() {
    rootCmd.AddCommand(VidxCmd)
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
