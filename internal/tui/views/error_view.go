package views

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	errorViewStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("196")).
			Padding(1, 2)

	errorTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	errorMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241"))

	errorHintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true)
)

type ErrorViewModel struct {
	err         error
	width       int
	height      int
	showTime    time.Time
	shouldClose bool
}

func NewErrorView(err error) ErrorViewModel {
	return ErrorViewModel{
		err:      err,
		showTime: time.Now(),
	}
}

func (m ErrorViewModel) Init() tea.Cmd {
	// Auto-close error after 3 seconds
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return errorTimeoutMsg(t)
	})
}

type errorTimeoutMsg time.Time

func (m ErrorViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case errorTimeoutMsg:
		m.shouldClose = true

	case tea.KeyMsg:
		if msg.String() == "esc" || msg.String() == "enter" {
			m.shouldClose = true
		}
	}
	return m, nil
}

func (m ErrorViewModel) View() string {
	var content strings.Builder

	// Error icon and title
	content.WriteString(errorTitleStyle.Render("❌ Error"))
	content.WriteString("\n\n")

	// Error message
	content.WriteString(errorMessageStyle.Render(m.err.Error()))
	content.WriteString("\n\n")

	// Hint
	content.WriteString(errorHintStyle.Render("Press Enter or Esc to continue"))

	// Center the modal
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		errorViewStyle.Render(content.String()),
	)
}

func (m ErrorViewModel) ShouldClose() bool {
	return m.shouldClose
}
