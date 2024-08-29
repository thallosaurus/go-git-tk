package main

import (
	"fmt"
	"go-git-tk/pkg/layouts"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	vp      viewport.Model
	borders bool
}

var (
	exit_program key.Binding = key.NewBinding(
		key.WithKeys("q"),
	)
	toggle_borders key.Binding = key.NewBinding(
		key.WithKeys("b"),
	)
)

func (m model) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m model) View() string {
	o := fmt.Sprintf("%s\n%s\n%s",
		layouts.HeaderStyle.Render("I am a header"),
		layouts.ContentStyle.Render(m.vp.View()),
		layouts.FooterStyle.Render("Footer"),
	)
	return layouts.MainStyle.Render(o)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		layouts.UpdateTermSize(msg)
		m.vp.Width = layouts.GetContentInnerWidth()
		m.vp.Height = layouts.GetContentInnerHeight()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, exit_program):
			return m, tea.Quit

		case key.Matches(msg, toggle_borders):
			if m.borders {
				layouts.TurnOnDebugBorders()
			} else {
				layouts.TurnOffDebugBorders()
			}

			m.borders = !m.borders
			return m, tea.WindowSize()
		}
	}

	var cmd tea.Cmd
	m.vp, cmd = m.vp.Update(msg)
	return m, cmd
}

func main() {
	//layouts.TurnOnDebugBorders()
	vp := viewport.New(layouts.GetContentInnerWidth(), layouts.GetContentInnerHeight())
	//initialModel.SetKeyMappingsShown(false)

	f, err := os.ReadFile("README.md")
	if err != nil {
		vp.SetContent(err.Error())
	} else {
		vp.SetContent(string(f))
	}

	initialModel := model{vp: vp, borders: true}

	if initialModel.borders {
		layouts.TurnOnDebugBorders()
	} else {
		layouts.TurnOffDebugBorders()
	}

	p := tea.NewProgram(initialModel, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
