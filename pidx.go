package main


import (
    "fmt"
    "errors"
    "sync"
    "encoding/xml"
)


var Pidx = new(PidxXML)


var pidxFileName = ""


//
//  This file contains all available packages from
//  all vendors.
//
type PidxXML struct {
    XMLName xml.Name `xml:"index"`
    Timestamp string `xml:"timestamp"`


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


func (p *PidxXML) addPdsc(pdsc Pdsc) error {
    idx := p.findPdsc(pdsc)
    if idx != -1 {
        message := fmt.Sprintf("Package %s/%s/%s already exists!", pdsc.Vendor, pdsc.Name, pdsc.Version)
        return errors.New(message)
    }
    p.Pindex.Pdscs = append(p.Pindex.Pdscs, pdsc)
    return nil
}


func (p *PidxXML) findPdsc(targetPdsc Pdsc) int {
    // Use map for this
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
    if *err = ReadXML(url, &incomingPidx); *err != nil {
        return
    }

    if vendorPidx.Timestamp == incomingPidx.Timestamp {
        // Nothing changed, avoid extra work
        return
    }

    for _, pdsc := range incomingPidx.ListPdsc() {
        pidx.addPdsc(pdsc)
    }
}


func (p *PidxXML) Update() error {

    fmt.Println("I: Updating list of packages (pdsc)")

    var wg sync.WaitGroup
    errs := make([]error, Vidx.Length())
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

    return nil
}
