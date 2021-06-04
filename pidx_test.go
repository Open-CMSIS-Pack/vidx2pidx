package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func HTTPServer(output string) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, output)
			},
		),
	)
}

func TestAddPdsc(t *testing.T) {
	t.Run("test fail if adding two existing packages", func(t *testing.T) {
		pidx := NewPidx()
		pdsc := Pdsc{
			Vendor:  "TheVendor",
			URL:     "http://the.url/",
			Name:    "TheName",
			Version: "0.0.1",
		}

		err := pidx.addPdsc(pdsc)
		if err != nil {
			t.Errorf("AddPdsc should not return error on valid pdsc: %s", err)
		}

		pdscList := pidx.ListPdsc()
		if len(pdscList) == 0 || pdscList[0] != pdsc {
			t.Error("AddPdsc should not fail to add a valid pdsc")
		}

		err = pidx.addPdsc(pdsc)
		if err == nil {
			t.Errorf("AddPdsc should fail to add existing pdsc")
		}
	})

	t.Run("test adding two different packages", func(t *testing.T) {
		pidx := NewPidx()
		pdsc1 := Pdsc{
			Vendor:  "TheVendor",
			URL:     "http://the.url/",
			Name:    "TheName",
			Version: "0.0.1",
		}

		pdsc2 := Pdsc{
			Vendor:  "TheVendor2",
			URL:     "http://the.url/",
			Name:    "TheName2",
			Version: "0.0.1",
		}

		err := pidx.addPdsc(pdsc1)
		if err != nil {
			t.Errorf("AddPdsc should not return error on valid pdsc: %s", err)
		}

		pdscList := pidx.ListPdsc()
		if len(pdscList) != 1 || pdscList[0] != pdsc1 {
			t.Error("AddPdsc should not fail to add a valid pdsc")
		}

		err = pidx.addPdsc(pdsc2)
		if err != nil {
			t.Errorf("AddPdsc should not return error on valid pdsc: %s", err)
		}

		pdscList = pidx.ListPdsc()
		if len(pdscList) != 2 || pdscList[1] != pdsc2 {
			t.Error("AddPdsc should not fail to add a valid pdsc")
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("test fail to read xml", func(t *testing.T) {
		pidxServer := HTTPServer("<bad-xml")

		vidx := NewVidx()
		vidx.Vindex.VendorPidxs = append(vidx.Vindex.VendorPidxs, VendorPidx{
			Vendor: "TheVendor",
			URL:    pidxServer.URL + "/",
		})

		pidx := NewPidx()
		err := pidx.Update(vidx)
		if err == nil {
			t.Error("PidxXML.Update() should fail on malformed-xml")
		}

		AssertEqual(t, err.Error(), "XML syntax error on line 1: unexpected EOF")
	})

	t.Run("test fail to add existing pdsc", func(t *testing.T) {
		xml := `<index>
			  <vendor>TheVendor</vendor>
			  <url>http://vendor.com/</url>
			  <timestamp></timestamp>
			  <pindex>
			    <!-- the two values below are intentionally duplicated -->
			    <pdsc url="http://vendor.com/" vendor="TheVendor" name="ThePack" version="1.2.3" />
			    <pdsc url="http://vendor.com/" vendor="TheVendor" name="ThePack" version="1.2.3" />
			  </pindex>
			</index>`
		pidxServer := HTTPServer(xml)

		vidx := NewVidx()
		vidx.Vindex.VendorPidxs = append(vidx.Vindex.VendorPidxs, VendorPidx{
			Vendor: "TheVendor",
			URL:    pidxServer.URL + "/",
		})

		pidx := NewPidx()
		err := pidx.Update(vidx)
		if err == nil {
			t.Error("PidxXML.Update() should fail when adding an existing pdsc")
		}

		AssertEqual(t, err.Error(), "Package TheVendor/ThePack/1.2.3 already exists!")
	})

	t.Run("test fail to add existing pdsc from pindex", func(t *testing.T) {
		xml := `<index>
			  <vendor>TheVendor</vendor>
			  <url>http://vendor.com/</url>
			  <timestamp></timestamp>
			  <pindex>
			    <pdsc url="http://vendor.com/" vendor="TheVendor" name="ThePack" version="1.2.3" />
			  </pindex>
			</index>`
		pidxServer := HTTPServer(xml)

		vidx := NewVidx()
		vidx.Vindex.VendorPidxs = append(vidx.Vindex.VendorPidxs, VendorPidx{
			Vendor: "TheVendor",
			URL:    pidxServer.URL + "/",
		})
		vidx.Pindex.Pdscs = append(vidx.Pindex.Pdscs, Pdsc{
			Vendor:  "TheVendor",
			URL:     "http://vendor.com/",
			Name:    "ThePack",
			Version: "1.2.3",
		})

		pidx := NewPidx()
		err := pidx.Update(vidx)
		if err == nil {
			t.Error("PidxXML.Update() should fail when adding an existing pdsc")
		}

		AssertEqual(t, err.Error(), "Package TheVendor/ThePack/1.2.3 already exists!")
	})
}

func ExampleUpdateSinglePidx() {
	xml := `<index>
                  <vendor>TheVendor</vendor>
                  <url>http://vendor.com/</url>
                  <timestamp></timestamp>
                  <pindex>
                    <pdsc url="http://vendor.com/" vendor="TheVendor" name="ThePack" version="1.2.3" />
                    <pdsc url="http://vendor.com/" vendor="TheVendor" name="TheOtherPack" version="1.1.0" />
                  </pindex>
                </index>`

	pidxServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, xml)
			},
		),
	)

	vidx := NewVidx()
	vidx.Vindex.VendorPidxs = append(vidx.Vindex.VendorPidxs, VendorPidx{
		Vendor: "TheVendor",
		URL:    pidxServer.URL + "/",
	})

	pidx := NewPidx()
	err := pidx.Update(vidx)
	if err != nil {
		fmt.Println("Unexpected: ", err)
		return
	}

	WriteXML("", pidx)
	// Output:
	// <index>
	//  <timestamp></timestamp>
	//  <pindex>
	//   <pdsc vendor="TheVendor" url="http://vendor.com/" name="ThePack" version="1.2.3" timestamp=""></pdsc>
	//   <pdsc vendor="TheVendor" url="http://vendor.com/" name="TheOtherPack" version="1.1.0" timestamp=""></pdsc>
	//  </pindex>
	// </index>
}

func ExampleUpdateWithPidxAndPdsc() {
	Logger.SetLevel(DEBUG)
	xml := `<index>
                  <vendor>TheVendor</vendor>
                  <url>http://vendor.com/</url>
                  <timestamp></timestamp>
                  <pindex>
                    <pdsc url="http://vendor.com/" vendor="TheVendor" name="ThePack" version="1.2.3" />
                    <pdsc url="http://vendor.com/" vendor="TheVendor" name="TheOtherPack" version="1.1.0" />
                  </pindex>
                </index>`

	pidxServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, xml)
			},
		),
	)

	vidx := NewVidx()
	vidx.Vindex.VendorPidxs = append(vidx.Vindex.VendorPidxs, VendorPidx{
		Vendor: "TheVendor",
		URL:    pidxServer.URL + "/",
	})
	vidx.Pindex.Pdscs = append(vidx.Pindex.Pdscs, Pdsc{
		Vendor:  "TheOtherVendor",
		URL:     "http://other-vendor.com/",
		Name:    "ThePackage",
		Version: "0.0.1",
	})

	pidx := NewPidx()
	err := pidx.Update(vidx)
	if err != nil {
		fmt.Println("Unexpected: ", err)
		return
	}

	WriteXML("", pidx)
	// Output:
	// <index>
	//  <timestamp></timestamp>
	//  <pindex>
	//   <pdsc vendor="TheVendor" url="http://vendor.com/" name="ThePack" version="1.2.3" timestamp=""></pdsc>
	//   <pdsc vendor="TheVendor" url="http://vendor.com/" name="TheOtherPack" version="1.1.0" timestamp=""></pdsc>
	//   <pdsc vendor="TheOtherVendor" url="http://other-vendor.com/" name="ThePackage" version="0.0.1" timestamp=""></pdsc>
	//  </pindex>
	// </index>
}