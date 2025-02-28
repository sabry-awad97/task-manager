package tui

import (
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sabry-awad97/task-manager/internal/storage"
	"github.com/sabry-awad97/task-manager/internal/tui/models"
	"github.com/sabry-awad97/task-manager/internal/tui/views"
)

type View int

const (
	MainView View = iota
	FormView
	DetailView
	ErrorView
)

func (v View) String() string {
	switch v {
	case MainView:
		return "main"
	case FormView:
		return "form"
	case DetailView:
		return "detail"
	default:
		return "unknown"
	}
}

type rootModel struct {
	currentView   View
	width, height int
	mainView      views.MainViewModel
	formView      views.FormViewModel
	detailView    views.DetailViewModel
	store         *storage.JSONStore
	tasks         []models.Task
	errorView     views.ErrorViewModel
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

	case views.ShowDetailMsg:
		m.detailView = views.NewDetailViewModel(msg.Task)
		m.currentView = DetailView
		return m, nil

	case views.ToggleTaskMsg:
		// Find and toggle the task
		for i := range m.tasks {
			if m.tasks[i].ID == msg.TaskID {
				m.tasks[i].Completed = !m.tasks[i].Completed
				break
			}
		}

		// Update storage
		m.store.Save(m.tasks)

		// Update main view
		m.mainView.UpdateTasks(m.tasks)
		return m, nil

	case views.EditTaskMsg:
		m.formView = views.NewFormViewModel()
		m.formView.InitForEdit(msg.Task)
		m.currentView = FormView
		return m, nil

	case views.DeleteTaskMsg:
		// Find and remove the task
		for i, task := range m.tasks {
			if task.ID == msg.TaskID {
				// Remove task from slice
				m.tasks = slices.Delete(m.tasks, i, i+1)
				break
			}
		}

		// Update storage
		m.store.Save(m.tasks)

		// Update main view
		m.mainView.UpdateTasks(m.tasks)
		return m, nil

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

		if m.currentView == DetailView && (msg.String() == "esc" || msg.String() == "q") {
			m.currentView = MainView
			return m, nil
		}
	case error:
		// Handle any error by showing the error view
		m.errorView = views.NewErrorView(msg)
		m.currentView = ErrorView
		return m, m.errorView.Init()
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
				newTask := newFormView.GetTask()

				if m.formView.IsEditing() {
					// Update existing task
					for i, task := range m.tasks {
						if task.ID == newTask.ID {
							m.tasks[i] = newTask
							break
						}
					}
				} else {
					// Add new task
					m.tasks = append(m.tasks, newTask)
				}

				// Update storage and view
				m.store.Save(m.tasks)
				m.mainView.UpdateTasks(m.tasks)
				m.currentView = MainView
				m.formView = views.NewFormViewModel()
			}
		}
		return m, newCmd

	case DetailView:
		newModel, cmd := m.detailView.Update(msg)
		if newDetailView, ok := newModel.(views.DetailViewModel); ok {
			m.detailView = newDetailView
			if m.detailView.ShouldReturn() {
				m.currentView = MainView
			}
		}
		return m, cmd

	case ErrorView:
		newModel, cmd := m.errorView.Update(msg)
		if newErrorView, ok := newModel.(views.ErrorViewModel); ok {
			m.errorView = newErrorView
			if m.errorView.ShouldClose() {
				m.currentView = MainView
			}
		}
		return m, cmd
	}

	return m, nil
}

func (m rootModel) View() string {
	switch m.currentView {
	case MainView:
		return m.mainView.View()
	case FormView:
		return m.formView.View()
	case DetailView:
		return m.detailView.View()
	case ErrorView:
		return m.errorView.View()
	default:
		return "Unknown View"
	}
}
