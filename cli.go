package tuiplayground

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type climodel struct {
	selectedView richmodel
}

func NewCliModel() climodel {
	return climodel{
		selectedView: nil,
	}
}

type TickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (c climodel) Init() tea.Cmd {
	//return nil
	return changeView(NewHomeView())
}

type (
	switchView struct {
		model richmodel
	}
	closeChild struct{}
)

func changeView(m richmodel) tea.Cmd {
	return func() tea.Msg {
		return switchView{
			model: m,
		}
	}
}

func closeView() tea.Cmd {
	return func() tea.Msg {
		return closeChild{}
	}
}

func (c climodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case switchView:
		c.selectedView = m.model
		return c, tea.Batch(c.selectedView.Init(), tea.ClearScreen)

	case closeChild:
		c.selectedView = nil
		return c, nil

	case tea.KeyMsg:
		switch m := msg.(tea.KeyMsg).String(); m {
		case "ctrl+c":
			return c, tea.Quit
		}
	}

	// if there is a selected child view, push update to the child
	if c.selectedView != nil {
		cc, cmd := c.selectedView.Update(msg)

		c.selectedView = cc.(richmodel)

		return c, cmd
	}

	return c, nil
}

func (c climodel) View() string {
	if c.selectedView != nil {
		return c.selectedView.View() + "\n\n" + c.GetKeymapString()
	} else {
		return fmt.Sprintf("root view, dbg: %v", c)
	}
}

func (c climodel) GetKeymapString() string {
	var s string
	s += "ctrl+c - quit"

	if c.selectedView != nil {
		s += ", "
		s += c.selectedView.GetKeymapString()
	}

	return helpStyle.Render(s)
}
