package content

import (
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
    content viewport.Model
    requestsFolderPath string
}

type ContentRead string

func NewContentModel(requestsFolderPath string) Model {
    if (!strings.HasSuffix(requestsFolderPath, "/")) {
        requestsFolderPath = requestsFolderPath + "/"
    }

    model := Model {
        requestsFolderPath: requestsFolderPath,
    }

    return model
}

func (self Model) Init() tea.Cmd {
    return nil
}

func (self Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case ContentRead:
        os.Exit(200)
        self.content.SetContent(string(msg))
    }

    newContentModel, cmd := self.content.Update(msg)
    self.content = newContentModel

    return self, cmd
}

func (self Model) View() string {
    return self.content.View()
}

func (self *Model) SetDimensions(width int, height int) {
    self.content.Width = width
    self.content.Height = height / 2
}

func (self *Model) SetContent(requestName string) {
    path := self.requestsFolderPath + requestName + ".hurl"
    self.content.SetContent(self.readSelectedRequestContent(path))
}

func (self *Model) readSelectedRequestContent(fullFilePath string) string {
    bytes, err := os.ReadFile(fullFilePath)

    if err != nil {
        log.Fatal(err)
    }

    return string(bytes)
}
