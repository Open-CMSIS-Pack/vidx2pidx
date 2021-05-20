package main


import (
    "github.com/chaws/cmpack-idx-gen/xml"
    cli "github.com/chaws/cmpack-idx-gen/commands"
)


func Init() {
    xml.Init()
}


func main(){
    Init()
    cli.Run()
}
