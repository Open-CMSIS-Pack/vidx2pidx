package main


import (
    "encoding/xml"
)


var Vidx = new(VidxXML)


//
//  This file maintain a list of all available vendors
//  and their pidx file.
//
//  A small change to the original vidx file is the "Timestamp"
//  field to determine when the last time tha vendor pidx file
//  was modified
//
type VidxXML struct {
    XMLName xml.Name `xml:"index"`
    Timestamp string `xml:"timestamp"`

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


func (v *VidxXML) Init(vidxFileName string) error {
   return ReadXML(vidxFileName, v)
}


func (v *VidxXML) ListPidx() []VendorPidx {
    return v.Vindex.VendorPidxs
}


//func (v *VidxXML) ListPdsc() []VendorPidx {
//    return v.Vindex.VendorPidxs
//}


func (v *VidxXML) Length() int {
    return len(v.Vindex.VendorPidxs)
}
