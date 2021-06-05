package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"bou.ke/monkey"
)

func AssertEqual(t *testing.T, got, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Wanted \"%s\", got \"%s\" instead", want, got)
	}
}

func TestAnyErr(t *testing.T) {
	var errs []error
	if AnyErr(errs) != nil {
		t.Error("AnyErr should return nothing when empty is given")
	}

	foo := errors.New("Foo error")
	errs = append(errs, foo)
	if AnyErr(errs) != foo {
		t.Error("AnyErr should return the first error in the array of errors")
	}

	bar := errors.New("Bar error")
	errs = append(errs, bar)
	if AnyErr(errs) != foo {
		t.Error("AnyErr should return the first error in the array of errors")
	}
}

func TestExitOnError(t *testing.T) {
	var exitCode = 0
	monkey.Patch(os.Exit, func(code int) {
		exitCode = code
	})

	ExitOnError(nil)

	if exitCode != 0 {
		t.Error("ExitOnError should not exit when no error is given")
	}

	foo := errors.New("Foo error")
	ExitOnError(foo)

	if exitCode != -1 {
		t.Error("ExitOnError should exit when error is given")
	}

	monkey.Unpatch(os.Exit)
}

func TestReadURL(t *testing.T) {
	t.Run("test no server", func(t *testing.T) {
		monkey.Patch(http.Get, func(string) (*http.Response, error) {
			return nil, nil
		})

		response, err := ReadURL("http://server.not.found")
		if err == nil || len(response) > 0 {
			t.Error("ReadURL should return an empty response and an error on bad URLs")
		}

		monkey.Unpatch(http.Get)
	})

	t.Run("test bad response", func(t *testing.T) {
		notFoundServer := httptest.NewServer(
			http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				},
			),
		)

		_, err := ReadURL(notFoundServer.URL)
		if err == nil {
			t.Error("ReadURL should return an error when request does not return 2xx")
		}
		AssertEqual(t, err.Error(), "request did not return successfully (404)")
	})

	t.Run("test bad body", func(t *testing.T) {
		bodyErrorServer := httptest.NewServer(
			http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Length", "1")
				},
			),
		)
		response, err := ReadURL(bodyErrorServer.URL)
		if err == nil || len(response) > 0 {
			t.Error("ReadURL should return an empty response and an error on falty URLs")
		}

		AssertEqual(t, err.Error(), "unexpected EOF")
	})

	t.Run("test all good", func(t *testing.T) {
		goodResponse := []byte("all good")
		goodServer := httptest.NewServer(
			http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintln(w, goodResponse)
				},
			),
		)
		response, err := ReadURL(goodServer.URL)
		if err != nil || len(response) == 0 {
			t.Error("ReadURL should return OK")
		}
	})
}

func TestReadXML(t *testing.T) {
	var dummyXML struct {
		Dummy    xml.Name `xml:"dummy"`
		Contents string   `xml:"contents"`
	}

	t.Run("test opening remote file not found 404", func(t *testing.T) {
		notFoundServer := httptest.NewServer(
			http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				},
			),
		)

		err := ReadXML(notFoundServer.URL, &dummyXML)
		if err == nil {
			t.Error("ReadXML should return error when xml is not found in the remote server")
		}

		AssertEqual(t, err.Error(), "request did not return successfully (404)")
	})

	t.Run("test local xml file not found or fail to open", func(t *testing.T) {
		fileName := time.Now().String()
		err := ReadXML(fileName, &dummyXML)
		if err == nil {
			t.Error("ReadXML should return error when local XML file is not found")
		}

		if !os.IsNotExist(err) {
			t.Errorf("Error should be not found, but got \"%s\" instead", err)
		}
	})

	t.Run("test local xml file fails to read", func(t *testing.T) {
		errString := "failed to read file"

		// ioutil.ReadAll is just an alias to io.ReadAll
		monkey.Patch(io.ReadAll, func(r io.Reader) ([]byte, error) {
			var empty []byte
			return empty, errors.New(errString)
		})

		err := ReadXML("test/dummy.xml", &dummyXML)
		if err == nil {
			t.Error("ReadXML should return error when local XML file fails to read")
		}

		AssertEqual(t, err.Error(), errString)

		monkey.Unpatch(io.ReadAll)
	})

	t.Run("test read malformed xml", func(t *testing.T) {
		monkey.Patch(io.ReadAll, func(r io.Reader) ([]byte, error) {
			return []byte("<unclosed-tag"), nil
		})

		err := ReadXML("test/dummy.xml", &dummyXML)
		if err == nil {
			t.Error("ReadXML should return error when local XML file fails to read")
		}

		AssertEqual(t, err.Error(), "XML syntax error on line 1: unexpected EOF")

		monkey.Unpatch(io.ReadAll)
	})

	t.Run("test all good", func(t *testing.T) {
		monkey.Patch(io.ReadAll, func(r io.Reader) ([]byte, error) {
			return []byte("<dummy><contents>Dummy content</contents></dummy>"), nil
		})

		err := ReadXML("test/dummy.xml", &dummyXML)
		if err != nil {
			t.Error("ReadXML should not return error on valid XML files:")
		}

		AssertEqual(t, dummyXML.Contents, "Dummy content")

		monkey.Unpatch(io.ReadAll)
	})
}

func TestWriteXML(t *testing.T) {
	type dummyXML struct {
		Dummy    xml.Name `xml:"dummy"`
		Contents string   `xml:"contents"`
	}

	t.Run("test fail to parse xml to bytes", func(t *testing.T) {
		// Creates an unmarshable type
		unmarshable := map[string]interface{}{
			"foo": make(chan int),
		}

		err := WriteXML("", unmarshable)
		if err == nil {
			t.Error("WriteXML should return error on unmarshable content")
		}

		AssertEqual(t, err.Error(), "xml: unsupported type: map[string]interface {}")
	})

	t.Run("test fail to write to file", func(t *testing.T) {
		// Tests if WriteXML raises error when attempting to write to file
		// It's meant to be a simple test. Actual write errors will be displayed
		// to the user during runtime

		errMessage := "fail to write to file"
		monkey.Patch(os.WriteFile, func(name string, data []byte, perm os.FileMode) error {
			return errors.New(errMessage)
		})

		xml := new(dummyXML)
		err := WriteXML("test/dummy-out.xml", xml)
		if err == nil {
			t.Error("WriteXML should return error when it's not able to write to file")
		}

		AssertEqual(t, err.Error(), errMessage)

		monkey.Unpatch(os.WriteFile)
	})

	// Test to stdout is covered in ExampleWriteXMLToStdout

	t.Run("test write to file", func(t *testing.T) {
		fileName := "test/dummy-out.xml"

		xml := new(dummyXML)
		xml.Contents = "dummy content"
		err := WriteXML(fileName, xml)
		if err != nil {
			t.Errorf("WriteXML should not return error on valid xml and valid file: %s", err)
		}

		written, err2 := ioutil.ReadFile(fileName)
		if err2 != nil {
			t.Fatalf("Can't open file %s to test if XML got actually written: %s", fileName, err2)
		}

		AssertEqual(t, written, []byte(`<dummyXML>
 <dummy></dummy>
 <contents>dummy content</contents>
</dummyXML>`))
	})
}

func ExampleWriteXML() {
	type dummyXML struct {
		Dummy    xml.Name `xml:"dummy"`
		Contents string   `xml:"contents"`
	}
	xml := new(dummyXML)
	xml.Contents = "dummy content"
	ExitOnError(WriteXML("", xml))
	// Output:
	// <dummyXML>
	//  <dummy></dummy>
	//  <contents>dummy content</contents>
	// </dummyXML>
}
