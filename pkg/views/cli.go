package views

import (
	"fmt"

	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	ssh_base_domain string = "localhost"
	ssh_user        string = "git"

	term_width  int = 20
	term_height int = 14

	cli_quit key.Binding = key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "exit"),
	)
)

type climodel struct {
	selectedView Richmodel
	show_keys    bool
	ready        bool
}

type Richmodel interface {
	tea.Model
	GetKeymapString() []key.Binding
}

func NewCliModel() climodel {
	return climodel{
		selectedView: nil,
		show_keys:    false,
		ready:        false,
	}
}

func (c *climodel) SetKeyMappingsShown(b bool) {
	c.show_keys = b
}

func (c climodel) Init() tea.Cmd {
	//return nil
	return ChangeView(MakeHomeList())
}

type (
	switchView struct {
		model Richmodel
	}
	closeChild struct{}
)

func ChangeView(m Richmodel) tea.Cmd {
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
	case tea.WindowSizeMsg:
		term_width = m.Width
		term_height = m.Height

		if !c.ready {
			c.ready = true
		}

	case switchView:
		c.selectedView = m.model
		return c, tea.Batch(c.selectedView.Init(), tea.ClearScreen)

	case closeChild:
		c.selectedView = nil
		return c, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(m, cli_quit):
			return c, tea.Batch(tea.ClearScreen, tea.Quit)
		}
	}

	// if there is a selected child view, push update to the child
	if c.selectedView != nil {
		cc, cmd := c.selectedView.Update(msg)

		c.selectedView = cc.(Richmodel)

		return c, cmd
	}

	return c, nil
}

func (c climodel) View() string {
	if c.selectedView != nil && c.ready {
		s := c.selectedView.View()
		/*if c.show_keys {
			s += "\n" + c.GetKeymapString()
		}*/

		km := append(c.GetKeymapString(), c.selectedView.GetKeymapString()...)

		var sb []string
		for _, val := range km {
			sb = append(sb, fmt.Sprintf("<%s>: %s", val.Help().Key, val.Help().Desc))
		}

		quickhelp := keymapStyle.Render(strings.Join(sb, " • "))
		v := fmt.Sprintf("%s\n%s", s, quickhelp)

		return v
	} else {
		return fmt.Sprintf("root view, dbg: %v", c)
	}
}

func (c climodel) GetKeymapString() []key.Binding {
	var s string

	s += "ctrl+c - quit"

	if c.selectedView != nil {
		s += " • "
		//s += c.selectedView.GetKeymapString() + "\n"
	}

	s = helpStyle.Render(s)

	//return s
	return []key.Binding{
		cli_quit,
	}
}
