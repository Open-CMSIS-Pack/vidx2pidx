package xml


import (
    "fmt"
    "os"
    "errors"
    "encoding/xml"
    "github.com/chaws/cmpack-idx-gen/utils"
)


//
//  cmpack-idx-gen.vidx
//
//  This file maintain a list of all available vendors
//  and their pidx file.
//
//  A small change to the original vidx file is the "Timestamp"
//  field to determine when the last time tha vendor pidx file
//  was modified
//
var vidxFileName = "cmpack-idx-gen.vidx"
type VidxXML struct {
    CommonIdx

    Vindex struct {
        XMLName xml.Name `xml:"vindex"`
        VendorPidxs []VendorPidx `xml:"pidx"`
    } `xml:"vindex"`
}


type VendorPidx struct {
    XMLName xml.Name `xml:"pidx"`
    Vendor string `xml:"vendor,attr"`
    URL string `xml:"url,attr"`
    Timestamp string `xml:"timestamp,attr"`
}


func (v *VidxXML) init() error {
    if _, err := os.Stat(vidxFileName); os.IsNotExist(err) {
      return v.save()
    }
   return utils.ReadXML(vidxFileName, v)
}


func (v *VidxXML) save() error {
    return utils.WriteXML(vidxFileName, v)
}


func (v *VidxXML) find(vendorName string) int {
    for i, pidx := range v.Vindex.VendorPidxs {
        if pidx.Vendor == vendorName {
            return i
        }
    }

    return -1
}


func (v *VidxXML) AddPidx(vendorName, pidxURL string) error {
    fmt.Printf("I: Adding '%s' (%s)\n", vendorName, pidxURL)

    if v.find(vendorName) != -1 {
        message := fmt.Sprintf("Vendor '%s' already exists.", vendorName)
        return errors.New(message)
    }

    vendorPidx := VendorPidx {
        Vendor: vendorName,
        URL: pidxURL,
        Timestamp: "",
    }
    v.Vindex.VendorPidxs = append(v.Vindex.VendorPidxs, vendorPidx)

    return v.save()
}


func (v *VidxXML) ListPidx() []VendorPidx {
    fmt.Println("I: Listing vendors pidx")
    return v.Vindex.VendorPidxs
}


func (v *VidxXML) RemovePidx(vendorName string) error {
    fmt.Printf("I: Removing vendor '%s'\n", vendorName)

    idx := v.find(vendorName)
    if idx == -1 {
        message := fmt.Sprintf("Vendor '%s' does not exist", vendorName)
        return errors.New(message)
    }

    pidxs := v.Vindex.VendorPidxs
    v.Vindex.VendorPidxs = append(pidxs[:idx], pidxs[idx+1:]...)
    return v.save()
}
