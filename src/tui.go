package main

import (
	"gurl/requests"
	"log"
	"os"
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

    return Model {
        requestsFolderPath: requestsFolderPath,
        requests: requests,
    }
}

type Model struct {
    requestsFolderPath string
    requests []requests.Request
    requestsList list.Model
    requestContent viewport.Model
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

        self.requestContent.Width = msg.Width
        self.requestContent.Height = msg.Height
    }

    var listCmd tea.Cmd
    self.requestsList, listCmd = self.requestsList.Update(msg)

    self.requestContent.SetContent(self.readSelectedRequestContent())
    var contentCmd tea.Cmd
    self.requestContent, contentCmd = self.requestContent.Update(msg)

    return self, tea.Batch(listCmd, contentCmd)
}

func (self *Model) readSelectedRequestContent() string {
    selectedRequestName := self.requestsFolderPath + self.requests[self.requestsList.Cursor()].Name + ".hurl"
    bytes, err := os.ReadFile(selectedRequestName)

    if err != nil {
        log.Fatal(err)
    }

    return string(bytes)
}

func (self Model) View() string {
    return lipgloss.JoinHorizontal(
        lipgloss.Left,
        self.requestsList.View(),
        self.requestContent.View(),
    )
}
