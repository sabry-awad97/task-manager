package views

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sabry-awad97/task-manager/internal/tui/models"
)

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))

	mainContainerStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240")).
				Margin(0).
				Padding(1)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1).
			Padding(0, 1).
			Align(lipgloss.Center)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Align(lipgloss.Center)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			MarginTop(1)
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	New    key.Binding
	Delete key.Binding
	Edit   key.Binding
	Help   key.Binding
	Quit   key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new task"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete task"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit task"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q", "quit"),
	),
}

// Add these methods after the keyMap struct definition
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},     // First column
		{k.New, k.Edit},    // Second column
		{k.Delete, k.Help}, // Third column
		{k.Quit},           // Fourth column
	}
}

type MainViewModel struct {
	table    table.Model
	tasks    []models.Task
	help     help.Model
	showHelp bool
	width    int
	height   int
}

func NewMainViewModel() MainViewModel {
	columns := []table.Column{
		{Title: "Title", Width: 30},
		{Title: "Due Date", Width: 12},
		{Title: "Priority", Width: 12},
		{Title: "Status", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	t.SetStyles(table.Styles{
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).
			Padding(0, 1),
		Selected: lipgloss.NewStyle().
			Background(lipgloss.Color("205")).
			Foreground(lipgloss.Color("0")).
			Bold(true),
		Cell: lipgloss.NewStyle().
			Padding(0, 1),
	})

	return MainViewModel{
		table: t,
		help:  help.New(),
	}
}

func (m MainViewModel) Init() tea.Cmd {
	return nil
}

func (m MainViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetWidth(m.width - 4)
		m.table.SetHeight(m.height - 6)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Help):
			m.showHelp = !m.showHelp
			return m, nil
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m MainViewModel) View() string {
	var b strings.Builder

	// Content container
	var content strings.Builder

	// Title
	content.WriteString(titleStyle.Render("✨ Task Manager ✨"))
	content.WriteString("\n")

	// Table (main content)
	content.WriteString(m.table.View())
	content.WriteString("\n")

	// Status bar
	status := fmt.Sprintf("%d tasks • Press ? for help", len(m.tasks))
	content.WriteString(statusStyle.Render(status))

	// Help menu
	if m.showHelp {
		content.WriteString("\n")
		content.WriteString(helpStyle.Render(m.help.View(keys)))
	}

	// Wrap content in main container
	mainContainer := mainContainerStyle.
		Width(m.width - 4).
		Height(m.height - 2).
		Render(content.String())

	b.WriteString(mainContainer)

	return baseStyle.
		Width(m.width).
		Height(m.height).
		Render(b.String())
}

func (m *MainViewModel) UpdateTasks(tasks []models.Task) {
	m.tasks = tasks
	rows := make([]table.Row, len(tasks))

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].DueDate.Before(tasks[j].DueDate)
	})

	for i, task := range tasks {
		priorityStyle := lipgloss.NewStyle()
		switch task.Priority {
		case models.Low:
			priorityStyle = priorityStyle.Foreground(lipgloss.Color("42"))
		case models.Medium:
			priorityStyle = priorityStyle.Foreground(lipgloss.Color("214"))
		case models.High:
			priorityStyle = priorityStyle.Foreground(lipgloss.Color("196"))
		}

		status := "Pending"
		if task.Completed {
			status = "Done"
		}

		rows[i] = table.Row{
			task.Title,
			task.DueDate.Format("2006-01-02"),
			priorityStyle.Render(task.Priority.String()),
			status,
		}
	}

	m.table.SetRows(rows)
}
