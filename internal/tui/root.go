package tui

import tea "github.com/charmbracelet/bubbletea"

type rootModel struct {
	currentView   string
	width, height int
}

func NewRootModel() rootModel {
	return rootModel{
		currentView: "main",
	}
}

func (m rootModel) Init() tea.Cmd {
	return nil
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m rootModel) View() string {
	return "Hello World!"
}
