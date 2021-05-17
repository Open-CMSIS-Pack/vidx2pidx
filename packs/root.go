package packs


//
//  This module is responsible for updating the list,
//  adding, removing, searching and upgrading packs.
//
//  There are two main aspects for packs:
//
//  1. External packs organization:
//       - Packs are distributed by vendors, and each vendor
//         provides a *.pidx file listing all available packs.
//       - Each pack has a pack description file (*.pdsc) that
//         gather useful info regarding the pack, including a
//         timestamp of package last modified time and a list of
//         all file components in that pack.
//
//  2. Internal packs organization (local layout)
//       - The cmpack config file contains a list of vendor
//         index sources (*.vidx), so cmpack needs to generate
//         a list of packs from them.
//


import (
    "fmt"
    "os"
    "github.com/chaws/cmpack/config"
)


var listsDir = "cmpack-lists"


func updateVidxList(vidx config.Vidx) error {
    fmt.Printf("I: Updating packs list for '%s' (reading from %s)\n", vidx.Name, vidx.Path)
    //packsList := ReadPacksList(vidx.Name)

    vidxXML, err := ReadVidx(vidx.Path)
    if err != nil {
       return err
    }

    fmt.Printf("I: Reading package indexes from %s (%s)\n", vidxXML.Vendor, vidxXML.URL)
    for _, pidx := range vidxXML.Vindex.Pidxs {
        url := fmt.Sprintf("%s/%s.pidx", pidx.URL, pidx.Vendor)
        pidxXML, err := ReadPidx(url)
        if err != nil {
           return err
        }

        for _, pdsc := range pidxXML.Pindex.Pdscs {
            fmt.Println(pdsc.URL)
        }
    }
    return nil
}


//
//  Get the vidx list and generate a list of packs
//  from each vidx source.
//
func UpdateList() error {
    fmt.Println("Updating list of packs")
    vidxSources := config.ListVidxs()
    if len(vidxSources) == 0 {
        fmt.Println("I: Nothing to do, no vendor index found.")
        return nil
    }

    for _, vidx := range vidxSources {
        updateVidxList(vidx)
    }
    return nil
}

func init() {
    // Makes sure lists directory exists
    _, err := os.Stat(listsDir)
    if os.IsNotExist(err) {
        os.Mkdir(listsDir, 0775)
    }
}
