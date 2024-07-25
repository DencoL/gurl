package requests

import tea "github.com/charmbracelet/bubbletea"

func handleEnter(model *Model) tea.Cmd {
	selectedRequest, err := model.selectedRequest()
	if err != nil {
		return nil
	}

	if selectedRequest.IsFolder {
		nextFolder := model.currentFolder() + selectedRequest.Name + "/"
		model.folders = append(model.folders, nextFolder)

		return model.readAllRequestsFromCurrentFolder
	} else {
		return model.executeRequest
	}
}

type ExecuteRequest string

func (self *Model) executeRequest() tea.Msg {
	return ExecuteRequest(self.selectedRequestFullPath())
}
