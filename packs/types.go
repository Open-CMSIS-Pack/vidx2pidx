package packs

import (
    "encoding/json"
    "encoding/xml"
    "io/ioutil"
    "fmt"
    "os"
    "path"
    "github.com/chaws/cmpack/utils"
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


//
//  Define common xml fields shared between vidx and pidx
//
type CommonIdx struct {
    XMLName xml.Name `xml:"index"`
    Vendor string `xml:"vendor"`
    URL string `xml:"url"`
    Timestamp string `xml:"timestamp"`
}


//
//  Struct that defines vidx XML file
//
type Vidx struct {
    CommonIdx

    Vindex struct {
        XMLName xml.Name `xml:"vindex"`
        Pidxs []struct {
            XMLName xml.Name `xml:"pidx"`
            Vendor string `xml:"vendor,attr"`
            URL string `xml:"url,attr"`
        } `xml:"pidx"`
    } `xml:"vindex"`
}


//
//  Struct that defines pidx XML file
//
type Pidx struct {
    CommonIdx

    Pindex struct {
        XMLName xml.Name `xml:"pindex"`
        Pdscs []struct {
            XMLName xml.Name `xml:"pdsc"`
            Vendor string `xml:"vendor,attr"`
            URL string `xml:"url,attr"`
            Name string `xml:"name,attr"`
            Version string `xml:"version,attr"`
        } `xml:"pdsc"`
    } `xml:"pindex"`
}


func ReadVidx(path string) (*Vidx, error) {
    vidx := new(Vidx)
    return vidx, utils.ReadXML(path, vidx)
}


func ReadPidx(path string) (*Pidx, error) {
    pidx := new(Pidx)
    return pidx, utils.ReadXML(path, pidx)
}
