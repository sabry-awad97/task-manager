package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sabry-awad97/task-manager/internal/storage"
	"github.com/sabry-awad97/task-manager/internal/tui/models"
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
	formView      views.FormViewModel
	store         *storage.JSONStore
	tasks         []models.Task
}

func NewRootModel() rootModel {
	store := storage.NewJSONStore("tasks.json")
	tasks, _ := store.Load() // Load existing tasks

	mainView := views.NewMainViewModel()
	mainView.UpdateTasks(tasks)

	return rootModel{
		currentView: MainView,
		mainView:    mainView,
		formView:    views.NewFormViewModel(),
		store:       store,
		tasks:       tasks,
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
		if m.currentView == MainView {
			newModel, newCmd := m.mainView.Update(msg)
			if newMainView, ok := newModel.(views.MainViewModel); ok {
				m.mainView = newMainView
			}
			return m, newCmd
		}

		newModel, newCmd := m.formView.Update(msg)
		if newFormView, ok := newModel.(views.FormViewModel); ok {
			m.formView = newFormView
		}
		return m, newCmd

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if m.currentView == MainView && msg.String() == "n" {
			m.currentView = FormView
			return m, nil
		}

		if m.currentView == FormView && msg.String() == "esc" {
			m.currentView = MainView
			return m, nil
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
		newModel, newCmd := m.formView.Update(msg)
		if newFormView, ok := newModel.(views.FormViewModel); ok {
			m.formView = newFormView
			if newFormView.Done() {
				// Create new task
				newTask := newFormView.GetTask()
				m.tasks = append(m.tasks, newTask)

				// Update storage
				m.store.Save(m.tasks)

				// Update main view
				m.mainView.UpdateTasks(m.tasks)

				// Reset form and return to main view
				m.currentView = MainView
				m.formView = views.NewFormViewModel()
			}
		}
		return m, newCmd
	}

	return m, nil
}

func (m rootModel) View() string {
	switch m.currentView {
	case MainView:
		return m.mainView.View()
	case FormView:
		return m.formView.View()
	default:
		return "Unknown View"
	}
}
