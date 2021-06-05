package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
		out, err = ioutil.ReadAll(output)
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
		out, err = ioutil.ReadAll(output)
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("vidx2pidx version %v\n%v", Version, License)

		AssertEqual(t, strings.TrimSpace(string(out)), expected)
	})
}

func ExampleSampleOutput() {

	cmd := NewCli()
	cmd.SetArgs([]string{"test/cypress.vidx", "-o", "-"})
	cmd.Execute()
	// Output:
	// <index>
	//  <timestamp></timestamp>
	//  <pindex>
	//   <pdsc vendor="Cypress" url="https://github.com/cypresssemiconductorco/cmsis-packs/raw/master/PSoC6_DFP/" name="PSoC6_DFP" version="1.2.0" timestamp=""></pdsc>
	//   <pdsc vendor="Cypress" url="https://github.com/cypresssemiconductorco/cmsis-packs/raw/master/PSoC4_DFP/" name="PSoC4_DFP" version="1.1.0" timestamp=""></pdsc>
	//   <pdsc vendor="Atmel" url="http://packs.download.atmel.com/" name="SAM3A_DFP" version="1.0.50" timestamp=""></pdsc>
	//  </pindex>
	// </index>
}
