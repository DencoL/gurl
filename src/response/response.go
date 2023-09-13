package response

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type HurlCommandDone string

type Model struct {
    content viewport.Model
}

func NewResponseModel(width int, height int) Model {
    viewPort := viewport.New(width, height)
    return Model {
        content: viewPort,
    }
}

func (self Model) Init() tea.Cmd {
    return nil
}

func (self Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg.(type) {
    case HurlCommandDone:
        self.content.SetContent("from hurl command done response")

        newContentModel, cmd := self.content.Update(msg)
        self.content = newContentModel
        return self, cmd
    }
    return self, nil
}

func (self Model) View() string {
    return self.content.View()
}

func (self *Model) SetDimensions(width int, height int) {
    self.content.Width = width
    self.content.Height = height
}

func (self *Model) SetContent(v string) {
    self.content.SetContent(v)
}
