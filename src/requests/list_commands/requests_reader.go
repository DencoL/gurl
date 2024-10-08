package listcommands

import (
	"gurl/data_models"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

var supportedHttpMethods = []string{
	"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH",
}

var readSize = len(findLongestSupportedHttpMethod())

func findLongestSupportedHttpMethod() string {
	return slices.MaxFunc[[]string, string](supportedHttpMethods, func(a string, b string) int {
		return len(a) - len(b)
	})
}

// Currently reads only from the first level, subfolders will be added later
func ReadRequestsInfo(folderPath string) []datamodels.Request {
	folderItems, err := os.ReadDir(folderPath)

	if err != nil {
		log.Fatal(err)
	}

	if len(folderItems) == 0 {
		return make([]datamodels.Request, 0)
	}

	var result []datamodels.Request
	filepath.WalkDir(folderPath, func(fullFilePath string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if fullFilePath == folderPath {
			return nil
		}

		if !isInTopLevelFolder(fullFilePath, folderPath) {
			return nil
		}

		if dirEntry.IsDir() {
			result = append(result, datamodels.Request{
				Name:     dirEntry.Name(),
				Method:   "",
				IsFolder: true,
			})

			return nil
		}

		if filepath.Ext(dirEntry.Name()) == ".hurl" {
			firstLine, _ := readFirstLine(fullFilePath)
			// if err == nil {
			httpMethod := parseHttpMethod(firstLine)

			result = append(result, datamodels.Request{
				Name:   strings.Split(dirEntry.Name(), ".")[0],
				Method: httpMethod,
			})
			// }
		}

		return nil
	})

	return result
}

func readFirstLine(fullFilePath string) (string, error) {
	file, err := os.Open(fullFilePath)
	if err != nil {
		log.Fatal("Error opening file!!!")
	}

	defer file.Close()
	bytes := make([]byte, readSize)

	actuallyReadSize, err := file.Read(bytes)

	return string(bytes[:actuallyReadSize]), err
}

func parseHttpMethod(value string) string {
	httpMethod := strings.Split(value, " ")[0]

	if slices.Contains(supportedHttpMethods, httpMethod) {
		return httpMethod
	}

	return ""
}

func isInTopLevelFolder(path string, expectedTopLevelFolderName string) bool {
	pathParts := strings.Split(path, "/")
	pathWithoutFile := strings.Join(pathParts[0:(len(pathParts)-1)], string(os.PathSeparator))

	return pathWithoutFile == expectedTopLevelFolderName || pathWithoutFile+"/" == expectedTopLevelFolderName
}
