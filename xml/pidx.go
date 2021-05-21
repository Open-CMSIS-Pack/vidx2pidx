package xml


import (
    "fmt"
    "os"
    "errors"
    "sync"
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


func updatePdscListTask(id int, vendorPidx VendorPidx, pidx *PidxXML, wg *sync.WaitGroup, err *error) {

    defer wg.Done()

    url := fmt.Sprintf("%s%s.pidx", vendorPidx.URL, vendorPidx.Vendor)
    fmt.Printf("I: [%d] Fetching packages list from %s\n", id, url)

    incomingPidx := new(PidxXML)
    if *err = utils.ReadXML(url, &incomingPidx); *err != nil {
        return
    }

    if vendorPidx.Timestamp == incomingPidx.Timestamp {
        // Nothing changed, avoid extra work
        return
    }

    // Save this timestamp to avoid extra work next time
    Vidx.setPidxTimestamp(id, incomingPidx.Timestamp)

    for _, pdsc := range incomingPidx.ListPdsc() {
        pidx.addPdsc(pdsc)
    }
}


func (p *PidxXML) Update() error {

    fmt.Println("I: Updating list of packages (pdsc)")

    var wg sync.WaitGroup
    errs := make([]error, Vidx.length())
    for i, vendorPidx := range Vidx.ListPidx() {
        wg.Add(1)
        go updatePdscListTask(i, vendorPidx, p, &wg, &errs[i])
    }

    wg.Wait()

    for _, err := range errs {
        if err != nil {
            return err
        }
    }

    err := Vidx.save()
    if err != nil {
        return err
    }

    return p.save()
}
