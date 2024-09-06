package requests

import (
	"errors"
	"gurl/data_models"
	"gurl/requests/help"
	"gurl/requests/list_commands"
	"strings"

	"github.com/76creates/stickers"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/exp/slices"
)

type Model struct {
	items   list.Model
	folders []string
	help    help.Model
	flexbox stickers.FlexBox
}

func New(folder string) Model {
	if !strings.HasSuffix(folder, "/") {
		folder = folder + "/"
	}

	model := Model{
		items:   list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0),
		folders: []string{folder},
		help:    help.New(),
		flexbox: *stickers.NewFlexBox(0, 0),
	}

	model.items.Title = "Requests"
	model.items.SetFilteringEnabled(false)
	model.items.SetShowStatusBar(false)
	model.items.SetShowHelp(false)

	return model
}

func (self *Model) selectedRequestFullPath() string {
	if selectedRequets, err := self.selectedRequest(); err == nil {
		return self.currentFolder() + selectedRequets.Name + ".hurl"
	}

	return ""
}

func (self *Model) selectedRequest() (request datamodels.Request, err error) {
	if selectedRequest, isRequest := self.items.SelectedItem().(datamodels.Request); isRequest {
		return selectedRequest, nil
	}

	return datamodels.Request{}, errors.New("no requests")
}

func (self *Model) currentFolder() string {
	return self.folders[len(self.folders)-1]
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

		self.flexbox.SetWidth(msg.Width)
		self.flexbox.SetHeight(msg.Height)

		return self, self.readAllRequestsFromCurrentFolder

	case tea.KeyMsg:
		switch {
		case msg.String() == "enter":
			return self, handleEnter(&self)
		case key.Matches(msg, self.help.Keys.Up), key.Matches(msg, self.help.Keys.Down):
			cmds = append(cmds, self.changeRequest)
		case key.Matches(msg, self.help.Keys.Back):
			return self, handleGoBack(&self)
		case key.Matches(msg, self.help.Keys.Edit):
			return self, handleEdit(&self)
		}

	case AllRequestRead:
		// TODO: this is called on WindowSizeMsg, so maybe check if some items exist and overwrite them or something
		requests := []datamodels.Request(msg)
		mappedRequests := make([]list.Item, 0)

		slices.SortFunc(requests, func(f datamodels.Request, s datamodels.Request) int {
			if f.IsFolder && !s.IsFolder {
				return -1
			}

			if !f.IsFolder && s.IsFolder {
				return 1
			}

			return 0
		})

		for _, request := range requests {
			mappedRequests = append(mappedRequests, request)
		}

		self.items.SetItems(mappedRequests)

		cmds = append(cmds, self.changeRequest)
	}

	self.items, cmd = self.items.Update(msg)
	cmds = append(cmds, cmd)

	return self, tea.Batch(cmds...)
}

func (self Model) View() string {
	self.flexbox.AddRows([]*stickers.FlexBoxRow{
		self.rowWithContent(self.items.View(), 9),
		self.rowWithContent(self.help.View(), 1),
	})

	return self.flexbox.Render()
}

func (self *Model) rowWithContent(content string, ratioY int) *stickers.FlexBoxRow {
	return self.flexbox.
		NewRow().
		AddCells([]*stickers.FlexBoxCell{
			stickers.NewFlexBoxCell(1, ratioY).SetContent(content),
		})
}

type AllRequestRead []datamodels.Request

func (self *Model) readAllRequestsFromCurrentFolder() tea.Msg {
	return AllRequestRead(listcommands.ReadRequestsInfo(self.currentFolder()))
}

type RequestChanged struct {
	RequestFilePath string
	IsFolder        bool
}

func (self *Model) changeRequest() tea.Msg {
	selectedRequest, err := self.selectedRequest()

	if err != nil {
		return nil
	}

	return RequestChanged{
		RequestFilePath: self.selectedRequestFullPath(),
		IsFolder:        selectedRequest.IsFolder,
	}
}
