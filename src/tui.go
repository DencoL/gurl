package main

import (
	requestcontent "gurl/request_content"
	requestresponse "gurl/request_response"
	"gurl/requests"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewTui(requestsFolderPath string) Model {
    if (!strings.HasSuffix(requestsFolderPath, "/")) {
        requestsFolderPath = requestsFolderPath + "/"
    }

    model := Model {
        requestsFolderPath: requestsFolderPath,
        requestsList: requests.New(requestsFolderPath),
        requestContent: requestcontent.New(),
        requestResponse: requestresponse.New(),
    }

    return model
}

type Model struct {
    requestsFolderPath string
    requestsList requests.Model
    requestContent requestcontent.Model
    requestResponse requestresponse.Model
}

func (self Model) Init() tea.Cmd {
    return nil
}

func (self Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    res, cmd := self.requestsList.Update(msg)
    self.requestsList = res.(requests.Model)
    cmds = append(cmds, cmd)

    res, cmd = self.requestContent.Update(msg)
    self.requestContent = res.(requestcontent.Model)
    cmds = append(cmds, cmd)

    res, cmd = self.requestResponse.Update(msg)
    self.requestResponse = res.(requestresponse.Model)
    cmds = append(cmds, cmd)

    return self, tea.Batch(cmds...)
}

func (self Model) View() string {
    return lipgloss.JoinHorizontal(
        lipgloss.Left,
        self.style(50, self.requestsList.View()),
        self.style(100, self.requestContent.View()),
        self.style(100, self.requestResponse.View()),
    )
}

func (self *Model) style(width int, view string) string {
    return lipgloss.NewStyle().
        BorderStyle(lipgloss.NormalBorder()).
        BorderRight(true).
        Width(width).
        Render(view)
}
