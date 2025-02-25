package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sabry-awad97/task-manager/internal/tui/views"
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
	mainView      views.MainViewModel
}

func NewRootModel() rootModel {
	return rootModel{
		currentView: MainView,
		mainView:    views.NewMainViewModel(),
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
		newModel, newCmd := m.mainView.Update(msg)
		if newMainView, ok := newModel.(views.MainViewModel); ok {
			m.mainView = newMainView
		}
		return m, newCmd

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	switch m.currentView {
	case MainView:
		newModel, newCmd := m.mainView.Update(msg)
		if newMainView, ok := newModel.(views.MainViewModel); ok {
			m.mainView = newMainView
		}
		return m, newCmd

	case FormView:
		return m, nil
	}

	return m, nil
}

func (m rootModel) View() string {
	switch m.currentView {
	case MainView:
		return m.mainView.View()
	case FormView:
		return "Form View"
	default:
		return "Unknown View"
	}
}
