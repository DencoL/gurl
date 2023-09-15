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
    cmds := []tea.Cmd {
        updateSubComponent[requests.Model](&self.requestsList, msg),
        updateSubComponent[requestcontent.Model](&self.requestContent, msg),
        updateSubComponent[requestresponse.Model](&self.requestResponse, msg),
    }

    return self, tea.Batch(cmds...)
}

func updateSubComponent[TModel tea.Model](model *TModel, msg tea.Msg) tea.Cmd {
    res, cmd := (*model).Update(msg)
    *model = res.(TModel)

    return cmd
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
