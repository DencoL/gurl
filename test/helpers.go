package test

import tea "github.com/charmbracelet/bubbletea"

func IsMsgOfType[TCommand tea.Msg](cmd tea.Cmd) bool {
    if cmd == nil {
        return false
    }

    _, ok := cmd().(TCommand)

    return ok
}
