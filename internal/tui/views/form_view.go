package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sabry-awad97/task-manager/internal/tui/models"
)

var (
	focusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))

	blurredStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	cursorStyle = focusedStyle

	_ = lipgloss.NewStyle()

	focusedButton = focusedStyle.
			Border(lipgloss.RoundedBorder()).
			Padding(0, 3)

	blurredButton = blurredStyle.
			Border(lipgloss.RoundedBorder()).
			Padding(0, 3)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Italic(true)

	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1).
			Width(42)

	activeInputStyle = inputStyle.
				BorderForeground(lipgloss.Color("205"))

	formContainerStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240")).
				Padding(0, 1).
				MarginLeft(2).
				MarginRight(2)

	labelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99"))

	inputContainerStyle = lipgloss.NewStyle().
				MarginBottom(1)

	buttonContainerStyle = lipgloss.NewStyle().
				MarginTop(1).
				Align(lipgloss.Center)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Align(lipgloss.Center).
			MarginTop(1)

	selectStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1).
			Width(42)

	activeSelectStyle = selectStyle.
				BorderForeground(lipgloss.Color("205"))

	optionStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1)

	selectedOptionStyle = optionStyle.
				Background(lipgloss.Color("205")).
				Foreground(lipgloss.Color("0"))

	priorityOptionStyle = map[models.PriorityLevel]lipgloss.Style{
		models.Low:    optionStyle.Foreground(lipgloss.Color("42")),
		models.Medium: optionStyle.Foreground(lipgloss.Color("214")),
		models.High:   optionStyle.Foreground(lipgloss.Color("196")),
	}
)

type FormViewModel struct {
	title       textinput.Model
	description textinput.Model
	dueDate     textinput.Model
	priority    int
	focusIndex  int
	errors      map[string]string
	width       int
	height      int
	done        bool
	isEditing   bool
	taskID      string // Store original task ID when editing
}

func (m FormViewModel) Done() bool { return m.done }
func (m FormViewModel) IsEditing() bool { return m.isEditing }


func NewFormViewModel() FormViewModel {
	title := textinput.New()
	title.Placeholder = "Enter task title"
	title.Focus()
	title.CharLimit = 50
	title.Width = 40
	title.Cursor.Style = cursorStyle

	description := textinput.New()
	description.Placeholder = "Enter task description"
	description.CharLimit = 100
	description.Width = 40
	description.Cursor.Style = cursorStyle

	dueDate := textinput.New()
	dueDate.Placeholder = "YYYY-MM-DD"
	dueDate.Width = 40
	dueDate.Cursor.Style = cursorStyle

	return FormViewModel{
		title:       title,
		description: description,
		dueDate:     dueDate,
		errors:      make(map[string]string),
		isEditing:   false,
	}
}

// Add method to initialize form for editing
func (m *FormViewModel) InitForEdit(task models.Task) {
	m.title.SetValue(task.Title)
	m.description.SetValue(task.Description)
	m.dueDate.SetValue(task.DueDate.Format("2006-01-02"))
	m.priority = int(task.Priority)
	m.isEditing = true
	m.taskID = task.ID
}

func (m FormViewModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m FormViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "up", "down":
			// Handle focus change
			s := msg.String()
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > 4 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = 4
			}

			// Blur all inputs
			m.title.Blur()
			m.description.Blur()
			m.dueDate.Blur()

			// Focus the active input
			switch m.focusIndex {
			case 0:
				cmd = m.title.Focus()
			case 1:
				cmd = m.description.Focus()
			case 2:
				cmd = m.dueDate.Focus()
			}
			return m, cmd

		case "enter":
			if m.focusIndex == 4 && m.validate() {
				m.done = true
				return m, nil
			}
			// Handle enter key for field navigation
			if m.focusIndex < 4 {
				// Blur current input
				switch m.focusIndex {
				case 0:
					m.title.Blur()
				case 1:
					m.description.Blur()
				case 2:
					m.dueDate.Blur()
				}

				// Move to next field
				m.focusIndex++

				// Focus next input
				switch m.focusIndex {
				case 1:
					cmd = m.description.Focus()
				case 2:
					cmd = m.dueDate.Focus()
				}
				return m, cmd
			}

		case "left", "right":
			if m.focusIndex == 3 {
				if msg.String() == "left" {
					m.priority--
					if m.priority < 0 {
						m.priority = 2
					}
				} else {
					m.priority++
					if m.priority > 2 {
						m.priority = 0
					}
				}
			}
		}
	}

	// Only update active input
	switch m.focusIndex {
	case 0:
		m.title, cmd = m.title.Update(msg)
	case 1:
		m.description, cmd = m.description.Update(msg)
	case 2:
		m.dueDate, cmd = m.dueDate.Update(msg)
	}

	return m, cmd
}

func (m FormViewModel) View() string {
	var b strings.Builder

	// Form content
	var content strings.Builder

	// Title header
	title := "✨ New Task"
	if m.isEditing {
		title = "✏️ Edit Task"
	}
	content.WriteString(titleStyle.Render(title))
	content.WriteString("\n")

	// Title input
	content.WriteString(inputContainerStyle.Render(
		labelStyle.Render("Title") + "\n" +
			m.renderInput(m.title, 0, "title"),
	))

	// Description input
	content.WriteString(inputContainerStyle.Render(
		labelStyle.Render("Description") + "\n" +
			m.renderInput(m.description, 1, ""),
	))

	// Due date input
	content.WriteString(inputContainerStyle.Render(
		labelStyle.Render("Due Date") + "\n" +
			m.renderInput(m.dueDate, 2, "dueDate"),
	))

	// Priority selection
	content.WriteString(inputContainerStyle.Render(
		labelStyle.Render("Priority") + "\n" +
			m.renderPriorities(),
	))

	// Save button
	content.WriteString(buttonContainerStyle.Render(
		m.renderSaveButton(),
	))

	// Wrap the content in a container
	formWidth := 46 // Adjust this value to control overall form width
	b.WriteString(formContainerStyle.Width(formWidth).Render(content.String()))

	// Footer with keyboard hints
	hint := "↑/↓: Navigate • Tab: Next • Esc: Cancel"
	b.WriteString("\n")
	b.WriteString(footerStyle.Render(hint))

	return baseStyle.Width(m.width).Height(m.height).Render(b.String())
}

func (m *FormViewModel) validate() bool {
	m.errors = make(map[string]string)
	valid := true

	if strings.TrimSpace(m.title.Value()) == "" {
		m.errors["title"] = "Title is required"
		valid = false
	}

	if strings.TrimSpace(m.dueDate.Value()) != "" {
		_, err := time.Parse("2006-01-02", m.dueDate.Value())
		if err != nil {
			m.errors["dueDate"] = "Invalid date format (YYYY-MM-DD)"
			valid = false
		}
	}

	return valid
}

func (m *FormViewModel) GetTask() models.Task {
	dueDate, _ := time.Parse("2006-01-02", m.dueDate.Value())
	task := models.NewTask(
		m.title.Value(),
		m.description.Value(),
		dueDate,
		models.PriorityLevel(m.priority),
	)

	// Preserve original ID if editing
	if m.isEditing {
		task.ID = m.taskID
	}

	return task
}

func (m FormViewModel) renderPriorities() string {
	priorities := []struct {
		value models.PriorityLevel
		label string
		icon  string
	}{
		{models.Low, "Low", "🟢"},
		{models.Medium, "Medium", "🟡"},
		{models.High, "High", "🔴"},
	}

	style := selectStyle
	if m.focusIndex == 3 {
		style = activeSelectStyle
	}

	// Build options list
	var options []string
	for i, p := range priorities {
		optStyle := priorityOptionStyle[p.value]
		if i == m.priority {
			if m.focusIndex == 3 {
				optStyle = selectedOptionStyle
			} else {
				optStyle = optStyle.Bold(true)
			}
		}
		option := fmt.Sprintf("%s %s", p.icon, p.label)
		options = append(options, optStyle.Render(option))
	}

	// Add navigation hint
	content := strings.Join(options, " │ ")
	if m.focusIndex == 3 {
		content += blurredStyle.Render(" (← → to select)")
	}

	return style.Render(content)
}

func (m FormViewModel) renderInput(input textinput.Model, index int, errorKey string) string {
	var b strings.Builder

	if m.focusIndex == index {
		b.WriteString(activeInputStyle.Render(input.View()))
	} else {
		b.WriteString(inputStyle.Render(input.View()))
	}

	if err := m.errors[errorKey]; err != "" && errorKey != "" {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render(err))
	}

	return b.String()
}

func (m FormViewModel) renderSaveButton() string {
	style := blurredButton
	if m.focusIndex == 4 {
		style = focusedButton
	}
	return style.Render("💾 Save")
}
