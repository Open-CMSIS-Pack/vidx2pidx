package main

import (
	"bou.ke/monkey"
	"errors"
	"os"
	"fmt"
	"testing"
	"net/http"
	"net/http/httptest"
)

func TestAnyErr(t *testing.T) {
	var errs []error
	if AnyErr(errs) != nil {
		t.Errorf("AnyErr should return nothing when empty is given")
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
}

func TestReadURL(t *testing.T) {
	t.Run("test no server", func(t *testing.T){
		monkey.Patch(http.Get, func(string) (*http.Response, error){
			return nil, nil
		})

		response, err := ReadURL("http://server.not.found")
		if err == nil || len(response) > 0 {
			t.Error("ReadURL should return an empty response and an error on bad URLs")
		}
	})

	t.Run("test bad body", func(t *testing.T){
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

		if err.Error() != "unexpected EOF" {
			t.Errorf("ReadURL should return 'unexpected EOF', got '%v' instead", err)
		}
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
