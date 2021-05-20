package xml


import (
    "encoding/xml"
)


type CommonIdx struct {
    XMLName xml.Name `xml:"index"`
    Timestamp string `xml:"timestamp"`
}


var Vidx = new(VidxXML)
var Pidx = new(PidxXML)


func Init(){
    Vidx.init()
    Pidx.init()
}
