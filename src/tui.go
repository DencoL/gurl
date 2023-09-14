package main

import (
	"gurl/requests"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewTui(requestsFolderPath string) Model {
    if (!strings.HasSuffix(requestsFolderPath, "/")) {
        requestsFolderPath = requestsFolderPath + "/"
    }

    model := Model {
        requestsFolderPath: requestsFolderPath,
        requestsList: requests.NewRequestsList(requestsFolderPath),
    }

    return model
}

type Model struct {
    requestsFolderPath string
    requestsList requests.Model
    requestContent viewport.Model
    response viewport.Model
}

func (self Model) Init() tea.Cmd {
    return nil
}

func (self Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    switch msg := msg.(type) {

    case tea.WindowSizeMsg:
        self.requestContent.Width = msg.Width
        self.requestContent.Height = msg.Height

        self.response.Width = msg.Width
        self.response.Height = msg.Height

    case requests.RequestRead:
        self.requestContent.SetContent(string(msg))

    case requests.RequestExecuted:
        res := string(msg) 
        if res == "" {
            res = "<EMPTY RESPONSE>"
        }
        self.response.SetContent(res)
    }

    res, cmd := self.requestsList.Update(msg)
    self.requestsList = res.(requests.Model)
    cmds = append(cmds, cmd)

    return self, tea.Batch(cmds...)
}

func (self Model) View() string {
    return lipgloss.JoinHorizontal(
        lipgloss.Left,
        self.style(50, self.requestsList.View()),
        self.style(100, self.requestContent.View()),
        self.style(100, self.response.View()),
    )
}

func (self *Model) style(width int, view string) string {
    return lipgloss.NewStyle().
        BorderStyle(lipgloss.NormalBorder()).
        BorderRight(true).
        Width(width).
        Render(view)
}
