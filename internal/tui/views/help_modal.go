package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	helpModalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205")).
			Padding(1, 2)

	helpHeadingStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Bold(true).
				Align(lipgloss.Center)

	helpSectionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("99")).
				Bold(true).
				PaddingTop(1)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

type HelpSection struct {
	Title string
	Items []HelpItem
}

type HelpItem struct {
	Key         string
	Description string
}

func RenderHelpModal(width, height int) string {
	sections := []HelpSection{
		{
			Title: "Navigation",
			Items: []HelpItem{
				{"↑/k", "Move up"},
				{"↓/j", "Move down"},
				{"tab", "Next field"},
				{"shift+tab", "Previous field"},
			},
		},
		{
			Title: "Tasks",
			Items: []HelpItem{
				{"n", "New task"},
				{"d", "Delete task"},
				{"e", "Edit task"},
				{"space", "Toggle complete"},
			},
		},
		{
			Title: "General",
			Items: []HelpItem{
				{"?", "Toggle help"},
				{"esc", "Exit/Cancel"},
				{"ctrl+c", "Quit"},
			},
		},
	}

	var content strings.Builder

	// Title
	content.WriteString(helpHeadingStyle.Render("🎯 Keyboard Shortcuts"))
	content.WriteString("\n\n")

	// Render each section
	for i, section := range sections {
		if i > 0 {
			content.WriteString("\n")
		}

		content.WriteString(helpSectionStyle.Render(section.Title))
		content.WriteString("\n")

		// Calculate maximum key length for alignment
		maxKeyLen := 0
		for _, item := range section.Items {
			if len(item.Key) > maxKeyLen {
				maxKeyLen = len(item.Key)
			}
		}

		// Render items
		for _, item := range section.Items {
			key := helpKeyStyle.Render(item.Key)
			padding := strings.Repeat(" ", maxKeyLen-len(item.Key)+2)
			desc := helpDescStyle.Render(item.Description)
			content.WriteString(key + padding + desc + "\n")
		}
	}

	// Center the modal
	modalWidth := 40
	modalHeight := strings.Count(content.String(), "\n") + 3

	modal := helpModalStyle.
		Width(modalWidth).
		Render(content.String())

	// Calculate position for centering
	leftPadding := (width - modalWidth) / 2
	if leftPadding < 0 {
		leftPadding = 0
	}
	topPadding := (height - modalHeight) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	// Apply padding for centering
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)
}
