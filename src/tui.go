package main

import (
	"gurl/requests"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func NewTui(requests []requests.Request) Model {
    return Model {
        requests: requests,
    }
}

type Model struct {
    requests []requests.Request
    requestsList list.Model
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
    }

    var cmd tea.Cmd
    self.requestsList, cmd = self.requestsList.Update(msg)

    return self, cmd
}

func (self Model) View() string {
    return self.requestsList.View()
}
