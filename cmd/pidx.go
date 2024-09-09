/* SPDX-License-Identifier: Apache-2.0 */
/* Copyright Contributors to the vidx2pidx project. */

package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// PidxXML maps the PIDX file format.
// Ref: https://github.com/ARM-software/CMSIS_5/blob/develop/CMSIS/Utilities/PackIndex.xsd
type PidxXML struct {
	XMLName   xml.Name `xml:"index"`
	Vendor    string   `xml:"vendor"`
	URL       string   `xml:"url"`
	Timestamp string   `xml:"timestamp"`

	Pindex struct {
		XMLName xml.Name  `xml:"pindex"`
		Pdscs   []PdscTag `xml:"pdsc"`
	} `xml:"pindex"`

	pdscList map[string]bool
	force    bool
}

// PdscTag maps a <pdsc> tag that goes in PIDX files.
type PdscTag struct {
	XMLName   xml.Name `xml:"pdsc"`
	Vendor    string   `xml:"vendor,attr"`
	URL       string   `xml:"url,attr"`
	Name      string   `xml:"name,attr"`
	Version   string   `xml:"version,attr"`
	Timestamp string   `xml:"timestamp,attr"`

	err error
}

// NewPidx creates a new instance of the PidXML struct.
// Additionally, it allocates memory to keep a map of existing PDSC tags.
func NewPidx() *PidxXML {
	p := new(PidxXML)
	p.pdscList = make(map[string]bool)
	return p
}

// addPdsc appends a PdscTag to the PidxXML struct.
// It makes sure that no duplicate PDSC tags are added.
// If PidxXML.force is true, it means that the PDSC tag info will be completely
// replaced by the info present in the actual PDSC file.
func (p *PidxXML) addPdsc(pdsc PdscTag) error {
	var err error
	pdscURL := pdsc.getURL()
	if p.pdscList[pdscURL] {
		message := fmt.Sprintf("Package %s/%s/%s already exists!", pdsc.Vendor, pdsc.Name, pdsc.Version)
		err = errors.New(message)
		pdsc.err = err
		return err
	}

	if p.force {
		// The pdsc info in the tag should be ignored
		// and the actual pdsc is retrieved to get info cross-checked

		incomingPdscXML := new(PdscXML)
		if err := ReadXML(pdscURL, &incomingPdscXML); err != nil {
			// If it can't get the pdsc file, consider the pdsc tag to be valid
			p.Pindex.Pdscs = append(p.Pindex.Pdscs, pdsc)
			p.pdscList[pdscURL] = true
			pdsc.err = err
			return err
		}

		// Validate tag against the actual pdsc file
		if err := incomingPdscXML.MatchTag(pdsc); err != nil {
			// Prioritize information from pdsc file rather than tag
			correctPdscTag := incomingPdscXML.Tag()
			p.Pindex.Pdscs = append(p.Pindex.Pdscs, correctPdscTag)

			// Mark both wrong and correct pdsc in pdscList
			// to avoid duplication
			p.pdscList[pdscURL] = true
			p.pdscList[correctPdscTag.getURL()] = true
			pdsc.err = err

			return err
		}
	}

	p.Pindex.Pdscs = append(p.Pindex.Pdscs, pdsc)
	p.pdscList[pdscURL] = true
	return nil
}

// ListPdsc returns a slice of PdscTag type.
func (p *PidxXML) ListPdsc() []PdscTag {
	Logger.Debug("Listing available packages")
	return p.Pindex.Pdscs
}

// updatePdscListTask is a thread that actually generates a flat list of all
// PDSC tags collected from the Pindex section of the VIDX file.
func updatePdscListTask(id int, vendorPidx VendorPidx, pidx *PidxXML, wg *sync.WaitGroup, errs [][]error) {
	defer wg.Done()

	errs[id] = make([]error, 1)
	url := vendorPidx.URL + vendorPidx.Vendor + ".pidx"
	Logger.Info("[%d] Fetching packages list from %s", id, url)

	incomingPidx := new(PidxXML)
	if err := ReadXML(url, &incomingPidx); err != nil {
		errs[id][0] = err
		return
	}

	Logger.Info("Adding pdscs")
	pdscs := incomingPidx.ListPdsc()
	errs[id] = make([]error, len(pdscs))
	for i, pdsc := range pdscs {
		if err := pidx.addPdsc(pdsc); err != nil {
			errs[id][i] = err
		}
	}
}

// Update generates a flattened PIDX file containing all PDSC tags collected from a given VIDX file.
func (p *PidxXML) Update(vidx *VidxXML, vidxFileName string, outputFileName string) error {
	Logger.Info("Updating list of packages")

	filename := filepath.Base(vidxFileName)
	p.Vendor = strings.TrimSuffix(filename, filepath.Ext(filename))
	p.URL, _ = filepath.Abs(outputFileName)
	p.URL = "file://" + filepath.ToSlash(p.URL)
	t := time.Now()
	p.Timestamp = t.Format("2006-01-02T15:04:05")

	var wg sync.WaitGroup

	// Process package index first
	errs := make([][]error, vidx.PidxLength()+vidx.PdscLength())
	for i, vendorPidx := range vidx.ListPidx() {
		wg.Add(1)
		go updatePdscListTask(i, vendorPidx, p, &wg, errs)
	}

	wg.Wait()

	// Now process package descriptors (vendors without pidx files)
	offset := vidx.PidxLength()
	for i, pdsc := range vidx.ListPdsc() {

		errs[i+offset] = make([]error, 1)
		errs[i+offset][0] = p.addPdsc(pdsc)
	}

	for _, e := range errs {
		if err := AnyErr(e); err != nil {
			return err
		}
	}

	return nil
}

// SetForce sets the force flag for the PidxXML struct.
func (p *PidxXML) SetForce(force bool) {
	p.force = force
}

// getURL returns a string representing the PdscTag URL.
func (p *PdscTag) getURL() string {
	return p.URL + p.Vendor + "." + p.Name + ".pdsc"
}
