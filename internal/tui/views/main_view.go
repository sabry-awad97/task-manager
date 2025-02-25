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
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("99")).
			Bold(true).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(0, 1)
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
	loaded   bool
}

func NewMainViewModel() MainViewModel {
	columns := []table.Column{
		{Title: "Title", Width: 20},
		{Title: "Due Date", Width: 15},
		{Title: "Priority", Width: 10},
		{Title: "Status", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	t.SetStyles(table.Styles{
		Header:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99")),
		Selected: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")),
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
		m.table.SetWidth(msg.Width - 2)
		m.table.SetHeight(msg.Height - 4) // Leave room for title and status
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

	// Title
	b.WriteString(titleStyle.Render("Task Manager"))
	b.WriteString("\n")

	// Table
	b.WriteString(m.table.View())
	b.WriteString("\n")

	// Status bar
	status := fmt.Sprintf("%d tasks • Press ? for help", len(m.tasks))
	b.WriteString(statusBarStyle.Render(status))

	// Help
	if m.showHelp {
		b.WriteString("\n")
		b.WriteString(m.help.View(keys))
	}

	return baseStyle.Width(m.width).Height(m.height).Render(b.String())
}

func (m *MainViewModel) UpdateTasks(tasks []models.Task) {
	m.tasks = tasks
	rows := make([]table.Row, len(tasks))

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].DueDate.Before(tasks[j].DueDate)
	})

	for i, task := range tasks {
		status := "Pending"
		if task.Completed {
			status = "Done"
		}

		rows[i] = table.Row{
			task.Title,
			task.DueDate.Format("2006-01-02"),
			task.Priority.String(),
			status,
		}
	}

	m.table.SetRows(rows)
}
