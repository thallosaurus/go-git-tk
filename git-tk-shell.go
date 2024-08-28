package main

import (
	"fmt"
	"go-git-tk/pkg/views"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	initialModel := views.NewCliModel()
	initialModel.SetKeyMappingsShown(false)

	p := tea.NewProgram(initialModel)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
