/* SPDX-License-Identifier: Apache-2.0 */
/* Copyright Contributors to the vidx2pidx project. */

package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCli(t *testing.T) {
	t.Run("test empty parameters list", func(t *testing.T) {
		currLogFile := Logger.file
		output := bytes.NewBufferString("")
		Logger.SetFile(output)

		cmd := NewCli()
		err := cmd.Execute()
		if err != nil {
			t.Errorf("Cli should not fail on empty arg list: %s", err)
		}

		var out []byte
		out, err = io.ReadAll(output)
		if err != nil {
			t.Fatal(err)
		}

		AssertEqual(t, strings.TrimSpace(string(out)), "E: Empty arguments list. See --help for more information.")

		Logger.SetFile(currLogFile)
	})

	t.Run("test display version and license", func(t *testing.T) {
		output := bytes.NewBufferString("")

		cmd := NewCli()
		cmd.SetOut(output)
		cmd.SetArgs([]string{"--version"})
		err := cmd.Execute()
		if err != nil {
			t.Errorf("Cli should not fail when requesting version: %s", err)
		}

		var out []byte
		out, err = io.ReadAll(output)
		if err != nil {
			t.Fatal(err)
		}

		AssertEqual(t, strings.Contains(string(out), "vidx2pidx version"), true)
	})

	t.Run("test continue despite errors", func(t *testing.T) {
		outputFileName := "test-continue-despite-errors.xml"

		currLogFile := Logger.file
		currLevel := Logger.level
		output := bytes.NewBufferString("")
		Logger.SetFile(output)
		Logger.SetLevel(ERROR)

		cmd := NewCli()
		cmd.SetArgs([]string{"../test/testing_vendor_index.vidx", "-f", "-o", outputFileName})
		ExitOnError(cmd.Execute())

		out, err := io.ReadAll(output)
		if err != nil {
			t.Fatal(err)
		}

		outStr := strings.TrimSpace(string(out))
		if len(outStr) == 0 || !strings.HasPrefix(outStr, "E: ") {
			t.Errorf("There should be an error log, instead got: '%s'", outStr)
		}

		out, err = os.ReadFile(outputFileName)
		if err != nil {
			t.Fatal(err)
		}

		expected := `<index>
 <vendor>testing_vendor_index</vendor>
 <url>file://XX/test-continue-despite-errors.xml</url>
 <>
 <pindex>
  <pdsc vendor="TheVendor" url="test/" name="ThePack2" version="1.1.0" timestamp=""></pdsc>
  <pdsc vendor="TheVendor" url="test/" name="ThePack1" version="1.2.3" timestamp=""></pdsc>
  <pdsc vendor="TheOtherVendor" url="test/" name="TheOtherPack" version="1.0.50" timestamp=""></pdsc>
  <pdsc vendor="TheVendor" url="non-existing-path/" name="ThePack" version="1.0.1" timestamp=""></pdsc>
 </pindex>
</index>`
		wd, _ := os.Getwd()
		wd = filepath.ToSlash(wd)
		expected = strings.ReplaceAll(expected, "XX", wd)
		s := string(out)
		sa, se, _ := strings.Cut(s, "timestamp") // cut out time, cannot compare
		_, se, _ = strings.Cut(se, "timestamp")
		AssertEqual(t, sa+se, expected)

		Logger.SetFile(currLogFile)
		Logger.SetLevel(currLevel)
		ExitOnError(os.RemoveAll(outputFileName))
	})
}
