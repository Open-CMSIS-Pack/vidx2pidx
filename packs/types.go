package packs

import (
    "encoding/json"
    "io/ioutil"
    "fmt"
    "os"
    "path"
)

//
//  This file defines structs that models packs
//

type PacksList struct {
    Name string `json:"name"`
    Packs []Pack `json:"packs"`

    // Keep a sha256 sum of Pidx file from the original pidx file
    // that is provided by vendors
    PidxSum string `json:"pidx_sum"`
}

//
//  The Pack struct hold meaningful information
//  on a Pack
type Pack struct {
    Name string `json:"name"`
    IsInstalled bool `json:"is_installed"`
}


func (p *PacksList) write() {
    file, _ := json.MarshalIndent(*p, "", " ")
	_ = ioutil.WriteFile(p.filePath(), file, 0644)
}


func (p *PacksList) filePath() string {
    return path.Join(listsDir, p.Name)
}


func ReadPacksList(listName string) *PacksList {

    packsList := new(PacksList)
    packsList.Name = listName

    jsonFile, err := os.Open(packsList.filePath())
    if err != nil {
        fmt.Printf("W: Packs list for '%s' not found, creating one\n", listName)
        packsList.Name = listName
        packsList.write()
    }

    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    json.Unmarshal(byteValue, packsList)

    return packsList
}
