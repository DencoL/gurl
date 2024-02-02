package requests

import (
	"errors"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/exp/slices"
)

type Model struct {
    items list.Model
    currentFolder string
}

func New(folder string) Model {
    if !strings.HasSuffix(folder, "/") {
        folder = folder + "/"
    }

    model := Model {
        items: list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0),
        currentFolder: folder,
    }

    model.items.Title = "Requests"
    model.items.SetFilteringEnabled(false)
    model.items.SetShowStatusBar(false)

    return model
}

func (self *Model) selectedRequestFullPath() string {
    selectedRequets, err := self.selectedRequest()
    if err != nil {
        return ""
    }

    return self.currentFolder + selectedRequets.Name + ".hurl"
}

func (self *Model) selectedRequest() (request Request, err error)  {
    selectedRequest, isRequest := self.items.SelectedItem().(Request)

    if !isRequest {
        return Request{}, errors.New("no requests")
    }

    return selectedRequest, nil
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
                selectedRequest, err := self.selectedRequest()
                if err != nil {
                    return self, nil
                }

                if selectedRequest.IsFolder {
                    self.currentFolder = self.currentFolder + selectedRequest.Name + "/"
                    return self, self.readAllRequestsFromCurrentFolder
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

        slices.SortFunc(requests, func(f Request, s Request) int {
            if (f.IsFolder && !s.IsFolder) {
                return -1
            }

            if (!f.IsFolder && s.IsFolder) {
                return 1
            }

            return 0
        });

        for _, request := range requests {
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
    return AllRequestRead(ReadRequestsInfo(self.currentFolder))
}

type RequestChanged struct {
    RequestFilePath string
    IsFolder bool
}

func (self *Model) changeRequest() tea.Msg {
    selectedRequest, err := self.selectedRequest()

    if err != nil {
        return nil
    }

    return RequestChanged {
        RequestFilePath: self.selectedRequestFullPath(),
        IsFolder: selectedRequest.IsFolder,
    }
}

type ExecuteRequest string
func (self *Model) executeRequest() tea.Msg {
    return ExecuteRequest(self.selectedRequestFullPath())
}
