package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bou.ke/monkey"
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
		pdsc := PdscTag{
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
		pdsc1 := PdscTag{
			Vendor:  "TheVendor",
			URL:     "http://the.url/",
			Name:    "TheName",
			Version: "0.0.1",
		}

		pdsc2 := PdscTag{
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

	t.Run("test force validating pdsc tag against actual pdsc file", func(t *testing.T) {
		xml := `<package>
			  <vendor>TheVendor</vendor>
			  <name>TheName</name>
			  <url>http://vendor.com/</url>
			  <timestamp></timestamp>
			  <releases>
			    <!-- The version is intentionally different, to show off a pdsc mismatch -->
			    <release version="0.0.2" />
			  </releases>
			</package>`
		pdscServer := HTTPServer(xml)

		pidx := NewPidx()

		// Force checking pdsc file
		pidx.SetForce(true)

		pdscTag := PdscTag{
			Vendor:  "TheVendor",
			URL:     pdscServer.URL + "/",
			Name:    "TheName",
			Version: "0.0.1",
		}

		err := pidx.addPdsc(pdscTag)
		if err == nil {
			t.Errorf("AddPdsc should return an error with the mismatchin version")
		}

		expected := fmt.Sprintf("Pdsc tag '%s%s' does not match the actual file:", pdscTag.URL, pdscTag.Name)
		expected += fmt.Sprintf(" URL('%s' != '%s')", "http://vendor.com/", pdscTag.URL)
		expected += fmt.Sprintf(" Version('%s' != '%s')", "0.0.2", pdscTag.Version)

		AssertEqual(t, err.Error(), expected)

		pdscList := pidx.ListPdsc()
		if len(pdscList) == 0 {
			t.Error("AddPdsc should still add the pdsc despite the mismatch")
		}

		if len(pdscList) != 1 {
			t.Error("AddPdsc should still add just one pdsc on mismatch")
		}

		newPdscTag := pdscList[0]
		if newPdscTag.Version != "0.0.2" {
			t.Error("AddPdsc on a mismatch should add the pdsc tag generated from the pdsc file")
		}
	})

	t.Run("test force validating is ignored when pdsc is unreachable", func(t *testing.T) {
		errMessage := "failed to read XML"
		monkey.Patch(ReadXML, func(string, interface{}) error {
			return errors.New(errMessage)
		})

		pidx := NewPidx()

		// Force checking pdsc file
		pidx.SetForce(true)

		pdscTag := PdscTag{
			Vendor:  "TheVendor",
			URL:     "http://vendor.com/",
			Name:    "TheName",
			Version: "0.0.1",
		}

		err := pidx.addPdsc(pdscTag)
		if err == nil {
			t.Errorf("AddPdsc should return an error when force checking a pdsc file that's unreachable")
		}

		AssertEqual(t, err.Error(), errMessage)

		pdscList := pidx.ListPdsc()
		if len(pdscList) == 0 {
			t.Error("AddPdsc should still add the pdsc despite an unreachable pdsc file")
		}

		if len(pdscList) != 1 {
			t.Error("AddPdsc should still add just one pdsc on unreachable pdsc file")
		}

		if pdscTag != pdscList[0] {
			t.Error("AddPdsc added something different than the original pdsc tag")
		}

		monkey.Unpatch(ReadXML)
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
		vidx.Pindex.Pdscs = append(vidx.Pindex.Pdscs, PdscTag{
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

func ExamplePidxXML_Update() {
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
	vidx.Pindex.Pdscs = append(vidx.Pindex.Pdscs, PdscTag{
		Vendor:  "TheOtherVendor",
		URL:     "http://other-vendor.com/",
		Name:    "ThePackage",
		Version: "0.0.1",
	})

	pidx := NewPidx()
	ExitOnError(pidx.Update(vidx))
	ExitOnError(WriteXML("", pidx))
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
