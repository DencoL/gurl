package requests

import (
	"bufio"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"golang.org/x/exp/slices"
)

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
                Name: dirEntry.Name(),
                Method: httpMethod,
            })
        }

        return nil
    })

    return result
}

func readFirstLine(fullFilePath string) string {
    readFile, err := os.Open(fullFilePath)

    if err != nil {
        log.Fatal(err)
    }

    fileScanner := bufio.NewScanner(readFile)
    fileScanner.Split(bufio.ScanLines)
    var firstLine string

    for fileScanner.Scan() {
        firstLine = fileScanner.Text()
        break
    }

    readFile.Close()

    if err := fileScanner.Err(); err != nil {
        log.Fatal(err)
    }

    return firstLine
}

func parseHttpMethod(value string) string {
    supportedHttpMethods := []string {
        "GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH",
    }

    httpMethod := strings.Split(value, " ")[0]

    if slices.Contains(supportedHttpMethods, httpMethod) {
        return httpMethod
    }

    return ""
}
