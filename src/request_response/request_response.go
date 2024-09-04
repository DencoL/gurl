package requestresponse

import (
	"bytes"
	"encoding/json"
	"gurl/requests"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	content viewport.Model
}

func New() Model {
	return Model{
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

	case requests.ExecuteRequest:
		return self, self.executeRequest(string(msg))

	case RequestExecuted:
		res := string(msg)
		if res == "" {
			res = "<EMPTY RESPONSE>"
		} else {
			self.content.SetContent(toPrettyJson(res))
		}
	}

	self.content, cmd = self.content.Update(msg)
	cmds = append(cmds, cmd)

	return self, tea.Batch(cmds...)
}

func toPrettyJson(response string) string {
	var prettyJSON bytes.Buffer
	toJsonError := json.Indent(&prettyJSON, []byte(response), "", "  ")
	if toJsonError == nil {
		return string(prettyJSON.Bytes())
	} else {
		return response
	}
}

func (self Model) View() string {
	return self.content.View()
}

type RequestExecuted string

func (self *Model) executeRequest(requestFilePath string) tea.Cmd {
	return func() tea.Msg {
		return RequestExecuted(runHurl(requestFilePath))
	}
}
