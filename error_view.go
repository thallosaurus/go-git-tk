package tuiplayground

import tea "github.com/charmbracelet/bubbletea"

type errorview struct {
	Parent richmodel
	Err    error
}

func (ev errorview) Init() tea.Cmd {
	return nil
}

func (ev errorview) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "enter":
			return ev, changeView(ev.Parent)
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

func (ev errorview) GetKeymapString() string {
	return "enter - back to prev screen"
}
