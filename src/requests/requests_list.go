package requests

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
    items list.Model
    requestsFolderPath string
}

func New(requestsFolderPath string) Model {
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
    if len(self.items.Items()) == 0 {
        return ""
    }

    return self.requestsFolderPath + self.selectedRequest().Name + ".hurl"
}

func (self *Model) selectedRequest() Request {
    return self.items.SelectedItem().(Request)
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
        return self, self.readAllRequestsFromCurrentFolder

    case tea.KeyMsg:
        switch {
            case msg.String() == "enter":
                if self.selectedRequest().IsFolder {
                    return self, nil
                } else {
                    return self, self.executeRequest
                }
            case key.Matches(msg, list.DefaultKeyMap().CursorDown), key.Matches(msg, list.DefaultKeyMap().CursorUp):
                cmds = append(cmds, self.changeRequest)
        }

    case AllRequestRead:
        // TODO: this is called on WindowSizeMsg, so maybe check if some items exist and overwrite them or something
        requests := []Request(msg)
        mappedRequests := make([]list.Item, 0)

        for _, request := range requests {
            if !request.IsFolder {
                continue
            }

            mappedRequests = append(mappedRequests, request)
        }

        for _, request := range requests {
            if request.IsFolder {
                continue
            }

            mappedRequests = append(mappedRequests, request)
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

type RequestChanged struct {
    RequestFilePath string
    IsFolder bool
}

func (self *Model) changeRequest() tea.Msg {
    return RequestChanged {
        RequestFilePath: self.selectedRequestFullPath(),
        IsFolder: self.selectedRequest().IsFolder,
    }
}

type ExecuteRequest string
func (self *Model) executeRequest() tea.Msg {
    return ExecuteRequest(self.selectedRequestFullPath())
}
