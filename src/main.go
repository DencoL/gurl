package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {
	homeFolder := os.Getenv("HOME")
	requestsFolderPath := homeFolder + "/hurl-requests"

	program := tea.NewProgram(NewTui(requestsFolderPath))

	if _, err := program.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
