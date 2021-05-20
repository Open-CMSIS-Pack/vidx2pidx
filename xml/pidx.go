package xml


import (
    "fmt"
    "os"
    "encoding/xml"
    "github.com/chaws/cmpack-idx-gen/utils"
)


//
//  cmpack-idx-gen.pidx
//
//  This file contains all available packages from
//  all vendors.
//
var pidxFileName = "cmpack-idx-gen.pidx"
type PidxXML struct {
    CommonIdx

    Pindex struct {
        XMLName xml.Name `xml:"pindex"`
        Pdscs []Pdsc `xml:"pdsc"`
    } `xml:"pindex"`
}

type Pdsc struct {
    XMLName xml.Name `xml:"pdsc"`
    Vendor string `xml:"vendor,attr"`
    URL string `xml:"url,attr"`
    Name string `xml:"name,attr"`
    Version string `xml:"version,attr"`
}


func (p *PidxXML) init() error {
    if _, err := os.Stat(pidxFileName); os.IsNotExist(err) {
      return p.save()
    }
   return utils.ReadXML(pidxFileName, p)
}


func (p *PidxXML) save() error {
    return utils.WriteXML(pidxFileName, p)
}


func (p *PidxXML) ListPdsc() []Pdsc{
    fmt.Println("I: Listing available packages (pdsc)")
    return p.Pindex.Pdscs
}


func (p *PidxXML) Update() error {
    fmt.Println("I: Updating list of packages (pdsc)")
    return nil
}
