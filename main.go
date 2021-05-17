package main

import (
    cli "github.com/chaws/cmpack/commands"
)


func Init() {
    InitConfig()
}


func main(){
    Init()
    cli.Run()
}
