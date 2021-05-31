package main


import (
    "fmt"
    "net/http"
    "io/ioutil"
    "os"
    "strings"
    "encoding/xml"
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
        fmt.Fprintf(os.Stderr, "E: %s\n", err)
        os.Exit(-1)
    }
}


func ReadURL(URL string) ([]byte, error) {

    var empty []byte
    resp, err := http.Get(URL)
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


func ReadXML(path string,  targetStruct interface{}) error {

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


func WriteXML(path string,  targetStruct interface{}) error {

    output, err := xml.MarshalIndent(targetStruct, "", " ")
    if err != nil {
        return err
    }

    if path == "" {
        os.Stdout.Write(output)
        return nil
    }

	err = ioutil.WriteFile(path, output, 0666)
    if err != nil {
        return err
    }

    return nil
}
