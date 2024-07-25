package requests

import (
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/exp/slices"
)

func handleGoBack(model *Model) tea.Cmd {
	if model.isInRootFolder() {
		return nil
	}

	foldersLength := len(model.folders)
	model.folders = slices.Delete(model.folders, foldersLength-1, foldersLength)

	return model.readAllRequestsFromCurrentFolder
}

func (self *Model) isInRootFolder() bool {
	return len(self.folders) == 1
}
