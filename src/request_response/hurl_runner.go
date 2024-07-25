package requestresponse

import (
	"os/exec"
)

func runHurl(hurlFilePath string) string {
	hurl := exec.Command("hurl", "--insecure", "--verbose", hurlFilePath)
	stdout, err := hurl.Output()

	if err != nil {
		return err.Error()
	}

	return string(stdout)
}
