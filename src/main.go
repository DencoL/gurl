package main

import (
	"fmt"
	"os"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
    homeFolder := os.Getenv("HOME")
    requestsFolderPath := homeFolder + "/hurl-requests"

    program := tea.NewProgram(NewTui(requestsFolderPath))

    if _, err := program.Run(); err != nil {
        fmt.Println(err)
        os.Exit(1);
    }
}
