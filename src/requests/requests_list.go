package requests

import (
	"os/exec"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
    items list.Model
    requestsFolderPath string
}

func NewRequestsList(requestsFolderPath string) Model {
    model := Model {
        items: list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0),
        requestsFolderPath: requestsFolderPath,
    }

    model.items.Title = "Requests"
    model.items.SetFilteringEnabled(false)
    model.items.SetShowStatusBar(false)

    return model
}

func (self *Model) selectedRequestFullPath() string {
    selectedRequest := self.items.SelectedItem().(Request)

    return self.requestsFolderPath + selectedRequest.Name + ".hurl"
}

func (self Model) Init() tea.Cmd {
    return nil
}

func (self Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd
    var cmd tea.Cmd

    switch msg := msg.(type) {

    case tea.WindowSizeMsg:
        self.items.SetSize(msg.Width, msg.Height)
        cmds = append(cmds, self.readAllRequestsFromCurrentFolder)

    case tea.KeyMsg:
        if msg.String() == "enter" {
            return self, self.executeRequest
        }

        switch {
            case key.Matches(msg, list.DefaultKeyMap().CursorDown):
                cmds = append(cmds, self.changeRequest)
            case key.Matches(msg, list.DefaultKeyMap().CursorUp):
                cmds = append(cmds, self.changeRequest)
        }

    case AllRequestRead:
        // TODO: this is called on WindowSizeMsg, so maybe check if some items exist and overwrite them or something
        requests := []Request(msg)
        mappedRequests := make([]list.Item, len(requests))

        for index, request := range requests {
            mappedRequests[index] = request
        }

        self.items.SetItems(mappedRequests)

        return self, self.changeRequest
    }

    self.items, cmd = self.items.Update(msg)
    cmds = append(cmds, cmd)

    return self, tea.Batch(cmds...)
}

func (self Model) View() string {
    return self.items.View()
}

type AllRequestRead []Request
func (self *Model) readAllRequestsFromCurrentFolder() tea.Msg {
    return AllRequestRead(ReadRequestsInfo(self.requestsFolderPath))
}

type RequestChanged string
func (self *Model) changeRequest() tea.Msg {
    return RequestChanged(self.selectedRequestFullPath())
}

type RequestExecuted string
func (self *Model) executeRequest() tea.Msg {
    hurl := exec.Command("hurl", self.selectedRequestFullPath())
    stdout, err := hurl.Output()

    if err != nil {
        return RequestExecuted(err.Error())
    }

    return RequestExecuted(stdout)
}
