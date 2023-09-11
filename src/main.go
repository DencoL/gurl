package main

import (
	"fmt"
	"gurl/requests"
	"os"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
    homeFolder := os.Getenv("HOME")
    requests := requests.ReadRequestsInfo(homeFolder + "/hurl-requests")

    program := tea.NewProgram(NewTui(requests))

    if _, err := program.Run(); err != nil {
        fmt.Println(err)
        os.Exit(1);
    }
}
