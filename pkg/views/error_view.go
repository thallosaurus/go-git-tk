package views

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	error_ok_key key.Binding = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter/⏎", "ok"),
	)
)

type errorview struct {
	Parent Richmodel
	Err    error
}

func (ev errorview) Init() tea.Cmd {
	return nil
}

func (ev errorview) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(m, error_ok_key):
			return ev, ChangeView(ev.Parent)
		}
	}

	return ev, nil
}

func (ev errorview) View() string {
	s := "An error occured:\n"
	s += ev.Err.Error()
	s += "\n\nPress Enter to continue"

	return s
}

func (ev errorview) GetKeymapString() []key.Binding {
	return []key.Binding{
		error_ok_key,
	}
}
