package requests

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type RequestInfo struct {
    name string
    method string
}

func (self *RequestInfo) GetName() string {
    return self.name
}

func (self *RequestInfo) GetMethod() string {
    return self.method
}

// Currently reads only from the first level, subfolders will be added later
func ReadRequestsInfo(folderPath string) []RequestInfo {
    folderItems, err := os.ReadDir(folderPath)

    if err != nil {
        log.Fatal(err)
    }

    if len(folderItems) == 0 {
        return make([]RequestInfo, 0)
    }

    var result []RequestInfo
    filepath.WalkDir(folderPath, func(_ string, dirEntry fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if filepath.Ext(dirEntry.Name()) == ".hurl" {
            result = append(result, RequestInfo{
                name: dirEntry.Name(),
                method: "",
            })
        }

        return nil
    })

    return result
}
