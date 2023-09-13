package main

import (
	"gurl/requests"
	"gurl/response"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewTui(requests []requests.Request, requestsFolderPath string) Model {
    if (!strings.HasSuffix(requestsFolderPath, "/")) {
        requestsFolderPath = requestsFolderPath + "/"
    }

    model := Model {
        requestsFolderPath: requestsFolderPath,
        requestContent: content.NewContentModel(requestsFolderPath),
    }

    model.requestsList = list.New(make([]list.Item, len(requests)), list.NewDefaultDelegate(), 0, 0)
    model.requestsList.Title = "Requests"
    model.requestsList.SetFilteringEnabled(false)
    model.requestsList.SetShowStatusBar(false)

    for index, request := range requests {
        model.requestsList.SetItem(index, request)
    }

    return model
}

type Model struct {
    requestsFolderPath string
    requestsList list.Model
    requestContent content.Model
    response viewport.Model
}

func (self *Model) selectedRequestFullPath() string {
    selectedRequest := self.requestsList.SelectedItem().(requests.Request)

    return self.requestsFolderPath + selectedRequest.Name + ".hurl"
}

func (self *Model) setListDimensions(width int, height int) {
    self.requestsList.SetWidth(width)
    self.requestsList.SetHeight(height)
}

func (self Model) Init() tea.Cmd {
    return nil
}

func (self Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        self.setListDimensions(msg.Width, msg.Height)
        self.requestContent.SetDimensions(msg.Width, msg.Height)
        self.response.Width = msg.Width
        self.response.Height = msg.Height
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyEnter:
            return self, self.runHurlCommand
        }
    case response.HurlCommandDone:
        self.response.SetContent(string(msg))
    }

    var listCmd tea.Cmd
    self.requestsList, listCmd = self.requestsList.Update(msg)

    self.requestContent.SetContent(self.selectedRequestFullPath())

    return self, listCmd
}

func (self *Model) runHurlCommand() tea.Msg {
    return response.HurlCommandDone(RunHurl(self.selectedRequestFullPath()))
}

func (self Model) View() string {
    return lipgloss.JoinHorizontal(
        lipgloss.Left,
        self.requestsList.View(),
        self.requestContent.View(),
        self.response.View(),
    )
}
