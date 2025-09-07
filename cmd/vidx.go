/* SPDX-License-Identifier: Apache-2.0 */
/* Copyright Contributors to the vidx2pidx project. */

package main

import (
	"encoding/xml"
)

// VidxXML maps the VIDX file format.
// Ref: https://github.com/ARM-software/CMSIS_5/blob/develop/CMSIS/Utilities/PackIndex.xsd
type VidxXML struct {
	XMLName   xml.Name `xml:"index"`
	Timestamp string   `xml:"timestamp"`

	Vindex struct {
		XMLName     xml.Name     `xml:"vindex"`
		VendorPidxs []VendorPidx `xml:"pidx"`
	} `xml:"vindex"`

	Pindex struct {
		XMLName xml.Name  `xml:"pindex"`
		Pdscs   []PdscTag `xml:"pdsc"`
	} `xml:"pindex"`
}

// VendorPidx maps the <pidx> tag in VIDX files.
type VendorPidx struct {
	XMLName   xml.Name `xml:"pidx"`
	Vendor    string   `xml:"vendor,attr"`
	URL       string   `xml:"url,attr"`
	Timestamp string   `xml:"timestamp,attr"`
}

// NewVidx creates a new instance of the VidxXML struct.
func NewVidx() *VidxXML {
	return new(VidxXML)
}

// Init parses the contents of the XML in "vidxFileName" into a VidxXML struct.
func (v *VidxXML) Init(vidxFileName string) error {
	return ReadXML(vidxFileName, v)
}

// ListPidx returns a slice with all <pidx> tags in a VIDX file.
func (v *VidxXML) ListPidx() []VendorPidx {
	return v.Vindex.VendorPidxs
}

// ListPdsc retuns a slice with all <pdsc> tags in a VIDX file.
func (v *VidxXML) ListPdsc() []PdscTag {
	return v.Pindex.Pdscs
}

// PidxLength returns the count of <pidx> tags in a VIDX file.
func (v *VidxXML) PidxLength() int {
	return len(v.Vindex.VendorPidxs)
}

// PdscLength returns the count of <pdsc> tags in a VIDX file.
func (v *VidxXML) PdscLength() int {
	return len(v.Pindex.Pdscs)
}
