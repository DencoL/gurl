package content

import (
	"log"
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
    content viewport.Model
}

func NewContentModel(requestsFolderPath string) Model {
    return Model{}
}

func (self Model) Init() tea.Cmd {
    return nil
}

func (self Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    newContentModel, cmd := self.content.Update(msg)
    self.content = newContentModel

    return self, cmd
}

func (self Model) View() string {
    return self.content.View()
}

func (self *Model) SetDimensions(width int, height int) {
    self.content.Width = width
    self.content.Height = height
}

func (self *Model) SetContent(requestFullPath string) {
    self.content.SetContent(self.readRequestContent(requestFullPath))
}

func (self *Model) readRequestContent(fullFilePath string) string {
    bytes, err := os.ReadFile(fullFilePath)

    if err != nil {
        log.Fatal(err)
    }

    return string(bytes)
}
