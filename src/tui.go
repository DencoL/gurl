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

    model := Model {
        requestsFolderPath: requestsFolderPath,
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
    requestContent viewport.Model
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

        self.requestContent.Width = msg.Width
        self.requestContent.Height = msg.Height

        self.response.Width = msg.Width
        self.response.Height = msg.Height

    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyEnter:
            return self, self.runHurlCommand
        case tea.KeyDown:
        case tea.KeyUp:
            return self, self.readRequestContent
        }
    case HurlCommandDone:
        self.response.SetContent(string(msg))
    case RequestRead:
        self.requestContent.SetContent(string(msg))
    }

    var listCmd tea.Cmd
    self.requestsList, listCmd = self.requestsList.Update(msg)

    return self, tea.Batch(listCmd, self.readRequestContent)
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

type HurlCommandDone string
func (self *Model) runHurlCommand() tea.Msg {
    return HurlCommandDone(RunHurl(self.selectedRequestFullPath()))
}

type RequestRead string
func (self *Model) readRequestContent() tea.Msg {
    bytes, err := os.ReadFile(self.selectedRequestFullPath())

    if err != nil {
        log.Fatal(err)
    }

    return RequestRead(bytes)
}
