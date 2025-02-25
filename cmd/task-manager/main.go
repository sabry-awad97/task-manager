package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sabry-awad97/task-manager/internal/tui"
)

func main() {
	p := tea.NewProgram(
		tui.NewRootModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(), // Enable mouse motion
		tea.WithMouseAllMotion(),  // Enable all mouse events
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running application: %v", err)
		os.Exit(1)
	}
}
