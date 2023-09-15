package requestcontent

import (
	"gurl/requests"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
    content viewport.Model
}

func New() Model {
    return Model {
        content: viewport.New(0, 0),
    }
}

func (self Model) Init() tea.Cmd {
    return nil
}

func (self Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd
    var cmd tea.Cmd

    switch msg := msg.(type) {

    case tea.WindowSizeMsg:
        self.content.Width = msg.Width
        self.content.Height = msg.Height

    case requests.RequestChanged:
        return self, self.readRequestContent(string(msg))

    case RequestRead:
        self.content.SetContent(string(msg))
    }

    self.content, cmd = self.content.Update(msg)
    cmds = append(cmds, cmd)

    return self, tea.Batch(cmds...)
}

func (self Model) View() string {
    return self.content.View()
}

type RequestRead string
func (self *Model) readRequestContent(requestFilePath string) tea.Cmd {
    return func() tea.Msg {
        return RequestRead(readRequestContent(requestFilePath))
    }
}
