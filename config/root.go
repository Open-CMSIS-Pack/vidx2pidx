package config


//
//  Handle all needed configuration for CMPack tool
//
//  Order of configurations to load:
//  1. From cmdline argument
//  2. From environment variable
//  3. From configuration file
//
//  This means that cmpack will first try to load configs from a file,
//  then overwrite it if there are any environment variables, and finally
//  overwrite again if there's any cmdline given regarding config
//
//  The config file should be in JSON format, as it's supported by default by GO's stdlib
//
//  Example of a config file:
//
//  $ cat .cmpack_config.json
//  {
//      "vidx_sources": [
//          {"name": "vidx-source-1", "path": "/path/to/vidx"},
//          {"name": "vidx-source-2", "path": "http://example.com/other.vidx"}
//      ]
//  }


import (
    "errors"
    "encoding/json"
    "io/ioutil"
    "fmt"
    "os"
)


var configFilename = ".cmpack.json"


//
//  Defines config file structure
//
type config struct {
    VidxSources []Vidx `json:"vidx_sources"`
}


//
//  Vidx stands for Vendor Index, it contains the sources
//  provided by vendors. Each vidx file can be obtained from
//  a local file or from a web URL.
//
type Vidx struct {
    Name string `json:"name"`
    Path string `json:"path"`
}


var Config config


func (c *config) write() {
    file, _ := json.MarshalIndent(*c, "", " ")
	_ = ioutil.WriteFile(configFilename, file, 0644)
}


func (c *config) vidxExist(name string) int {
    for i := 0; i < len(c.VidxSources); i++ {
        if name == c.VidxSources[i].Name {
            return i
        }
    }
    return -1
}


func AddVidx(name, path string) error {
    if exist := Config.vidxExist(name); exist != -1 {
        message := fmt.Sprintf("There is already a vidx for name '%s'", name)
        return errors.New(message)
    }
    vidx := Vidx{
        Name: name,
        Path: path,
    }

    Config.VidxSources = append(Config.VidxSources, vidx)
    Config.write()

    return nil
}


func ListVidxs() []Vidx {
    return Config.VidxSources
}


func RemoveVidx(name string) error {
    exist := Config.vidxExist(name)
    if exist == -1 {
        message := fmt.Sprintf("Vidx '%s' does not seem to exist", name)
        return errors.New(message)
    }

    Config.VidxSources = append(Config.VidxSources[:exist], Config.VidxSources[exist + 1:]...)
    Config.write()
    return nil
}


//
//  Read config file into memory
//
func Init() {

    jsonFile, err := os.Open(configFilename)
    if err != nil {
        fmt.Printf("W: Config file '%s' not found, creating one\n", configFilename)
        Config.write()
    }

    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    json.Unmarshal(byteValue, &Config)
}
