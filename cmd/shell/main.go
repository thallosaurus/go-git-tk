package main

import (
	"flag"
	"fmt"
	"go-git-tk/pkg/config"
	"go-git-tk/pkg/views"

	tea "github.com/charmbracelet/bubbletea"
)

// var debugFlag = flag.Bool("d", false, "enable debug mode")
var configFlag = flag.String("C", "/etc/gittk/config.json", "specify config to use (default: \"/etc/gittk/config.json\")")

func main() {
	flag.Parse()
	config.InitConfig(*configFlag)

	initialModel := views.NewCliModel()
	initialModel.SetKeyMappingsShown(false)

	p := tea.NewProgram(initialModel, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
