package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type View int

const (
	MainView View = iota
	FormView
)

func (v View) String() string {
	switch v {
	case MainView:
		return "main"
	case FormView:
		return "form"
	default:
		return "unknown"
	}
}

type rootModel struct {
	currentView   View
	width, height int
}

func NewRootModel() rootModel {
	return rootModel{
		currentView: MainView,
	}
}

func (m rootModel) Init() tea.Cmd {
	return nil
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	switch m.currentView {
	case MainView:
		return m, nil

	case FormView:
		return m, nil
	}

	return m, nil
}

func (m rootModel) View() string {
	return "Hello World!"
}
