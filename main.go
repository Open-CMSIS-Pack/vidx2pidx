package main

import (
    "github.com/chaws/cmpack/config"
    cli "github.com/chaws/cmpack/commands"
)


func Init() {
    config.Init()
}


func main(){
    Init()
    cli.Run()
}
