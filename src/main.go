package main

import (
	"fmt"
	"gurl/requests"
	"os"
)

func main() {
    homeFolder := os.Getenv("HOME")
    requests := requests.ReadRequestsInfo(homeFolder + "/hurl-requests")

    // Call to TUI will be here later
    for _, r := range requests {
        fmt.Println(r.Name)
    }
}
