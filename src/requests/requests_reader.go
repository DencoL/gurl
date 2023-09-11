package requests

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"golang.org/x/exp/slices"
)

var supportedHttpMethods = []string {
    "GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH",
}

var readSize = len(findLongestSupportedHttpMethod())

func findLongestSupportedHttpMethod() string {
    return slices.MaxFunc[[]string, string](supportedHttpMethods, func(a string, b string) int { 
        return len(a) - len(b)
    })
}

// Currently reads only from the first level, subfolders will be added later
func ReadRequestsInfo(folderPath string) []Request {
    folderItems, err := os.ReadDir(folderPath)

    if err != nil {
        log.Fatal(err)
    }

    if len(folderItems) == 0 {
        return make([]Request, 0)
    }

    var result []Request
    filepath.WalkDir(folderPath, func(fullFilePath string, dirEntry fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if filepath.Ext(dirEntry.Name()) == ".hurl" {
            firstLine := readFirstLine(fullFilePath)
            httpMethod := parseHttpMethod(firstLine)
            
            result = append(result, Request{
                Name: strings.Split(dirEntry.Name(), ".")[0],
                Method: httpMethod,
            })
        }

        return nil
    })

    return result
}

func readFirstLine(fullFilePath string) string {
    file, err := os.Open(fullFilePath)
    if err != nil {
        log.Fatal("Error opening file!!!")
    }

    defer file.Close()
    bytes := make([]byte, readSize)

    actuallyReadSize, err := file.Read(bytes)
    if err != nil {
        log.Fatal(err)
    }

    return string(bytes[:actuallyReadSize])
}

func parseHttpMethod(value string) string {
    httpMethod := strings.Split(value, " ")[0]

    if slices.Contains(supportedHttpMethods, httpMethod) {
        return httpMethod
    }

    return ""
}
