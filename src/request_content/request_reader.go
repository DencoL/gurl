package requestcontent

import "os"

func readRequestContent(requestFilePath string) string {
	bytes, err := os.ReadFile(requestFilePath)

	if err != nil {
		return err.Error()
	}

	return string(bytes)
}
