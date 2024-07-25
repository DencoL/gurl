package requests

import (
	listcommands "gurl/requests/list_commands"

	tea "github.com/charmbracelet/bubbletea"
)

func handleEdit(model *Model) tea.Cmd {
	if selectedRequest, _ := model.selectedRequest(); selectedRequest.IsFolder {
		return nil
	}
	return tea.Batch(tea.HideCursor, listcommands.OpenEditor(model.selectedRequestFullPath()))
}
