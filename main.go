package main


import (
    "github.com/chaws/cmpack-idx-gen/config"
    cli "github.com/chaws/cmpack-idx-gen/commands"
)


func Init() {
    config.Init()
}


func main(){
    Init()
    cli.Run()
}
