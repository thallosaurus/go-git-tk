package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thallosaurus/go-git-tk/pkg/views"
)

func main() {
	initialModel := views.NewCliModel()

	p := tea.NewProgram(initialModel)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
