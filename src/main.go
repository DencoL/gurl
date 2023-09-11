package main

import (
	"fmt"
	"gurl/requests"
	"os"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
    homeFolder := os.Getenv("HOME")
    requestsFolderPath := homeFolder + "/hurl-requests"
    requests := requests.ReadRequestsInfo(requestsFolderPath)

    program := tea.NewProgram(NewTui(requests, requestsFolderPath))

    if _, err := program.Run(); err != nil {
        fmt.Println(err)
        os.Exit(1);
    }
}
