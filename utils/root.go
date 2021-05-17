package utils


import (
    "net/http"
    "io/ioutil"
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
