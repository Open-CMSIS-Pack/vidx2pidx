package main

import (
	"bou.ke/monkey"
	"errors"
	"os"
	"testing"
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
