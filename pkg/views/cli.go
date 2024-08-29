package views

import (
	"fmt"
	"go-git-tk/pkg/config"
	"go-git-tk/pkg/layouts"

	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

var conf = config.ReadConfig("./scripts/gittk-shell/config.json")

var (
	cli_quit key.Binding = key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "exit"),
	)
)

type climodel struct {
	selectedView Richmodel
	show_keys    bool
	rootPane     viewport.Model
}

type Richmodel interface {
	tea.Model
	GetKeymapString() []key.Binding
	GetHeaderString() string
}

func NewCliModel() climodel {
	vp := viewport.New(layouts.GetContentInnerWidth(), layouts.GetContentInnerHeight())
	vp.KeyMap.Down.SetEnabled(false)
	vp.KeyMap.Up.SetEnabled(false)

	return climodel{
		selectedView: nil,
		show_keys:    false,
		rootPane:     vp,
	}
}

func (c *climodel) SetKeyMappingsShown(b bool) {
	c.show_keys = b
}

func (c climodel) Init() tea.Cmd {
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

func (c climodel) GetHeaderString() string {
	return ""
}

func (c climodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		layouts.UpdateTermSize(m)
		c.rootPane.Width = layouts.GetContentInnerWidth()
		c.rootPane.Height = layouts.GetContentInnerHeight()

	case switchView:
		c.selectedView = m.model
		return c, tea.Batch(c.selectedView.Init(), tea.WindowSize(), tea.ClearScreen)

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
	if c.selectedView != nil {
		c.rootPane.SetContent(c.selectedView.View())

		s := layouts.HeaderStyle.Render(c.selectedView.GetHeaderString()) + "\n"
		//s += layouts.ContentStyle.Render(c.selectedView.View())
		s += layouts.ContentStyle.Render(c.rootPane.View())
		/*if c.show_keys {
			s += "\n" + c.GetKeymapString()
		}*/

		km := append(c.GetKeymapString(), c.selectedView.GetKeymapString()...)

		var sb []string
		for _, val := range km {

			if val.Enabled() {
				sb = append(sb, fmt.Sprintf("<%s>: %s", val.Help().Key, val.Help().Desc))
			}
		}

		quickhelp := layouts.FooterStyle.Render(strings.Join(sb, " • "))
		v := fmt.Sprintf("%s\n%s", s, quickhelp)

		return layouts.MainStyle.Render(v)
	} else {
		return fmt.Sprintf("root view, dbg: %v", c)
	}
}

func (c climodel) GetKeymapString() []key.Binding {
	return []key.Binding{
		cli_quit,
	}
}
