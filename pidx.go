package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"sync"
)

//
//  This file contains all available packages from
//  all vendors.
//
type PidxXML struct {
	XMLName   xml.Name `xml:"index"`
	Timestamp string   `xml:"timestamp"`

	Pindex struct {
		XMLName xml.Name  `xml:"pindex"`
		Pdscs   []PdscTag `xml:"pdsc"`
	} `xml:"pindex"`

	pdscList map[string]bool
}

type PdscTag struct {
	XMLName   xml.Name `xml:"pdsc"`
	Vendor    string   `xml:"vendor,attr"`
	URL       string   `xml:"url,attr"`
	Name      string   `xml:"name,attr"`
	Version   string   `xml:"version,attr"`
	Timestamp string   `xml:"timestamp,attr"`
}

func NewPidx() *PidxXML {
	p := new(PidxXML)
	p.pdscList = make(map[string]bool)
	return p
}

func (p *PidxXML) addPdsc(pdsc PdscTag) error {
	if p.pdscList[pdsc.getURL()] {
		message := fmt.Sprintf("Package %s/%s/%s already exists!", pdsc.Vendor, pdsc.Name, pdsc.Version)
		return errors.New(message)
	}
	p.Pindex.Pdscs = append(p.Pindex.Pdscs, pdsc)
	p.pdscList[pdsc.getURL()] = true
	return nil
}

func (p *PidxXML) ListPdsc() []PdscTag {
	Logger.Debug("Listing available packages")
	return p.Pindex.Pdscs
}

func updatePdscListTask(id int, vendorPidx VendorPidx, pidx *PidxXML, wg *sync.WaitGroup, err *error) {
	defer wg.Done()

	url := vendorPidx.URL + vendorPidx.Vendor + ".pidx"
	Logger.Info("[%d] Fetching packages list from %s", id, url)

	incomingPidx := new(PidxXML)
	if *err = ReadXML(url, &incomingPidx); *err != nil {
		return
	}

	Logger.Info("Adding pdscs")
	for _, pdsc := range incomingPidx.ListPdsc() {
		if *err = pidx.addPdsc(pdsc); *err != nil {
			return
		}
	}
}

func (p *PidxXML) Update(vidx *VidxXML) error {
	Logger.Info("Updating list of packages")

	var wg sync.WaitGroup
	var err error
	var errs []error

	// Process package index first
	errs = make([]error, vidx.PidxLength())
	for i, vendorPidx := range vidx.ListPidx() {
		wg.Add(1)
		go updatePdscListTask(i, vendorPidx, p, &wg, &errs[i])
	}

	wg.Wait()

	if err = AnyErr(errs); err != nil {
		return err
	}

	// Now process package descriptors (vendors without pidx files)
	errs = make([]error, vidx.PdscLength())
	for i, pdsc := range vidx.ListPdsc() {
		errs[i] = p.addPdsc(pdsc)
	}

	if err = AnyErr(errs); err != nil {
		return err
	}

	return nil
}

func (p *PdscTag) getURL() string {
	return p.URL + p.Vendor + "." + p.Name + ".pdsc"
}
