package utils


import (
    "net/http"
    "io/ioutil"
    "os"
    "strings"
    "encoding/xml"
)


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

    xml.Unmarshal(contents, targetStruct)

    return nil
}
