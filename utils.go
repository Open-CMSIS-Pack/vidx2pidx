package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"os"
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

func ReadURL(url string) ([]byte, error) {
	var empty []byte
	resp, err := http.Get(url)
	if err != nil {
		return empty, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return empty, err
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

	err = ioutil.WriteFile(path, output, 0666)
	if err != nil {
		return err
	}

	return nil
}
