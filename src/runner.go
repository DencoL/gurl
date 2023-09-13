package main

import (
	"os/exec"
)

func RunHurl(hurlFilePath string) string {
    hurl := exec.Command("hurl", hurlFilePath)
    stdout, err := hurl.Output()

    if err != nil {
        return err.Error()
    }

    return string(stdout)
}
