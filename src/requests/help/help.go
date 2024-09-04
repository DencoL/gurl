package help

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type keyMap struct {
	Up           key.Binding
	Down         key.Binding
	Confirm      key.Binding
	Back         key.Binding
	Edit         key.Binding
	Help         key.Binding
	YankResponse key.Binding
}

func newListKeyMap() *keyMap {
	return &keyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "Up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "Down"),
		),
		Confirm: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↵ ", "Open/Send"),
		),
		Back: key.NewBinding(
			key.WithKeys("-"),
			key.WithHelp("-", "Go up folder"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "Edit"),
		),
		YankResponse: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "Yank response"),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Up, k.Down, k.Confirm, k.Edit,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Confirm, k.Edit},
		{k.Back, k.Help},
	}
}

func Keys() *keyMap {
	return newListKeyMap()
}

type Model struct {
	help help.Model
	Keys *keyMap
}

func New() Model {
	h := help.New()
	h.ShowAll = true

	model := Model{
		help: h,
	}

	model.Keys = newListKeyMap()

	return model
}

func (self Model) Init() tea.Cmd {
	return nil
}

func (self Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		self.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		}
	}

	return self, nil
}

func (self Model) View() string {
	return self.help.View(self.Keys)
}
