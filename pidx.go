package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"sync"
)

var Pidx = new(PidxXML)

//
//  This file contains all available packages from
//  all vendors.
//
type PidxXML struct {
	XMLName   xml.Name `xml:"index"`
	Timestamp string   `xml:"timestamp"`

	Pindex struct {
		XMLName xml.Name `xml:"pindex"`
		Pdscs   []Pdsc   `xml:"pdsc"`
	} `xml:"pindex"`
}

type Pdsc struct {
	XMLName   xml.Name `xml:"pdsc"`
	Vendor    string   `xml:"vendor,attr"`
	URL       string   `xml:"url,attr"`
	Name      string   `xml:"name,attr"`
	Version   string   `xml:"version,attr"`
	Timestamp string   `xml:"timestamp,attr"`
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

func (p *PidxXML) ListPdsc() []Pdsc {
	Logger.Debug("Listing available packages")
	return p.Pindex.Pdscs
}

func updatePdscListTask(id int, vendorPidx VendorPidx, pidx *PidxXML, wg *sync.WaitGroup, err *error) {
	defer wg.Done()

	url := vendorPidx.URL + vendorPidx.Vendor + ".pidx"
	Logger.Info("[%d] Fetching packages list from %s\n", id, url)

	incomingPidx := new(PidxXML)
	if *err = ReadXML(url, &incomingPidx); *err != nil {
		return
	}

	if vendorPidx.Timestamp == incomingPidx.Timestamp {
		// Nothing changed, avoid extra work
		return
	}

	for _, pdsc := range incomingPidx.ListPdsc() {
		if *err = pidx.addPdsc(pdsc); *err != nil {
			return
		}
	}
}

func (p *PidxXML) Update() error {
	Logger.Info("Updating list of packages")

	var wg sync.WaitGroup
	var err error
	var errs []error

	// Process package index first
	errs = make([]error, Vidx.PidxLength())
	for i, vendorPidx := range Vidx.ListPidx() {
		wg.Add(1)
		go updatePdscListTask(i, vendorPidx, p, &wg, &errs[i])
	}

	wg.Wait()

	if err = AnyErr(errs); err != nil {
		return err
	}

	// Now process package descriptors (vendors without pidx files)
	errs = make([]error, Vidx.PdscLength())
	for i, pdsc := range Vidx.ListPdsc() {
		errs[i] = Pidx.addPdsc(pdsc)
	}

	if err = AnyErr(errs); err != nil {
		return err
	}

	return nil
}
