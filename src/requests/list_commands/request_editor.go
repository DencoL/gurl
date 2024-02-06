package listcommands

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

func OpenEditor(file string) tea.Cmd {
	cmd := exec.Command("nvim", file)
	return tea.ExecProcess(cmd, func(_ error) tea.Msg {
        return nil
	})
}
