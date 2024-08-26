package main

import (
	"fmt"
	"tuiplayground"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	initialModel := tuiplayground.NewCliModel()

	p := tea.NewProgram(initialModel)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
