package views

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sabry-awad97/task-manager/internal/tui/models"
)

var (
	detailContainerStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("205")).
				Padding(1, 2)

	detailHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("205")).
				MarginBottom(1)

	detailLabelStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("99"))

	detailValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255"))

	detailTimeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true)

	detailFooterStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Align(lipgloss.Center).
				MarginTop(1)
)

type DetailViewModel struct {
	task         models.Task
	width        int
	height       int
	shouldReturn bool
}

func NewDetailViewModel(task models.Task) DetailViewModel {
	return DetailViewModel{task: task}
}

func (m DetailViewModel) Init() tea.Cmd {
	return nil
}

func (m DetailViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.shouldReturn = true
		}
	}
	return m, nil
}

func (m DetailViewModel) View() string {
	var content strings.Builder

	// Header with task title
	content.WriteString(detailHeaderStyle.Render("📝 Task Details"))
	content.WriteString("\n\n")

	// Format task details
	details := []struct {
		label string
		value string
	}{
		{"Title", m.task.Title},
		{"Description", m.task.Description},
		{"Priority", getPriorityWithIcon(m.task.Priority)},
		{"Status", getStatusWithIcon(m.task.Completed)},
		{"Due Date", formatDate(m.task.DueDate)},
		{"Created", formatDate(m.task.CreatedAt)},
	}

	// Render details
	for _, detail := range details {
		if detail.value != "" {
			content.WriteString(detailLabelStyle.Render(detail.label))
			content.WriteString("\n")
			content.WriteString(detailValueStyle.Render(detail.value))
			content.WriteString("\n\n")
		}
	}

	// Footer
	content.WriteString(detailFooterStyle.Render("Press q or esc to return"))

	// Center the modal
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		detailContainerStyle.Render(content.String()),
	)
}

func (m DetailViewModel) ShouldReturn() bool {
	return m.shouldReturn
}

func formatDate(t time.Time) string {
	return detailTimeStyle.Render(t.Format("Monday, January 2, 2006"))
}

func getPriorityWithIcon(p models.PriorityLevel) string {
	icons := map[models.PriorityLevel]string{
		models.Low:    "🟢",
		models.Medium: "🟡",
		models.High:   "🔴",
	}
	return fmt.Sprintf("%s %s", icons[p], p.String())
}

func getStatusWithIcon(completed bool) string {
	if completed {
		return "✅ Done"
	}
	return "⏳ Pending"
}
