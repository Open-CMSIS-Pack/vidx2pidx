/* SPDX-License-Identifier: Apache-2.0 */
/* Copyright Contributors to the vidx2pidx project. */

package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

func AnyErr(errs []error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

func ExitOnError(err error) {
	if err != nil {
		Logger.Error(err.Error())
		os.Exit(-1)
	}
}

var CacheDir string

func ReadURL(url string) ([]byte, error) {
	var empty []byte
	resp, err := http.Get(url) // #nosec
	if err != nil {
		return empty, err
	}

	if resp.StatusCode/100 != 2 {
		message := fmt.Sprintf("request did not return successfully (%v)", resp.StatusCode)
		return empty, errors.New(message)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return empty, err
	}

	if len(CacheDir) > 0 {
		fileName := path.Join(CacheDir, path.Base(url))
		err = ioutil.WriteFile(fileName, body, 0600)
		if err != nil {
			return body, err
		}
	}

	return body, nil
}

func ReadXML(path string, targetStruct interface{}) error {
	var contents []byte
	var err error
	var xmlFile *os.File

	if strings.HasPrefix(path, "http") {
		contents, err = ReadURL(path)
		if err != nil {
			return err
		}
	} else {
		xmlFile, err = os.Open(path)
		if err != nil {
			return err
		}

		contents, err = ioutil.ReadAll(xmlFile)
		if err != nil {
			return err
		}
	}

	if err = xml.Unmarshal(contents, targetStruct); err != nil {
		return err
	}

	return nil
}

func WriteXML(path string, targetStruct interface{}) error {
	output, err := xml.MarshalIndent(targetStruct, "", " ")
	if err != nil {
		return err
	}

	if path == "" || path == "-" {
		os.Stdout.Write(output)
		return nil
	}

	err = ioutil.WriteFile(path, output, 0600)
	if err != nil {
		return err
	}

	return nil
}

func EnsureDir(dirName string) error {
	err := os.MkdirAll(dirName, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}
