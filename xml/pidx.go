package xml


import (
    "fmt"
    "os"
    "errors"
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


func (p *Pdsc) equalTo(pdsc Pdsc) bool {
    return p.Vendor == pdsc.Vendor &&
           p.URL == pdsc.URL &&
           p.Name == pdsc.Name &&
           p.Version == pdsc.Version
}


func (p *PidxXML) init() error {
    if _, err := os.Stat(pidxFileName); os.IsNotExist(err) {
      return p.save()
    }
   return utils.ReadXML(pidxFileName, p)
}


func (p *PidxXML) addPdsc(pdsc Pdsc) error {
    idx := p.findPdsc(pdsc)
    if idx != -1 {
        message := fmt.Sprintf("Package %s/%s/%s already exists!", pdsc.Vendor, pdsc.Name, pdsc.Version)
        return errors.New(message)
    }
    p.Pindex.Pdscs = append(p.Pindex.Pdscs, pdsc)
    return nil
}


func (p *PidxXML) save() error {
    return utils.WriteXML(pidxFileName, p)
}


func (p *PidxXML) findPdsc(targetPdsc Pdsc) int {
    for i, pdsc := range p.Pindex.Pdscs {
        if pdsc.equalTo(targetPdsc) {
            return i
        }
    }

    return -1
}


func (p *PidxXML) ListPdsc() []Pdsc{
    fmt.Println("D: Listing available packages (pdsc)")
    return p.Pindex.Pdscs
}


func (p *PidxXML) Update() error {
    fmt.Println("I: Updating list of packages (pdsc)")
    for i, pidx := range Vidx.ListPidx() {
        url := fmt.Sprintf("%s%s.pidx", pidx.URL, pidx.Vendor)
        fmt.Printf("I: [%d] Fetching packages list from %s\n", i, url)

        incomingPidx := new(PidxXML)
        if err := utils.ReadXML(url, &incomingPidx); err != nil {
            return err
        }

        if pidx.Timestamp == incomingPidx.Timestamp {
            // Nothing changed, avoid extra work
            continue
        }

        // Save this timestamp to avoid extra work next time
        Vidx.Vindex.VendorPidxs[i].Timestamp = incomingPidx.Timestamp
        Vidx.save()

        for _, pdsc := range incomingPidx.ListPdsc() {
            p.addPdsc(pdsc)
        }
    }

    return p.save()
}
