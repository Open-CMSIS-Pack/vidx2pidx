package main

import (
	"errors"
	"testing"

	"bou.ke/monkey"
)

func TestVidx(t *testing.T) {
	t.Run("test fail to read xml", func(t *testing.T) {
		errMessage := "failed to read XML"
		monkey.Patch(ReadXML, func(string, interface{}) error {
			return errors.New(errMessage)
		})

		vidx := NewVidx()
		err := vidx.Init("../test/dummy.xml")
		if err == nil {
			t.Error("VidxXML.Init() should fail if XML cannot be read")
		}

		AssertEqual(t, err.Error(), errMessage)

		monkey.Unpatch(ReadXML)
	})
}
