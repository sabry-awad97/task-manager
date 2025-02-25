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
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	New    key.Binding
	Delete key.Binding
	Edit   key.Binding
	Help   key.Binding
	Quit   key.Binding
	Enter  key.Binding
	Space  key.Binding
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
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "view details"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "toggle completed"),
	),
}

// Add these methods after the keyMap struct definition
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter},
		{k.New, k.Edit, k.Space},
		{k.Delete, k.Help},
		{k.Quit},
	}
}

// Add event for view transition
type ShowDetailMsg struct {
	Task models.Task
}

type ToggleTaskMsg struct {
	TaskID string
}

// Add edit message type
type EditTaskMsg struct {
	Task models.Task
}

type MainViewModel struct {
	table    table.Model
	tasks    []models.Task
	help     help.Model
	showHelp bool
	width    int
	height   int
	mouseX   int // Add mouse position tracking
	mouseY   int // Add mouse position tracking
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
		case msg.Type == tea.KeyEnter:
			if task, ok := m.SelectedTask(); ok {
				return m, func() tea.Msg {
					return ShowDetailMsg{Task: task}
				}
			}
		case key.Matches(msg, keys.Space):
			if task, ok := m.SelectedTask(); ok {
				return m, func() tea.Msg {
					return ToggleTaskMsg{TaskID: task.ID}
				}
			}
		case key.Matches(msg, keys.Edit):
			if task, ok := m.SelectedTask(); ok {
				return m, func() tea.Msg {
					return EditTaskMsg{Task: task}
				}
			}
		}

	case tea.MouseMsg:
		m.mouseX = msg.X
		m.mouseY = msg.Y

		switch msg.Action {
		case tea.MouseActionPress:
			switch msg.Button {
			case tea.MouseButtonLeft:
				if m.isClickInTable(msg) {
					rowIdx := m.getClickedRowIndex(msg)
					if rowIdx >= 0 && rowIdx < len(m.tasks) {
						m.table.SetCursor(rowIdx)
						if task, ok := m.SelectedTask(); ok {
							return m, func() tea.Msg {
								return ShowDetailMsg{Task: task}
							}
						}
					}
				}
			case tea.MouseButtonRight:
				if m.isClickInTable(msg) {
					rowIdx := m.getClickedRowIndex(msg)
					if rowIdx >= 0 && rowIdx < len(m.tasks) {
						m.table.SetCursor(rowIdx)
						if task, ok := m.SelectedTask(); ok {
							return m, func() tea.Msg {
								return ToggleTaskMsg{TaskID: task.ID}
							}
						}
					}
				}
			}
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// Add helper methods for mouse interaction
func (m MainViewModel) isClickInTable(msg tea.MouseMsg) bool {
	// Adjust these values based on your layout
	tableTop := 2 // Account for title and padding
	tableBottom := tableTop + m.table.Height()

	return msg.Y >= tableTop && msg.Y <= tableBottom
}

func (m MainViewModel) getClickedRowIndex(msg tea.MouseMsg) int {
	tableTop := 2               // Same as in isClickInTablesClickInTable
	return msg.Y - tableTop - 1 // -1 for header row
}

func (m MainViewModel) View() string {
	if m.showHelp {
		return RenderHelpModal(m.width, m.height)
	}

	// Pre-allocate builders with estimated capacity
	content := strings.Builder{}
	content.Grow(m.width * m.height)

	// Build content in single pass
	content.WriteString(titleStyle.Render("✨ Task Manager ✨"))
	content.WriteByte('\n')
	content.WriteString(m.table.View())
	content.WriteByte('\n')
	content.WriteString(statusStyle.Render(fmt.Sprintf("%d tasks • Press ? for help", len(m.tasks))))

	// Apply container styles in sequence
	return baseStyle.
		Width(m.width).
		Height(m.height).
		Render(
			mainContainerStyle.
				Width(m.width - 4).
				Height(m.height - 2).
				Render(content.String()),
		)
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

func (m MainViewModel) SelectedTask() (models.Task, bool) {
	if len(m.tasks) == 0 {
		return models.Task{}, false
	}
	selected := m.table.SelectedRow()
	for _, task := range m.tasks {
		if task.Title == selected[0] {
			return task, true
		}
	}
	return models.Task{}, false
}
