package main

import (
	"testing"
	"errors"
)

func TestAnyErr(t *testing.T){
	var errs []error
	if AnyErr(errs) != nil {
		t.Errorf("AnyErr should return nothing when empty is given")
	}

	foo := errors.New("Foo error")
	errs = append(errs, foo)
	if AnyErr(errs) != foo {
		t.Errorf("AnyErr should return the first error in the array of errors")
	}

	bar := errors.New("Bar error")
	errs = append(errs, bar)
	if AnyErr(errs) != foo {
		t.Errorf("AnyErr should return the first error in the array of errors")
	}
}
