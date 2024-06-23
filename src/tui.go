package main

import (
	requestcontent "gurl/request_content"
	requestresponse "gurl/request_response"
	"gurl/requests"
	"strings"
	"github.com/76creates/stickers"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewTui(requestsFolderPath string) Model {
    if (!strings.HasSuffix(requestsFolderPath, "/")) {
        requestsFolderPath = requestsFolderPath + "/"
    }

    model := Model {
        flexbox: *stickers.NewFlexBox(0, 0),
        requestsFolderPath: requestsFolderPath,
        requestsList: requests.New(requestsFolderPath),
        requestContent: requestcontent.New(),
        requestResponse: requestresponse.New(),
    }

    return model
}

type Model struct {
    flexbox stickers.FlexBox
    requestsFolderPath string
    requestsList requests.Model
    requestContent requestcontent.Model
    requestResponse requestresponse.Model
}

func (self Model) Init() tea.Cmd {
    return nil
}

func (self Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        self.flexbox.SetWidth(msg.Width)
        self.flexbox.SetHeight(msg.Height)
    }

    cmds := []tea.Cmd {
        updateSubComponent[requests.Model](&self.requestsList, msg),
        updateSubComponent[requestcontent.Model](&self.requestContent, msg),
        updateSubComponent[requestresponse.Model](&self.requestResponse, msg),
    }

    return self, tea.Batch(cmds...)
}

func updateSubComponent[TModel tea.Model](model *TModel, msg tea.Msg) tea.Cmd {
    res, cmd := (*model).Update(msg)
    *model = res.(TModel)

    return cmd
}

func (self Model) View() string {
    mainRow := self.flexbox.
        NewRow().
        AddCells([]*stickers.FlexBoxCell {
            createRightBorderedCell(1, self.requestsList.View()),
            createRightBorderedCell(2, self.requestContent.View()),
            createCell(2, self.requestResponse.View()),
        })

    self.flexbox.AddRows([]*stickers.FlexBoxRow {
        mainRow,
    })

    return self.flexbox.Render()
}

func createCell(ratioX int, content string) *stickers.FlexBoxCell  {
    return stickers.
        NewFlexBoxCell(ratioX, 1).
        SetContent(content)
}

func createRightBorderedCell(ratioX int, content string) *stickers.FlexBoxCell  {
    return createCell(ratioX, content).
        SetStyle(lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderRight(true))
}

func (self *Model) style(width int, view string) string {
    return lipgloss.NewStyle().
        BorderStyle(lipgloss.NormalBorder()).
        BorderRight(true).
        Render(view)
}
