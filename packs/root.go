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
    "github.com/chaws/cmpack/config"
    "os"
)


var listsDir = "cmpack-lists"


func updateVidxList(vidx config.Vidx) {
    fmt.Printf("I: Updating packs list for '%s' (reading from %s)\n", vidx.Name, vidx.Path)
    packsList := ReadPacksList(vidx.Name)
    fmt.Printf("Pack %s", packsList.Name)
    // Read vidx.Name list file into memory
    // Read pidx-es from vidx.Path (either local file or url)
    // If there's any new package or package removal, update local list
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
