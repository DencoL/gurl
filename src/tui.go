package main

import (
	"gurl/content"
	"gurl/requests"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewTui(requests []requests.Request, requestsFolderPath string) Model {
    return Model {
        requests: requests,
        requestContent: content.NewContentModel(requestsFolderPath),
    }
}

type Model struct {
    requestsFolderPath string
    requests []requests.Request
    requestsList list.Model
    requestContent content.Model
}

func (self *Model) SelectedRequest() requests.Request {
    return self.requests[self.requestsList.Cursor()]
}

func (self *Model) initRequests(width int, height int) {
    self.requestsList = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)

    self.requestsList.SetFilteringEnabled(false)
    self.requestsList.SetShowStatusBar(false)
    self.requestsList.Title = "Requests"

    self.requestsList.SetItems(make([]list.Item, len(self.requests)))
    for index, request := range self.requests {
        self.requestsList.SetItem(index, request)
    }
}

func (self Model) Init() tea.Cmd {
    return nil
}

func (self Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        self.initRequests(msg.Width, msg.Height)

        self.requestContent.SetDimensions(msg.Width, msg.Height)
    }

    var listCmd tea.Cmd
    self.requestsList, listCmd = self.requestsList.Update(msg)

    self.requestContent.SetContent(self.SelectedRequest().Name)

    return self, listCmd
}

func (self Model) View() string {
    return lipgloss.JoinHorizontal(
        lipgloss.Left,
        self.requestsList.View(),
        self.requestContent.View(),
    )
}
