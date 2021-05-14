// This module will be responsible for handling cli commands
package commands

import (
    "fmt"
    "github.com/alecthomas/kong"
)

/**
 *  Define cmdline parameters here
 */
var CMPack struct {
    Vidx struct {
        Add struct {
            Name string `arg help:"Vidx name. Ex: MyVidx"`
            Path string `arg help:"Vidx file. Ex: /path/to/my.vidx or http//example.com/my.vidx"`
        } `cmd help:"Add a vidx source"`
        // List struct {} `cmd help:"List vidx sources"`
        // Rm struct {} `cmd help:"Remove a vidx source"`
    } `cmd help:"Update packs."`

    // Update struct {
    // } `cmd help:"Update packs."`

    // Search struct {
    //    Filter []string `arg name:"filter" help:"What to search"`
    //} `cmd help:"Search packs."`
}

/**
 *  Parse cmdline arguments and invoke user's commands
 */
func Run() {
    fmt.Println("Generating cmdline args")

    ctx := kong.Parse(&CMPack)
    switch ctx.Command() {
        case "vidx add <name> <path>":
            fmt.Println("Adding new vidx file")
            args := CMPack.Vidx.Add
            Vidx.Add(args.Name, args.Path)

        // case "update":
        //    fmt.Println("Updating list of packs")

        // case "search <filter>":
        //    fmt.Println("Calling search", CMPack.Search.Filter)

        default:
            panic(ctx.Command())
    }
}
