package requests

import (
	"errors"
	"gurl/data_models"
	"gurl/requests/list_commands"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/exp/slices"
)

type listKeyMap struct {
	selectRequest key.Binding
	changeRequest key.Binding
    back key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		selectRequest: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("Enter", "Execute/Open"),
		),
		changeRequest: key.NewBinding(
			key.WithKeys("j", "k"),
			key.WithHelp("j/k", "Down/Up"),
		),
		back: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "Go back"),
		),
	}
}

type Model struct {
    items list.Model
    folders []string
    keys *listKeyMap
}

func New(folder string) Model {
    if !strings.HasSuffix(folder, "/") {
        folder = folder + "/"
    }

    model := Model {
        items: list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0),
        folders: []string { folder },
    }

    model.keys = newListKeyMap()
    // TODO: the help menu looks terrible now
    model.items.AdditionalShortHelpKeys = func() []key.Binding {
        return []key.Binding {
            model.keys.selectRequest,
            model.keys.back,
        }
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

    return self.currentFolder() + selectedRequets.Name + ".hurl"
}

func (self *Model) selectedRequest() (request datamodels.Request, err error)  {
    selectedRequest, isRequest := self.items.SelectedItem().(datamodels.Request)

    if !isRequest {
        return datamodels.Request{}, errors.New("no requests")
    }

    return selectedRequest, nil
}

func (self *Model) currentFolder() string {
    return self.folders[len(self.folders) - 1]
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
            case key.Matches(msg, self.keys.selectRequest):
                selectedRequest, err := self.selectedRequest()
                if err != nil {
                    return self, nil
                }

                if selectedRequest.IsFolder {
                    nextFolder := self.currentFolder() + selectedRequest.Name + "/"
                    self.folders = append(self.folders, nextFolder)

                    return self, self.readAllRequestsFromCurrentFolder
                } else {
                    return self, self.executeRequest
                }
            case key.Matches(msg, self.keys.changeRequest):
                cmds = append(cmds, self.changeRequest)
            case key.Matches(msg, self.keys.back):
                foldersLength := len(self.folders)
                
                if foldersLength == 1 {
                    return self, nil
                }
                
                self.folders = slices.Delete(self.folders, foldersLength - 1, foldersLength)
                
                return self, self.readAllRequestsFromCurrentFolder
        }

    case AllRequestRead:
        // TODO: this is called on WindowSizeMsg, so maybe check if some items exist and overwrite them or something
        requests := []datamodels.Request(msg)
        mappedRequests := make([]list.Item, 0)

        slices.SortFunc(requests, func(f datamodels.Request, s datamodels.Request) int {
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

type AllRequestRead []datamodels.Request
func (self *Model) readAllRequestsFromCurrentFolder() tea.Msg {
    return AllRequestRead(listcommands.ReadRequestsInfo(self.currentFolder()))
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
